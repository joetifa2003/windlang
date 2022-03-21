package evaluator

import (
	"fmt"
	"io/ioutil"
	"wind-vm-go/ast"
	"wind-vm-go/lexer"
	"wind-vm-go/object"
	"wind-vm-go/parser"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type Evaluator struct {
	envManager *object.EnvironmentManager
	includes   []*object.Environment
}

func New(envManager *object.EnvironmentManager) *Evaluator {
	return &Evaluator{
		envManager: envManager,
	}
}

func (e *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node.Statements, env)

	case *ast.BlockStatement:
		return e.evalBlockStatement(node, object.NewEnclosedEnvironment(env))

	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env)

	case *ast.LetStatement:
		val := e.Eval(node.Value, env)
		if isError(val) {
			return val
		}

		env.Let(node.Name.Value, val)

	case *ast.ReturnStatement:
		val := e.Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}

	case *ast.ForStatement:
		return e.evalForStatement(node, env)

	case *ast.IncludeStatement:
		return e.evalIncludeStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return boolToBoolObject(node.Value)

	case *ast.PrefixExpression:
		right := e.Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return e.evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := e.Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := e.Eval(node.Right, env)
		if isError(right) {
			return left
		}

		return e.evalInfixExpression(node.Operator, left, right)

	case *ast.PostfixExpression:
		switch left := node.Left.(type) {
		case *ast.Identifier:
			return e.evalPostfixExpression(node.Operator, left, env)
		default:
			return newError("postfix expression must be identifier")
		}

	case *ast.IfExpression:
		return e.evalIfExpression(node, env)

	case *ast.Identifier:
		return e.evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := e.Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := e.evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return e.applyFunction(function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.AssignExpression:
		val := e.Eval(node.Value, env)
		if isError(val) {
			return val
		}

		_, ok := env.Get(node.Name.Value)
		if !ok {
			return newError("identifier not found: " + node.Name.Value)
		}

		env.Set(node.Name.Value, val)

		return val
	}

	return NULL
}

func (e *Evaluator) applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := e.extendFunctionEnv(fn, args)
		evaluated := e.Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}

}

func (e *Evaluator) extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Let(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := e.Eval(exp, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = e.Eval(statement, env)

		switch result := result.(type) {
		case *object.Error:
			return result

		case *object.ReturnValue:
			return result.Value
		}
	}

	return result
}

func (e *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = e.Eval(statement, env)

		if isReturn(result) {
			return result
		}
	}

	return result
}

func (e *Evaluator) evalForStatement(node *ast.ForStatement, env *object.Environment) object.Object {
	enclosedEnv := object.NewEnclosedEnvironment(env)

	init := e.Eval(node.Initializer, enclosedEnv)
	if isError(init) {
		return init
	}

	for {
		condition := e.Eval(node.Condition, enclosedEnv)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result := e.Eval(node.Body, enclosedEnv)
		if isReturn(result) {
			return result
		}

		increment := e.Eval(node.Increment, enclosedEnv)
		if isError(increment) {
			return increment
		}
	}

	return NULL
}

func (e *Evaluator) evalIncludeStatement(node *ast.IncludeStatement, env *object.Environment) object.Object {

	path := node.Path
	fileEnv, evaluated := e.envManager.Get(path)

	if !evaluated {
		file, _ := ioutil.ReadFile(path)
		input := string(file)
		lexer := lexer.New(input)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		evaluator := New(e.envManager)
		evaluator.Eval(program, fileEnv)
	}

	e.includes = append(e.includes, fileEnv)

	return NULL
}

func (e *Evaluator) evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return e.evalBangOperatorExpression(right)
	case "-":
		return e.evalMinusPrefixOperatorExpression(right)

	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func (e *Evaluator) evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE

	default:
		return FALSE
	}
}

func (e *Evaluator) evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func (e *Evaluator) evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return e.evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return e.evalStringInfixExpression(operator, left, right)

	case operator == "==":
		return boolToBoolObject(left == right)

	case operator == "!=":
		return boolToBoolObject(left != right)

	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "<":
		return boolToBoolObject(leftVal < rightVal)
	case ">":
		return boolToBoolObject(leftVal > rightVal)
	case "==":
		return boolToBoolObject(leftVal == rightVal)
	case "!=":
		return boolToBoolObject(leftVal != rightVal)
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func (e *Evaluator) evalPostfixExpression(operator string, left *ast.Identifier, env *object.Environment) object.Object {
	value := e.Eval(left, env)
	if isError(value) {
		return value
	}

	switch value.Type() {
	case object.INTEGER_OBJ:
		switch operator {
		case "++":
			newValue := value.(*object.Integer).Value + 1
			env.Set(left.Value, &object.Integer{Value: newValue})
			return &object.Integer{Value: newValue}

		default:
			return newError("unknown operator: %s%s", operator, value.Type())
		}
	default:
		return newError("unknown operator: %s%s", operator, value.Type())
	}
}

func (e *Evaluator) evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := e.Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Eval(ie.ThenBranch, env)
	} else if ie.ElseBranch != nil {
		return e.Eval(ie.ElseBranch, env)
	} else {
		return NULL
	}
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	for _, include := range e.includes {
		if val, ok := include.Get(node.Value); ok {
			return val
		}
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false

	case TRUE:
		return true

	case FALSE:
		return false

	default:
		return true
	}
}

func boolToBoolObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}

	return FALSE
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func isReturn(obj object.Object) bool {
	rt := obj.Type()

	if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
		return true
	}

	return false
}
