package evaluator

import (
	"fmt"
	"io/ioutil"
	"math"
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
		enclosed := object.NewEnclosedEnvironment(env)
		obj := e.evalBlockStatement(node, enclosed)
		enclosed.Dispose()

		return obj

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

	case *ast.WhileStatement:
		return e.evalWhileStatement(node, env)

	case *ast.IncludeStatement:
		return e.evalIncludeStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

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
			return right
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

		_, ok := env.Set(node.Name.Value, val)
		if !ok {
			return newError("identifier not found: " + node.Name.Value)
		}

		return val

	case *ast.ArrayLiteral:
		objects := make([]object.Object, len(node.Value))

		for index, expr := range node.Value {
			object := e.Eval(expr, env)
			if isError(object) {
				return object
			}

			objects[index] = object
		}

		return &object.Array{Value: objects}

	case *ast.IndexExpression:
		return e.evalIndexExpression(node, env)
	}

	return NULL
}

func (e *Evaluator) applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newError("expected %d arg(s) got %d", len(fn.Parameters), len(args))
		}
		extendedEnv := e.extendFunctionEnv(fn, args)
		evaluated := e.Eval(fn.Body, extendedEnv)
		extendedEnv.Dispose()
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Inspect())
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

	switch body := node.Body.(type) {
	case *ast.BlockStatement: // to optimize for block statements
		bodyEnv := object.NewEnclosedEnvironment(enclosedEnv)

		for {
			condition := e.Eval(node.Condition, enclosedEnv)
			if isError(condition) {
				return condition
			}

			if !isTruthy(condition) {
				break
			}

			result := e.evalBlockStatement(body, bodyEnv)
			if isReturn(result) {
				return result
			}

			bodyEnv.ClearStore()

			increment := e.Eval(node.Increment, enclosedEnv)
			if isError(increment) {
				return increment
			}
		}

		bodyEnv.Dispose()

	default:
		for {
			condition := e.Eval(node.Condition, enclosedEnv)
			if isError(condition) {
				return condition
			}

			if !isTruthy(condition) {
				break
			}

			result := e.Eval(body, enclosedEnv)
			if isReturn(result) {
				return result
			}

			increment := e.Eval(node.Increment, enclosedEnv)
			if isError(increment) {
				return increment
			}
		}
	}

	enclosedEnv.Dispose()

	return NULL
}

func (e *Evaluator) evalWhileStatement(node *ast.WhileStatement, env *object.Environment) object.Object {
	for {
		condition := e.Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result := e.Eval(node.Body, env)
		if isReturn(result) {
			return result
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
		e.Eval(program, fileEnv)
	}

	env.Includes = append(env.Includes, fileEnv)

	return NULL
}

func (e *Evaluator) evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return e.evalBangOperatorExpression(right)
	case "-":
		return e.evalMinusPrefixOperatorExpression(right)

	default:
		return newError("unknown operator: %s%s", operator, right.Inspect())
	}
}

func (e *Evaluator) evalBangOperatorExpression(right object.Object) object.Object {
	return boolToBoolObject(!isTruthy(right))
}

func (e *Evaluator) evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Inspect())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func (e *Evaluator) evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value

		return e.evalIntegerInfixExpression(operator, leftVal, rightVal)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value

		return e.evalFloatInfixExpression(operator, leftVal, rightVal)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Integer).Value

		return e.evalFloatInfixExpression(operator, leftVal, float64(rightVal))

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Float).Value

		return e.evalFloatInfixExpression(operator, float64(leftVal), rightVal)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return e.evalStringInfixExpression(operator, left, right)

	case operator == "==":
		return boolToBoolObject(left == right)

	case operator == "!=":
		return boolToBoolObject(left != right)

	case operator == "&&":
		return boolToBoolObject(isTruthy(left) && isTruthy(right))

	case operator == "||":
		return boolToBoolObject(isTruthy(left) || isTruthy(right))

	default:
		return newError("unknown operator: %s %s %s",
			left.Inspect(), operator, right.Inspect())
	}
}

func (e *Evaluator) evalIntegerInfixExpression(operator string, left, right int64) object.Object {
	switch operator {
	case "<":
		return boolToBoolObject(left < right)
	case "<=":
		return boolToBoolObject(left <= right)
	case ">":
		return boolToBoolObject(left > right)
	case ">=":
		return boolToBoolObject(left >= right)
	case "==":
		return boolToBoolObject(left == right)
	case "!=":
		return boolToBoolObject(left != right)
	case "+":
		return &object.Integer{Value: left + right}
	case "-":
		return &object.Integer{Value: left - right}
	case "*":
		return &object.Integer{Value: left * right}
	case "/":
		return &object.Integer{Value: left / right}
	case "%":
		return &object.Integer{Value: left % right}
	default:
		return newError("unknown operator: %d %s %d",
			left, operator, right)
	}
}

func (e *Evaluator) evalFloatInfixExpression(operator string, left, right float64) object.Object {
	switch operator {
	case "<":
		return boolToBoolObject(left < right)
	case "<=":
		return boolToBoolObject(left <= right)
	case ">":
		return boolToBoolObject(left > right)
	case ">=":
		return boolToBoolObject(left >= right)
	case "==":
		return boolToBoolObject(left == right)
	case "!=":
		return boolToBoolObject(left != right)
	case "+":
		return &object.Float{Value: left + right}
	case "-":
		return &object.Float{Value: left - right}
	case "*":
		return &object.Float{Value: left * right}
	case "/":
		return &object.Float{Value: left / right}
	case "%":
		return &object.Float{Value: math.Mod(left, right)}
	default:
		return newError("unknown operator: %f %s %f",
			left, operator, right)
	}
}

func (e *Evaluator) evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s",
			left.Inspect(), operator, right.Inspect())
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
			intObj := value.(*object.Integer)
			intObj.Value++

			env.Set(left.Value, intObj)
			return intObj

		default:
			return newError("unknown operator: %s%s", operator, value.Inspect())
		}
	default:
		return newError("unknown operator: %s%s", operator, value.Inspect())
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

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func (e *Evaluator) evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := e.Eval(node.Left, env)
	if isError(left) {
		return left
	}

	index := e.Eval(node.Index, env)
	if isError(index) {
		return index
	}

	switch left.Type() {
	case object.ARRAY_OBJ:
		return e.evalArrayIndexExpression(node, left, index)
	case object.STRING_OBJ:
		return e.evalStringIndexExpression(node, left, index)
	default:
		return newError("index operator not supported: %s", left.Inspect())
	}
}

func (e *Evaluator) evalArrayIndexExpression(node *ast.IndexExpression, array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObj.Value[idx]
}

func (e *Evaluator) evalStringIndexExpression(node *ast.IndexExpression, str, index object.Object) object.Object {
	strObj := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObj.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return &object.String{Value: string([]rune(strObj.Value)[idx])}
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
