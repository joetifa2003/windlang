package evaluator

import (
	"fmt"
	"io/ioutil"
	"math"
	"wind-vm-go/ast"
	"wind-vm-go/lexer"
	"wind-vm-go/parser"
	"wind-vm-go/token"
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
func (e *Evaluator) Eval(node ast.Node, env *Environment) (Object, *Error) {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node.Statements, env)

	case *ast.BlockStatement:
		return e.evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env)

	case *ast.LetStatement:
		return e.evalLetStatement(node, env)

	case *ast.ReturnStatement:
		return e.evalReturnStatement(node, env)

	case *ast.ForStatement:
		return e.evalForStatement(node, env)

	case *ast.WhileStatement:
		return e.evalWhileStatement(node, env)

	case *ast.IncludeStatement:
		return e.evalIncludeStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}, nil

	case *ast.FloatLiteral:
		return &Float{Value: node.Value}, nil

	case *ast.Boolean:
		return boolToBoolObject(node.Value), nil

	case *ast.PrefixExpression:
		return e.evalPrefixExpression(node, env)

	case *ast.InfixExpression:
		return e.evalInfixExpression(node, env)

	case *ast.PostfixExpression:
		return e.evalPostfixExpression(node, env)

	case *ast.IfExpression:
		return e.evalIfExpression(node, env)

	case *ast.Identifier:
		return e.evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		return &Function{Parameters: node.Parameters, Body: node.Body, Env: env}, nil

	case *ast.CallExpression:
		return e.evalCallExpression(node, env)

	case *ast.StringLiteral:
		return &String{Value: node.Value}, nil

	case *ast.AssignExpression:
		return e.evalAssignExpression(node, env)

	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node, env)

	case *ast.IndexExpression:
		return e.evalIndexExpression(node, env)

	case *ast.NilLiteral:
		return NIL, nil

	case *ast.HashLiteral:
		return e.evalHashLiteral(node, env)
	}

	return NIL, nil
}

