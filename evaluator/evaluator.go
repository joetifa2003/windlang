package evaluator

import (
	"fmt"
	"io/ioutil"
	"math"

	"github.com/joetifa2003/windlang/ast"
	"github.com/joetifa2003/windlang/lexer"
	"github.com/joetifa2003/windlang/parser"
	"github.com/joetifa2003/windlang/token"
)

var (
	NIL   = &Nil{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Evaluator struct {
	envManager *EnvironmentManager
	filePath   string
}

func New(envManager *EnvironmentManager, filePath string) *Evaluator {
	return &Evaluator{
		envManager: envManager,
		filePath:   filePath,
	}
}

// Eval returns the result of the evaluation and potential error
func (e *Evaluator) Eval(node ast.Node, env *Environment, this Object) (Object, *Error) {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node.Statements, env, this)

	case *ast.BlockStatement:
		return e.evalBlockStatement(node, env, this)

	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env, this)

	case *ast.LetStatement:
		return e.evalLetStatement(node, env, this)

	case *ast.ReturnStatement:
		return e.evalReturnStatement(node, env, this)

	case *ast.ForStatement:
		return e.evalForStatement(node, env, this)

	case *ast.WhileStatement:
		return e.evalWhileStatement(node, env, this)

	case *ast.IncludeStatement:
		return e.evalIncludeStatement(node, env, this)

	// Expressions
	case *ast.IntegerLiteral:
		return Integer{Value: node.Value}, nil

	case *ast.FloatLiteral:
		return &Float{Value: node.Value}, nil

	case *ast.Boolean:
		return boolToBoolObject(node.Value), nil

	case *ast.PrefixExpression:
		return e.evalPrefixExpression(node, env, this)

	case *ast.InfixExpression:
		return e.evalInfixExpression(node, env, this)

	case *ast.PostfixExpression:
		return e.evalPostfixExpression(node, env, this)

	case *ast.IfExpression:
		return e.evalIfExpression(node, env, this)

	case *ast.Identifier:
		return e.evalIdentifier(node, env, this)

	case *ast.FunctionLiteral:
		return &Function{Parameters: node.Parameters, Body: node.Body, Env: env, This: this}, nil

	case *ast.CallExpression:
		return e.evalCallExpression(node, env, this)

	case *ast.StringLiteral:
		return &String{Value: node.Value}, nil

	case *ast.AssignExpression:
		return e.evalAssignExpression(node, env, this)

	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node, env, this)

	case *ast.IndexExpression:
		return e.evalIndexExpression(node, env, this)

	case *ast.NilLiteral:
		return NIL, nil

	case *ast.HashLiteral:
		return e.evalHashLiteral(node, env)

	case *ast.EchoStatement:
		val, err := e.Eval(node.Value, env, this)
		if err != nil {
			return nil, err
		}

		fmt.Println(val.Inspect())
	}

	return NIL, nil
}

func (e *Evaluator) evalCallExpression(node *ast.CallExpression, env *Environment, this Object) (Object, *Error) {
	function, err := e.Eval(node.Function, env, this)
	if err != nil {
		return nil, err
	}

	args, err := e.evalExpressions(node.Arguments, env, this)
	if err != nil {
		return nil, err
	}

	return e.applyFunction(node, function, args)
}

