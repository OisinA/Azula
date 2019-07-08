package tests

import (
	"azula/lexer"
	"azula/object"
	"azula/parser"
	"azula/evaluator"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-5", -5},
		{"5 + 5 + 5 + 5 + 5", 25},
		{"(5 + 10) * 2 + 4", 34},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 == 1", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{"if(true) { 10 }", 10},
		{"if(1 > 2) { 10 } else { 20 }", 20},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestEvalReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 20;", 10},
		{"return 2 * 5; 5;", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		expectedMessage string
	}{
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestEvalLetStatements(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"int x = 5; x;", 5},
		{"int x = 5 * 25; x;", 125},
		{"int x = 5; x = 10; x;", 10},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEvalErrorStatements(t *testing.T) {
	tests := []struct {
		input string
		expected string
	}{
		{"error x = _\"error\"_; x;", "error"},
	}
	for _, tt := range tests {
		eval := testEval(tt.input)
		result, ok := eval.(*object.Error)
		if !ok {
			t.Errorf("object is not error. got=%T (%+v)", eval, eval)
		}
		if result.Message != tt.expected {
			t.Errorf("error message incorrect. got=%q", result.Message)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func function(int x): array(int) { [1, 2, 3, 4]; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Paramaters=%+v", fn.Parameters)
	}

	expectedBody := "[1, 2, 3, 4]"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"func identity(int x): int { return x; }; identity(5);", 5},
		{"func double(int x): int { return x * 2; }; double(5);", 10},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!" + 17`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!17" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T {%+v}", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{
			"[1, 2, 3, 4][0]",
			1,
		},
		{
			`array(int) x = [400, 1000]; x[0];`,
			400,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestForLoopExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{
			"int i = 0; for(x in [1, 2, 3, 4]) { x; }",
			4,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, int64(tt.expected.(int)))
	}
}

func TestClassLiterals(t *testing.T) {
	input := "class TestClass(int x) { x = 10; }"

	evaluated := testEval(input)
	class, ok := evaluated.(*object.Class)
	if !ok {
		t.Fatalf("object is not Class. got=%T (%+v)", evaluated, evaluated)
	}
	if len(class.Parameters) != 1 {
		t.Fatalf("class has wrong parameters. Parameters=%+v", class.Parameters)
	}
	if class.Body.String() != "x = 10;" {
		t.Fatalf("body is not %q. got=%q", "x = 10", class.Body.String())
	}
}

func TestClassCalls(t *testing.T) {
	input := `
	class TestClass(int x) {
		func getx(): int {
			return x;
		}
	}
	TestClass c = TestClass(5);
	c.getx();
	`

	evaluated := testEval(input)
	_, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", evaluated, evaluated)
	}
}

func TestHashLiterals(t *testing.T) {
	input := `
	string two = "two";
	{
		"one"=>10-9,
		two=>1 + 1,
		"thr" + "ee"=>6 / 2,
		4=>4,
	};
	`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey(): 1,
		(&object.String{Value: "two"}).HashKey(): 2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey(): 4,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{
			`{"foo"=>5}["foo"]`,
			5,
		},
		{
			`{"foo"=>5}["bar"]`,
			nil,
		},
		{
			`string key = "foo"; {"foo"=>5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5=>5}[5]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			_, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("evaluated not error")
			}
		}
	}
}