func (e *Evaluator) evalCallExpression(node *ast.CallExpression, env *Environment) (Object, *Error) {
	function, err := e.Eval(node.Function, env)
	if err != nil {
		return nil, err
	}

	args, err := e.evalExpressions(node.Arguments, env)
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
		evaluated, err := e.Eval(fn.Body, extendedEnv)
		if err != nil {
			return nil, err
		}

		return unwrapReturnValue(evaluated), nil

	case *Builtin:
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

func (e *Evaluator) evalLetStatement(node *ast.LetStatement, env *Environment) (Object, *Error) {
	val, err := e.Eval(node.Value, env)
	if err != nil {
		return nil, err
	}

	env.Let(node.Name.Value, val)

	return NIL, nil
}

func (e *Evaluator) evalReturnStatement(node *ast.ReturnStatement, env *Environment) (Object, *Error) {
	val, err := e.Eval(node.ReturnValue, env)
	if err != nil {
		return nil, err
	}

	return &ReturnValue{Value: val}, nil
}

func (e *Evaluator) evalExpressions(exps []ast.Expression, env *Environment) ([]Object, *Error) {
	var result []Object

	for _, exp := range exps {
		evaluated, err := e.Eval(exp, env)
		if err != nil {
			return nil, err
		}

		result = append(result, evaluated)
	}

	return result, nil
}

func (e *Evaluator) evalProgram(statements []ast.Statement, env *Environment) (Object, *Error) {
	var result Object
	var err *Error

	for _, statement := range statements {
		result, err = e.Eval(statement, env)
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

func (e *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *Environment) (Object, *Error) {
	enclosedEnv := NewEnclosedEnvironment(env)

	var result Object
	var err *Error

	for _, statement := range block.Statements {
		result, err = e.Eval(statement, enclosedEnv)
		if err != nil {
			return nil, err
		}

		if isReturn(result) {
			return result, nil
		}
	}

	return result, nil
}

func (e *Evaluator) evalForStatement(node *ast.ForStatement, env *Environment) (Object, *Error) {
	enclosedEnv := NewEnclosedEnvironment(env)

	_, err := e.Eval(node.Initializer, enclosedEnv)
	if err != nil {
		return nil, err
	}

	switch body := node.Body.(type) {
	case *ast.BlockStatement: // to optimize for block statements
		bodyEnv := NewEnclosedEnvironment(enclosedEnv)

		for {
			condition, err := e.Eval(node.Condition, enclosedEnv)
			if err != nil {
				return nil, err
			}

			if !isTruthy(condition) {
				break
			}

			result, err := e.evalBlockStatement(body, bodyEnv)
			if err != nil {
				return nil, err
			}

			if isReturn(result) {
				return result, nil
			}

			bodyEnv.ClearStore()

			_, err = e.Eval(node.Increment, enclosedEnv)
			if err != nil {
				return nil, err
			}
		}

	default:
		for {
			condition, err := e.Eval(node.Condition, enclosedEnv)
			if err != nil {
				return nil, err
			}

			if !isTruthy(condition) {
				break
			}

			result, err := e.Eval(body, enclosedEnv)
			if err != nil {
				return nil, err
			}

			if isReturn(result) {
				return result, nil
			}

			_, err = e.Eval(node.Increment, enclosedEnv)
			if err != nil {
				return nil, err
			}
		}
	}

	return NIL, nil
}

func (e *Evaluator) evalWhileStatement(node *ast.WhileStatement, env *Environment) (Object, *Error) {
	for {
		condition, err := e.Eval(node.Condition, env)
		if err != nil {
			return nil, err
		}

		if !isTruthy(condition) {
			break
		}

		result, err := e.Eval(node.Body, env)
		if err != nil {
			return nil, err
		}

		if isReturn(result) {
			return result, nil
		}
	}

	return NIL, nil
}

func (e *Evaluator) evalIncludeStatement(node *ast.IncludeStatement, env *Environment) (Object, *Error) {
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

		_, err := e.Eval(program, fileEnv)
		if err != nil {
			return nil, err
		}
	}

	if node.Alias != nil {
		hashMap := HashMapFromEnv(fileEnv)
		includeEnv := NewEnvironment()
		includeEnv.Let(node.Alias.Value, hashMap)

		env.Includes = append(env.Includes, includeEnv)
	} else {
		env.Includes = append(env.Includes, fileEnv)
	}

	return NIL, nil
}

func (e *Evaluator) evalPrefixExpression(node *ast.PrefixExpression, env *Environment) (Object, *Error) {
	right, err := e.Eval(node.Right, env)
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
	case *Integer:
		return &Integer{Value: -right.Value}, nil
	case *Float:
		return &Float{Value: -right.Value}, nil
	default:
		return nil, e.newError(node.Token, "unknown operator: -%s", right.Inspect())
	}
}

func (e *Evaluator) evalInfixExpression(node *ast.InfixExpression, env *Environment) (Object, *Error) {
	left, err := e.Eval(node.Left, env)
	if err != nil {
		return nil, err
	}

	right, err := e.Eval(node.Right, env)
	if err != nil {
		return nil, err
	}

	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		leftVal := left.(*Integer).Value
		rightVal := right.(*Integer).Value

		return e.evalIntegerInfixExpression(node, node.Operator, leftVal, rightVal)

	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		leftVal := left.(*Float).Value
		rightVal := right.(*Float).Value

		return e.evalFloatInfixExpression(node, node.Operator, leftVal, rightVal)

	case left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ:
		leftVal := left.(*Float).Value
		rightVal := right.(*Integer).Value

		return e.evalFloatInfixExpression(node, node.Operator, leftVal, float64(rightVal))

	case left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ:
		leftVal := left.(*Integer).Value
		rightVal := right.(*Float).Value

		return e.evalFloatInfixExpression(node, node.Operator, float64(leftVal), rightVal)

	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
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

func (e *Evaluator) evalIntegerInfixExpression(node *ast.InfixExpression, operator string, left, right int64) (Object, *Error) {
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
		return &Integer{Value: left + right}, nil
	case "-":
		return &Integer{Value: left - right}, nil
	case "*":
		return &Integer{Value: left * right}, nil
	case "/":
		return &Integer{Value: left / right}, nil
	case "%":
		return &Integer{Value: left % right}, nil
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

func (e *Evaluator) evalIfExpression(ie *ast.IfExpression, env *Environment) (Object, *Error) {
	condition, err := e.Eval(ie.Condition, env)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		return e.Eval(ie.ThenBranch, env)
	} else if ie.ElseBranch != nil {
		return e.Eval(ie.ElseBranch, env)
	} else {
		return NIL, nil
	}
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier, env *Environment) (Object, *Error) {
	if val, ok := env.Get(node.Value); ok {
		return val, nil
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin, nil
	}

	return nil, e.newError(node.Token, "identifier not found: "+node.Value)
}

func (e *Evaluator) evalArrayLiteral(node *ast.ArrayLiteral, env *Environment) (Object, *Error) {
	objects := make([]Object, len(node.Value))

	for index, expr := range node.Value {
		object, err := e.Eval(expr, env)
		if err != nil {
			return nil, err
		}

		objects[index] = object
	}

	return &Array{Value: objects}, nil
}

func (e *Evaluator) evalIndexExpression(node *ast.IndexExpression, env *Environment) (Object, *Error) {
	left, err := e.Eval(node.Left, env)
	if err != nil {
		return nil, err
	}

	index, err := e.Eval(node.Index, env)
	if err != nil {
		return nil, err
	}

	switch left.Type() {
	case ARRAY_OBJ:
		return e.evalArrayIndexExpression(node, left, index)
	case STRING_OBJ:
		return e.evalStringIndexExpression(node, left, index)
	case HASH_OBJ:
		return e.evalHashIndexExpression(node, left, index)
	default:
		return nil, e.newError(node.Token, "index operator not supported: %s", left.Inspect())
	}
}

func (e *Evaluator) evalHashLiteral(node *ast.HashLiteral, env *Environment) (Object, *Error) {
	hash := make(map[HashKey]Object)

	for key, value := range node.Pairs {
		hashKey, err := e.Eval(key, env)
		if err != nil {
			return nil, err
		}

		key, ok := hashKey.(Hashable)
		if !ok {
			return nil, e.newError(node.Token, "unusable as hash key: %s", hashKey.Inspect())
		}

		hashValue, err := e.Eval(value, env)
		if err != nil {
			return nil, err
		}

		hash[key.HashKey()] = hashValue
	}

	return &Hash{Pairs: hash}, nil
}

func (e *Evaluator) evalArrayIndexExpression(node *ast.IndexExpression, array, index Object) (Object, *Error) {
	arrayObj := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObj.Value) - 1)

	if idx < 0 || idx > max {
		return NIL, nil
	}

	return arrayObj.Value[idx], nil
}