func (e *Evaluator) applyFunction(node *ast.CallExpression, fn Object, args []Object) (Object, *Error) {
	switch fn := fn.(type) {
	case *Function:
		if len(args) != len(fn.Parameters) {
			return nil, e.newError(node.Token, "expected %d arg(s) got %d", len(fn.Parameters), len(args))
		}

		extendedEnv := e.extendFunctionEnv(fn, args)
		evaluated, err := e.Eval(fn.Body, extendedEnv, fn.This)
		if err != nil {
			return nil, err
		}

		return unwrapReturnValue(evaluated), nil

	case *GoFunction:
		if fn.ArgsCount != -1 && len(args) != fn.ArgsCount {
			return nil, e.newError(node.Token, "expected %d arg(s) got %d", fn.ArgsCount, len(args))
		}

		for i, t := range fn.ArgsTypes {
			if t != Any && t != args[i].Type() {
				return nil, e.newError(node.Token, "expected arg %d to be of type %s got %s", i, t, args[i].Type())
			}
		}

		return fn.Fn(e, node, args...)

	default:
		return nil, e.newError(node.Token, "not a function: %s", fn.Inspect())
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

func (e *Evaluator) evalLetStatement(node *ast.LetStatement, env *Environment, this Object) (Object, *Error) {
	val, err := e.Eval(node.Value, env, this)
	if err != nil {
		return nil, err
	}

	if node.Constant {
		env.LetConstant(node.Name.Value, val)
	} else {
		env.Let(node.Name.Value, val)
	}

	return NIL, nil
}

func (e *Evaluator) evalReturnStatement(node *ast.ReturnStatement, env *Environment, this Object) (Object, *Error) {
	val, err := e.Eval(node.ReturnValue, env, this)
	if err != nil {
		return nil, err
	}

	return &ReturnValue{Value: val}, nil
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, env *Environment, this Object) ([]Object, *Error) {
	var result []Object

	for _, exp := range exps {
		evaluated, err := e.Eval(exp, env, this)
		if err != nil {
			return nil, err
		}

		result = append(result, evaluated)
	}

	return result, nil
}

func (e *Evaluator) evalProgram(statements []ast.Statement, env *Environment, this Object) (Object, *Error) {
	var result Object
	var err *Error

	for _, statement := range statements {
		result, err = e.Eval(statement, env, this)
		if err != nil {
			return nil, err
		}

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value, nil
		}
	}

	return result, nil
}

func (e *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *Environment, this Object) (Object, *Error) {
	enclosedEnv := NewEnclosedEnvironment(env)

	var result Object
	var err *Error

	for _, statement := range block.Statements {
		result, err = e.Eval(statement, enclosedEnv, this)
		if err != nil {
			return nil, err
		}

		if isReturn(result) {
			return result, nil
		}
	}

	return result, nil
}

func (e *Evaluator) evalForStatement(node *ast.ForStatement, env *Environment, this Object) (Object, *Error) {
	enclosedEnv := NewEnclosedEnvironment(env)

	_, err := e.Eval(node.Initializer, enclosedEnv, this)
	if err != nil {
		return nil, err
	}

	switch body := node.Body.(type) {
	case *ast.BlockStatement: // to optimize for block statements
		bodyEnv := NewEnclosedEnvironment(enclosedEnv)

		for {
			condition, err := e.Eval(node.Condition, enclosedEnv, this)
			if err != nil {
				return nil, err
			}

			if !isTruthy(condition) {
				break
			}

			result, err := e.evalBlockStatement(body, bodyEnv, this)
			if err != nil {
				return nil, err
			}

			if isReturn(result) {
				return result, nil
			}

			bodyEnv.ClearStore()

			_, err = e.Eval(node.Increment, enclosedEnv, this)
			if err != nil {
				return nil, err
			}
		}

	default:
		for {
			condition, err := e.Eval(node.Condition, enclosedEnv, this)
			if err != nil {
				return nil, err
			}

			if !isTruthy(condition) {
				break
			}

			result, err := e.Eval(body, enclosedEnv, this)
			if err != nil {
				return nil, err
			}

			if isReturn(result) {
				return result, nil
			}

			_, err = e.Eval(node.Increment, enclosedEnv, this)
			if err != nil {
				return nil, err
			}
		}
	}

	return NIL, nil
}

func (e *Evaluator) evalWhileStatement(node *ast.WhileStatement, env *Environment, this Object) (Object, *Error) {
	for {
		condition, err := e.Eval(node.Condition, env, this)
		if err != nil {
			return nil, err
		}

		if !isTruthy(condition) {
			break
		}

		result, err := e.Eval(node.Body, env, this)
		if err != nil {
			return nil, err
		}

		if isReturn(result) {
			return result, nil
		}
	}

	return NIL, nil
}

func (e *Evaluator) evalIncludeStatement(node *ast.IncludeStatement, env *Environment, this Object) (Object, *Error) {
	path := node.Path
	fileEnv, evaluated := e.envManager.Get(path)

	if !evaluated {
		file, ioErr := ioutil.ReadFile(path)
		if ioErr != nil {
			return nil, e.newError(node.Token, "cannot read file: %s", path)
		}

		input := string(file)
		lexer := lexer.New(input)
		parser := parser.New(lexer, path)
		program := parser.ParseProgram()
		parser.ReportErrors()

		_, err := e.Eval(program, fileEnv, this)
		if err != nil {
			return nil, err
		}
	}

	if node.Alias != nil {
		includeObject := &IncludeObject{
			Value: fileEnv,
		}

		env.AddAlias(node.Alias.Value, includeObject)
	} else {
		env.Includes = append(env.Includes, fileEnv)
	}

	return NIL, nil
}

func (e *Evaluator) evalPrefixExpression(node *ast.PrefixExpression, env *Environment, this Object) (Object, *Error) {
	right, err := e.Eval(node.Right, env, this)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case "!":
		return e.evalBangOperatorExpression(node, right)
	case "-":
		return e.evalMinusPrefixOperatorExpression(node, right)

	default:
		return nil, e.newError(node.Token, "unknown operator: %s%s", node.Operator, right.Inspect())
	}
}

