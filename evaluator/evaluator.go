package evaluator

import (
	"fmt"
	"io/ioutil"
	"math"
	"wind-vm-go/ast"
	"wind-vm-go/lexer"
	"wind-vm-go/parser"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Evaluator struct {
	envManager *EnvironmentManager
}

func New(envManager *EnvironmentManager) *Evaluator {
	return &Evaluator{
		envManager: envManager,
	}
}

func (e *Evaluator) Eval(node ast.Node, env *Environment) Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node.Statements, env)

	case *ast.BlockStatement:
		enclosed := NewEnclosedEnvironment(env)
		obj := e.evalBlockStatement(node, enclosed)

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

		return &ReturnValue{Value: val}

	case *ast.ForStatement:
		return e.evalForStatement(node, env)

	case *ast.WhileStatement:
		return e.evalWhileStatement(node, env)

	case *ast.IncludeStatement:
		return e.evalIncludeStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &Float{Value: node.Value}

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

		return &Function{Parameters: params, Env: env, Body: body}

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
		return &String{Value: node.Value}

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
		objects := make([]Object, len(node.Value))

		for index, expr := range node.Value {
			object := e.Eval(expr, env)
			if isError(object) {
				return object
			}

			objects[index] = object
		}

		return &Array{Value: objects}

	case *ast.IndexExpression:
		return e.evalIndexExpression(node, env)
	}

	return NULL
}

func (e *Evaluator) applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		if len(args) != len(fn.Parameters) {
			return newError("expected %d arg(s) got %d", len(fn.Parameters), len(args))
		}
		extendedEnv := e.extendFunctionEnv(fn, args)
		evaluated := e.Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *Builtin:
		return fn.Fn(e, args...)

	default:
		return newError("not a function: %s", fn.Inspect())
	}

}

func (e *Evaluator) extendFunctionEnv(
	fn *Function,
	args []Object,
) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Let(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, env *Environment) []Object {
	var result []Object

	for _, exp := range exps {
		evaluated := e.Eval(exp, env)

		if isError(evaluated) {
			return []Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) evalProgram(statements []ast.Statement, env *Environment) Object {
	var result Object

	for _, statement := range statements {
		result = e.Eval(statement, env)

		switch result := result.(type) {
		case *Error:
			return result

		case *ReturnValue:
			return result.Value
		}
	}

	return result
}

func (e *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *Environment) Object {
	var result Object

	for _, statement := range block.Statements {
		result = e.Eval(statement, env)

		if isReturn(result) {
			return result
		}
	}

	return result
}

func (e *Evaluator) evalForStatement(node *ast.ForStatement, env *Environment) Object {
	enclosedEnv := NewEnclosedEnvironment(env)

	init := e.Eval(node.Initializer, enclosedEnv)
	if isError(init) {
		return init
	}

	switch body := node.Body.(type) {
	case *ast.BlockStatement: // to optimize for block statements
		bodyEnv := NewEnclosedEnvironment(enclosedEnv)

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

	return NULL
}

func (e *Evaluator) evalWhileStatement(node *ast.WhileStatement, env *Environment) Object {
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

func (e *Evaluator) evalIncludeStatement(node *ast.IncludeStatement, env *Environment) Object {

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

func (e *Evaluator) evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return e.evalBangOperatorExpression(right)
	case "-":
		return e.evalMinusPrefixOperatorExpression(right)

	default:
		return newError("unknown operator: %s%s", operator, right.Inspect())
	}
}

func (e *Evaluator) evalBangOperatorExpression(right Object) Object {
	return boolToBoolObject(!isTruthy(right))
}

func (e *Evaluator) evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Inspect())
	}

	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func (e *Evaluator) evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		leftVal := left.(*Integer).Value
		rightVal := right.(*Integer).Value

		return e.evalIntegerInfixExpression(operator, leftVal, rightVal)

	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		leftVal := left.(*Float).Value
		rightVal := right.(*Float).Value

		return e.evalFloatInfixExpression(operator, leftVal, rightVal)

	case left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ:
		leftVal := left.(*Float).Value
		rightVal := right.(*Integer).Value

		return e.evalFloatInfixExpression(operator, leftVal, float64(rightVal))

	case left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ:
		leftVal := left.(*Integer).Value
		rightVal := right.(*Float).Value

		return e.evalFloatInfixExpression(operator, float64(leftVal), rightVal)

	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
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

func (e *Evaluator) evalIntegerInfixExpression(operator string, left, right int64) Object {
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
		return &Integer{Value: left + right}
	case "-":
		return &Integer{Value: left - right}
	case "*":
		return &Integer{Value: left * right}
	case "/":
		return &Integer{Value: left / right}
	case "%":
		return &Integer{Value: left % right}
	default:
		return newError("unknown operator: %d %s %d",
			left, operator, right)
	}
}

func (e *Evaluator) evalFloatInfixExpression(operator string, left, right float64) Object {
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
		return &Float{Value: left + right}
	case "-":
		return &Float{Value: left - right}
	case "*":
		return &Float{Value: left * right}
	case "/":
		return &Float{Value: left / right}
	case "%":
		return &Float{Value: math.Mod(left, right)}
	default:
		return newError("unknown operator: %f %s %f",
			left, operator, right)
	}
}

func (e *Evaluator) evalStringInfixExpression(operator string, left, right Object) Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s",
			left.Inspect(), operator, right.Inspect())
	}

	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	return &String{Value: leftVal + rightVal}
}

func (e *Evaluator) evalPostfixExpression(operator string, left *ast.Identifier, env *Environment) Object {
	value := e.Eval(left, env)
	if isError(value) {
		return value
	}

	switch value.Type() {
	case INTEGER_OBJ:
		switch operator {
		case "++":
			intObj := value.(*Integer)
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

func (e *Evaluator) evalIfExpression(ie *ast.IfExpression, env *Environment) Object {
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

func (e *Evaluator) evalIdentifier(node *ast.Identifier, env *Environment) Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func (e *Evaluator) evalIndexExpression(node *ast.IndexExpression, env *Environment) Object {
	left := e.Eval(node.Left, env)
	if isError(left) {
		return left
	}

	index := e.Eval(node.Index, env)
	if isError(index) {
		return index
	}

	switch left.Type() {
	case ARRAY_OBJ:
		return e.evalArrayIndexExpression(node, left, index)
	case STRING_OBJ:
		return e.evalStringIndexExpression(node, left, index)
	default:
		return newError("index operator not supported: %s", left.Inspect())
	}
}

func (e *Evaluator) evalArrayIndexExpression(node *ast.IndexExpression, array, index Object) Object {
	arrayObj := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObj.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObj.Value[idx]
}

func (e *Evaluator) evalStringIndexExpression(node *ast.IndexExpression, str, index Object) Object {
	strObj := str.(*String)
	idx := index.(*Integer).Value
	max := int64(len(strObj.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return &String{Value: string([]rune(strObj.Value)[idx])}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func isTruthy(obj Object) bool {
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

func boolToBoolObject(value bool) *Boolean {
	if value {
		return TRUE
	}

	return FALSE
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}

	return false
}

func isReturn(obj Object) bool {
	rt := obj.Type()

	if rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
		return true
	}

	return false
}