func (e *Evaluator) evalStringIndexExpression(node *ast.IndexExpression, str, index Object) (Object, *Error) {
	strObj := str.(*String)
	idx := index.(*Integer).Value
	max := int64(len(strObj.Value) - 1)

	if idx < 0 || idx > max {
		return NIL, nil
	}

	return &String{Value: string([]rune(strObj.Value)[idx])}, nil
}

func (e *Evaluator) evalHashIndexExpression(node *ast.IndexExpression, hash, index Object) (Object, *Error) {
	hashObj := hash.(*Hash)
	key, ok := index.(Hashable)
	if !ok {
		return nil, e.newError(node.Token, "unusable as hash key: %s", index.Inspect())
	}

	if val, ok := hashObj.Pairs[key.HashKey()]; ok {
		return val, nil
	}

	return NIL, nil
}

func (e *Evaluator) evalAssignExpression(node *ast.AssignExpression, env *Environment) (Object, *Error) {
	val, err := e.Eval(node.Value, env)
	if err != nil {
		return nil, err
	}

	switch left := node.Name.(type) {
	case *ast.Identifier:
		return e.evalAssingIdentifierExpression(node, left, val, env)

	case *ast.IndexExpression:
		return e.evalAssingIndexExpression(node, left, val, env)
	}

	return val, nil
}

func (e *Evaluator) evalAssingIdentifierExpression(node *ast.AssignExpression, left *ast.Identifier, val Object, env *Environment) (Object, *Error) {
	_, ok := env.Set(left.Value, val)
	if !ok {
		return nil, e.newError(node.Token, "identifier not found: "+left.Value)
	}

	return val, nil
}

func (e *Evaluator) evalAssingIndexExpression(node *ast.AssignExpression, left *ast.IndexExpression, val Object, env *Environment) (Object, *Error) {
	leftObj, err := e.Eval(left.Left, env)
	if err != nil {
		return nil, err
	}

	index, err := e.Eval(left.Index, env)
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
	idx := index.(*Integer).Value
	max := int64(len(leftObj.Value) - 1)

	if idx < 0 || idx > max {
		return NIL, nil
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

func (e *Evaluator) evalPostfixExpression(node *ast.PostfixExpression, env *Environment) (Object, *Error) {
	left, err := e.Eval(node.Left, env)
	if err != nil {
		return nil, err
	}

	switch left := left.(type) {
	case *Integer:
		return e.evalPostfixIntegerExpression(node, node.Operator, left)
	default:
		return nil, e.newError(node.Token, "postfix operator not supported: %s", left.Inspect())
	}
}

func (e *Evaluator) evalPostfixIntegerExpression(node *ast.PostfixExpression, operator string, left *Integer) (Object, *Error) {
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

	if rt == RETURN_VALUE_OBJ {
		return true
	}

	return false
}