func (e *Evaluator) evalBangOperatorExpression(node *ast.PrefixExpression, right Object) (Object, *Error) {
	return boolToBoolObject(!isTruthy(right)), nil
}

func (e *Evaluator) evalMinusPrefixOperatorExpression(node *ast.PrefixExpression, right Object) (Object, *Error) {
	switch right := right.(type) {
	case Integer:
		return Integer{Value: -right.Value}, nil
	case *Float:
		return &Float{Value: -right.Value}, nil
	default:
		return nil, e.newError(node.Token, "unknown operator: -%s", right.Inspect())
	}
}

func (e *Evaluator) evalInfixExpression(node *ast.InfixExpression, env *Environment, this Object) (Object, *Error) {
	left, err := e.Eval(node.Left, env, this)
	if err != nil {
		return nil, err
	}

	right, err := e.Eval(node.Right, env, this)
	if err != nil {
		return nil, err
	}

	switch {
	case left.Type() == IntegerObj && right.Type() == IntegerObj:
		leftVal := left.(Integer).Value
		rightVal := right.(Integer).Value

		return e.evalIntegerInfixExpression(node, node.Operator, leftVal, rightVal)

	case left.Type() == FloatObj && right.Type() == FloatObj:
		leftVal := left.(*Float).Value
		rightVal := right.(*Float).Value

		return e.evalFloatInfixExpression(node, node.Operator, leftVal, rightVal)

	case left.Type() == FloatObj && right.Type() == IntegerObj:
		leftVal := left.(*Float).Value
		rightVal := right.(Integer).Value

		return e.evalFloatInfixExpression(node, node.Operator, leftVal, float64(rightVal))

	case left.Type() == IntegerObj && right.Type() == FloatObj:
		leftVal := left.(Integer).Value
		rightVal := right.(*Float).Value

		return e.evalFloatInfixExpression(node, node.Operator, float64(leftVal), rightVal)

	case left.Type() == StringObj && right.Type() == StringObj:
		return e.evalStringInfixExpression(node, node.Operator, left, right)

	case node.Operator == "==":
		return boolToBoolObject(left == right), nil

	case node.Operator == "!=":
		return boolToBoolObject(left != right), nil

	case node.Operator == "&&":
		return boolToBoolObject(isTruthy(left) && isTruthy(right)), nil

	case node.Operator == "||":
		return boolToBoolObject(isTruthy(left) || isTruthy(right)), nil

	default:
		return nil, e.newError(node.Token, "unknown operator: %s %s %s",
			left.Inspect(), node.Operator, right.Inspect())
	}
}

func (e *Evaluator) evalIntegerInfixExpression(node *ast.InfixExpression, operator string, left, right int) (Object, *Error) {
	switch operator {
	case "<":
		return boolToBoolObject(left < right), nil
	case "<=":
		return boolToBoolObject(left <= right), nil
	case ">":
		return boolToBoolObject(left > right), nil
	case ">=":
		return boolToBoolObject(left >= right), nil
	case "==":
		return boolToBoolObject(left == right), nil
	case "!=":
		return boolToBoolObject(left != right), nil
	case "+":
		return Integer{Value: left + right}, nil
	case "-":
		return Integer{Value: left - right}, nil
	case "*":
		return Integer{Value: left * right}, nil
	case "/":
		return Integer{Value: left / right}, nil
	case "%":
		return Integer{Value: left % right}, nil
	default:
		return nil, e.newError(node.Token, "unknown operator: %d %s %d",
			left, operator, right)
	}
}

func (e *Evaluator) evalFloatInfixExpression(node *ast.InfixExpression, operator string, left, right float64) (Object, *Error) {
	switch operator {
	case "<":
		return boolToBoolObject(left < right), nil
	case "<=":
		return boolToBoolObject(left <= right), nil
	case ">":
		return boolToBoolObject(left > right), nil
	case ">=":
		return boolToBoolObject(left >= right), nil
	case "==":
		return boolToBoolObject(left == right), nil
	case "!=":
		return boolToBoolObject(left != right), nil
	case "+":
		return &Float{Value: left + right}, nil
	case "-":
		return &Float{Value: left - right}, nil
	case "*":
		return &Float{Value: left * right}, nil
	case "/":
		return &Float{Value: left / right}, nil
	case "%":
		return &Float{Value: math.Mod(left, right)}, nil
	default:
		return nil, e.newError(node.Token, "unknown operator: %f %s %f",
			left, operator, right)
	}
}

