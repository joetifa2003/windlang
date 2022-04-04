package evaluator

import (
	"testing"

	"github.com/joetifa2003/windlang/lexer"
	"github.com/joetifa2003/windlang/parser"
	"github.com/stretchr/testify/assert"
)

const fileName = "main-test.wind"

func TestEvalBooleanInfixExpression(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"3.1 > 3.0", true},
		{"3.1 > 3", true},
		{"3 < 3.1", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(&Boolean{}, evaluated)
		boolVal := evaluated.(*Boolean).Value
		assert.Equal(tc.expected, boolVal)
	}
}

func TestBangOperator(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected Object
	}{
		{"!true", &Boolean{Value: false}},
		{"!false", &Boolean{Value: true}},
		{"!5", &Boolean{Value: false}},
		{"!!true", &Boolean{Value: true}},
		{"!!false", &Boolean{Value: false}},
		{"!!5", &Boolean{Value: true}},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(tc.expected, evaluated)
		assert.Equal(tc.expected, evaluated)
	}
}

func TestIfElseExpressions(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected Object
	}{
		{"if (true) { 10 }", &Integer{Value: 10}},
		{"if (1) { 10 }", &Integer{Value: 10}},
		{"if (1 < 2) { 10 }", &Integer{Value: 10}},
		{"if (1 > 2) { 10 } else { 20 }", &Integer{Value: 20}},
		{"if (1 < 2) { 10 } else { 20 }", &Integer{Value: 10}},
		{"if (false) { 1 }", &Nil{}},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(tc.expected, evaluated)
		assert.Equal(tc.expected, evaluated)
	}
}

func TestLetStatements(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected Object
	}{
		{"let a = 5; a;", &Integer{Value: 5}},
		{"let a = 5 * 5; a;", &Integer{Value: 25}},
		{"let a = 5; let b = a; b;", &Integer{Value: 5}},
		{"let a = 5; let b = a; let c = a + b + 5; c;", &Integer{Value: 15}},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(tc.expected, evaluated)
		assert.Equal(tc.expected, evaluated)
	}
}

func TestFunctions(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected Object
	}{
		{"let double = fn(x) { x * 2; }; double(5);", &Integer{Value: 10}},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", &Integer{Value: 10}},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", &Integer{Value: 20}},
		{"fn(x) { x; }(5)", &Integer{Value: 5}},
		{"fn(x) { return x; }(5)", &Integer{Value: 5}},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(tc.expected, evaluated)
		assert.Equal(tc.expected, evaluated)
	}
}

func TestClosures(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected Object
	}{
		{"let newAdder = fn(x) { fn(y) { x + y }; }; let addTwo = newAdder(2); addTwo(2);", &Integer{Value: 4}},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(tc.expected, evaluated)
		assert.Equal(tc.expected, evaluated)
	}
}

func TestStringLiteral(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected Object
	}{
		{`"Hello World!"`, &String{Value: "Hello World!"}},
		{`"Hello" + " World!"`, &String{Value: "Hello World!"}},
		{`"Hello" + " " + "World!"`, &String{Value: "Hello World!"}},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(tc.expected, evaluated)
		assert.Equal(tc.expected, evaluated)
	}
}

func TestHashLiterals(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected Object
	}{
		{
			`let x = { "foo": 1, "bar": 2 }; x["foo"]`,
			&Integer{Value: 1},
		},
		{
			`let x = { "foo": 1, "bar": 2 }; x.bar`,
			&Integer{Value: 2},
		},
		{
			`let x = { "foo": fn() { return 1; } }; x["foo"]()`,
			&Integer{Value: 1},
		},
		{
			`let x = { "foo": fn() { return 1; } }; x.foo()`,
			&Integer{Value: 1},
		},
		{
			`let x = { "foo": 1 }; x.foo++; x.foo`,
			&Integer{Value: 2},
		},
		{
			`let x = { "foo": 1 }; x["foo"]++; x["foo"]`,
			&Integer{Value: 2},
		},
		{
			`let x = {"foo": { "bar": 1 } }; x.foo.bar`,
			&Integer{Value: 1},
		},
		{
			`let x = {"foo": { "bar": 1 } }; x.foo.bar = 2; x.foo.bar`,
			&Integer{Value: 2},
		},
		{
			`let x = { "foo": { "bar": fn() { return { "baz": 1 }; } } }; x.foo.bar().baz`,
			&Integer{Value: 1},
		},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(tc.expected, evaluated)
		assert.Equal(tc.expected, evaluated)
	}
}

func TestIntInfixExpression(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected int64
	}{
		{"5 + 5;", 10},
		{"5 - 5;", 0},
		{"5 * 5;", 25},
		{"5 / 5;", 1},
		{"4 % 2;", 0},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(&Integer{}, evaluated)
		intVal := evaluated.(*Integer).Value
		assert.Equal(tc.expected, intVal)
	}
}

func TestFloatInfixExpression(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected float64
	}{
		{"5.0 + 5.0;", 10},
		{"5.5 - 5.5;", 0},
		{"5.5 * 5.5;", 30.25},
		{"5.8 / 5.8;", 1},
		{"4.0 % 2.0;", 0},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(&Float{}, evaluated)
		intVal := evaluated.(*Float).Value
		assert.Equal(tc.expected, intVal)
	}
}

func TestFloatInfixIntExpression(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected float64
	}{
		{"5.0 + 5;", 10},
		{"5.5 - 5;", 0.5},
		{"5.5 * 5;", 27.5},
		{"5.8 / 5;", 1.16},
		{"4.0 % 2;", 0},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(&Float{}, evaluated)
		intVal := evaluated.(*Float).Value
		assert.Equal(tc.expected, intVal)
	}
}

func TestIntPrefixOperators(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected int64
	}{
		{"-5;", -5},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(&Integer{}, evaluated)
		intVal := evaluated.(*Integer).Value
		assert.Equal(tc.expected, intVal)
	}
}

func TestFloatPrefixOperators(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected float64
	}{
		{"-5.0;", -5},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(&Float{}, evaluated)
		intVal := evaluated.(*Float).Value
		assert.Equal(tc.expected, intVal)
	}
}

func TestArrayLiteral(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected []Object
	}{
		{
			`[1, 2, 3]`,
			[]Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}},
		},
		{
			`[1, 2.5, 3, true, "Hello"]`,
			[]Object{&Integer{Value: 1}, &Float{Value: 2.5}, &Integer{Value: 3}, &Boolean{Value: true}, &String{Value: "Hello"}},
		},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(&Array{}, evaluated)
		arrayVal := evaluated.(*Array).Value
		assert.Equal(tc.expected, arrayVal)
	}
}

func TestArrayIndex(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected Object
	}{
		{
			`[1, 2, 3][0]`,
			&Integer{Value: 1},
		},
		{
			`[1, 5, true, "Hello", fn() { return "Hello"; }][4]()`,
			&String{Value: "Hello"},
		},
	}

	for _, tc := range tests {
		evaluated, err := testEval(tc.input)
		assert.Nil(err)
		assert.IsType(tc.expected, evaluated)
		assert.Equal(tc.expected, evaluated)
	}
}

func testEval(input string) (Object, *Error) {
	l := lexer.New(input)
	p := parser.New(l, fileName)
	program := p.ParseProgram()

	envManager := NewEnvironmentManager()
	env, _ := envManager.Get(fileName)
	evaluator := New(envManager, fileName)

	return evaluator.Eval(program, env)
}