func (e *Evaluator) evalStringInfixExpression(node *ast.InfixExpression, operator string, left, right Object) (Object, *Error) {
	if operator != "+" {
		return nil, e.newError(node.Token, "unknown operator: %s %s %s",
			left.Inspect(), operator, right.Inspect())
	}

	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	return &String{Value: leftVal + rightVal}, nil
}

func (e *Evaluator) evalIfExpression(ie *ast.IfExpression, env *Environment, this Object) (Object, *Error) {
	condition, err := e.Eval(ie.Condition, env, this)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		return e.Eval(ie.ThenBranch, env, this)
	} else if ie.ElseBranch != nil {
		return e.Eval(ie.ElseBranch, env, this)
	} else {
		return NIL, nil
	}
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier, env *Environment, this Object) (Object, *Error) {
	if val, ok := env.Get(node.Value); ok {
		return val, nil
	}

	if node.Value == "this" {
		return this, nil
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin, nil
	}

	return nil, e.newError(node.Token, "identifier not found: "+node.Value)
}

func (e *Evaluator) evalArrayLiteral(node *ast.ArrayLiteral, env *Environment, this Object) (Object, *Error) {
	objects := make([]Object, len(node.Value))

	for index, expr := range node.Value {
		object, err := e.Eval(expr, env, this)
		if err != nil {
			return nil, err
		}

		objects[index] = object
	}

	return &Array{Value: objects}, nil
}

func (e *Evaluator) evalIndexExpression(node *ast.IndexExpression, env *Environment, this Object) (Object, *Error) {
	left, err := e.Eval(node.Left, env, this)
	if err != nil {
		return nil, err
	}

	index, err := e.Eval(node.Index, env, this)
	if err != nil {
		return nil, err
	}

	switch left := left.(type) {
	case *Array:
		switch index.Type() {
		case IntegerObj:
			return e.evalArrayIndexExpression(node, left, index)
		default:
			return e.evalWithFunctionsIndexExpression(node, left, index)
		}

	case *Hash:
		return e.evalHashIndexExpression(node, left, index)

	case *IncludeObject:
		return e.evalIncludeIndexExpression(node, left, index)

	case ObjectWithFunctions:
		return e.evalWithFunctionsIndexExpression(node, left, index)

	default:
		return nil, e.newError(node.Token, "index operator not supported: %s", left.Inspect())
	}
}

func (e *Evaluator) evalHashLiteral(node *ast.HashLiteral, env *Environment) (Object, *Error) {
	hash := &Hash{Pairs: make(map[HashKey]Object)}

	for key, value := range node.Pairs {
		hashKey, err := e.Eval(key, env, hash)
		if err != nil {
			return nil, err
		}

		key, ok := hashKey.(Hashable)
		if !ok {
			return nil, e.newError(node.Token, "unusable as hash key: %s", hashKey.Inspect())
		}

		hashValue, err := e.Eval(value, env, hash)
		if err != nil {
			return nil, err
		}

		hash.Pairs[key.HashKey()] = hashValue
	}

	return hash, nil
}

func (e *Evaluator) evalArrayIndexExpression(node *ast.IndexExpression, array *Array, index Object) (Object, *Error) {
	idx := index.(Integer).Value
	max := len(array.Value) - 1

	if idx < 0 || idx > max {
		return NIL, nil
	}

	return array.Value[idx], nil
}

func (e *Evaluator) evalHashIndexExpression(node *ast.IndexExpression, hash *Hash, index Object) (Object, *Error) {
	key, ok := index.(Hashable)
	if !ok {
		return nil, e.newError(node.Token, "unusable as hash key: %s", index.Inspect())
	}

	if val, ok := hash.Pairs[key.HashKey()]; ok {
		return val, nil
	}

	return NIL, nil
}

func (e *Evaluator) evalIncludeIndexExpression(node *ast.IndexExpression, include, index Object) (Object, *Error) {
	includeObj := include.(*IncludeObject)
	key, ok := index.(*String)
	if !ok {
		return nil, e.newError(node.Token, "unusable as include key: %s", key.Inspect())
	}

	obj, ok := includeObj.Value.Store[key.Value]
	if !ok {
		return nil, e.newError(node.Token, "include key not found: %s", key.Inspect())
	}

	return obj, nil
}

func (e *Evaluator) evalWithFunctionsIndexExpression(node *ast.IndexExpression, obj ObjectWithFunctions, index Object) (Object, *Error) {
	name, ok := index.(*String)
	if !ok {
		return nil, e.newError(node.Token, "cannot use %s as an index", index.Type().String())
	}

	fn, ok := obj.GetFunction(name.Value)
	if !ok {
		return nil, e.newError(
			node.Token, "cannot find '%s' function on type %s",
			name.Value,
			obj.(Object).Type().String(),
		)
	}

	return fn, nil
}

func (e *Evaluator) evalAssignExpression(node *ast.AssignExpression, env *Environment, this Object) (Object, *Error) {
	val, err := e.Eval(node.Value, env, this)
	if err != nil {
		return nil, err
	}

	switch left := node.Name.(type) {
	case *ast.Identifier:
		return e.evalAssingIdentifierExpression(node, left, val, env)

	case *ast.IndexExpression:
		return e.evalAssingIndexExpression(node, left, val, env, this)
	}

	return nil, e.newError(node.Token, "cannot assign to %s", node.Name.String())
}

func (e *Evaluator) evalAssingIdentifierExpression(node *ast.AssignExpression, left *ast.Identifier, val Object, env *Environment) (Object, *Error) {
	if env.IsConstant(left.Value) {
		return nil, e.newError(node.Token, "cannot assign to a constant variable %s", left.Value)
	}

	_, ok := env.Set(left.Value, val)
	if !ok {
		return nil, e.newError(node.Token, "identifier not found: "+left.Value)
	}

	return val, nil
}

func (e *Evaluator) evalAssingIndexExpression(node *ast.AssignExpression, left *ast.IndexExpression, val Object, env *Environment, this Object) (Object, *Error) {
	leftObj, err := e.Eval(left.Left, env, this)
	if err != nil {
		return nil, err
	}

	index, err := e.Eval(left.Index, env, this)
	if err != nil {
		return nil, err
	}

	switch leftObj := leftObj.(type) {
	case *Array:
		return e.evalAssingArrayIndexExpression(node, leftObj, index, val)
	case *Hash:
		return e.evalAssingHashIndexExpression(node, leftObj, index, val)
	default:
		return nil, e.newError(node.Token, "index operator not supported: %s", leftObj.Inspect())
	}
}

func (e *Evaluator) evalAssingArrayIndexExpression(node *ast.AssignExpression, leftObj *Array, index Object, val Object) (Object, *Error) {
	idx := index.(Integer).Value
	max := len(leftObj.Value) - 1

	if idx < 0 || idx > max {
		return nil, e.newError(node.Token, "index out of bounds")
	}

	leftObj.Value[idx] = val

	return val, nil
}

func (e *Evaluator) evalAssingHashIndexExpression(node *ast.AssignExpression, leftObj *Hash, index Object, val Object) (Object, *Error) {
	key, ok := index.(Hashable)
	if !ok {
		return nil, e.newError(node.Token, "unusable as hash key: %s", index.Inspect())
	}

	leftObj.Pairs[key.HashKey()] = val

	return val, nil
}

func (e *Evaluator) evalPostfixExpression(node *ast.PostfixExpression, env *Environment, this Object) (Object, *Error) {
	left, err := e.Eval(node.Left, env, this)
	if err != nil {
		return nil, err
	}

	switch left := left.(type) {
	case Integer:
		return e.evalPostfixIntegerExpression(node, node.Operator, left)
	default:
		return nil, e.newError(node.Token, "postfix operator not supported: %s", left.Inspect())
	}
}

func (e *Evaluator) evalPostfixIntegerExpression(node *ast.PostfixExpression, operator string, left Integer) (Object, *Error) {
	switch operator {
	case "++":
		left.Value++
		return left, nil
	case "--":
		left.Value--
		return left, nil
	default:
		return nil, e.newError(node.Token, "postfix operator not supported: %s", operator)
	}
}

func (e *Evaluator) newError(token token.Token, format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf("[file %s:%d] %s", e.filePath, token.Line, fmt.Sprintf(format, a...))}
}

func isTruthy(obj Object) bool {
	switch obj {
	case NIL:
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

func isReturn(obj Object) bool {
	rt := obj.Type()

	return rt == ReturnValueObj
}
