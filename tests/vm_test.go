package tests

import (
	"fmt"
	"testing"

	"azula/compiler"
	"azula/object"
	"azula/vm"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		vm := vm.New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testCompileIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testVMBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case string:
		err := testVMStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object not array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
			return
		}

		for i, expectedElem := range expected {
			err := testCompileIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
			return
		}

		if len(hash.Pairs) != len(expected) {
			t.Errorf("hash has wrong number of pairs. want=%d, got=%d", len(expected), len(hash.Pairs))
		}

		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in pairs")
			}

			err := testCompileIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testingIntegerObject failed: %s", err)
			}
		}
	case *object.Null:
		if actual != vm.Null {
			t.Errorf("object is not null: %T (%+v)", actual, actual)
		}
	}
}

func testVMBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}

func testVMStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not string. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}

	return nil
}

func TestVMIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * (2 + 10)", 60},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	runVmTests(t, tests)
}

func TestVMBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
	}

	runVmTests(t, tests)
}

func TestVMConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if(true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", vm.Null},
		{"if (false) { 10 }", vm.Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}

	runVmTests(t, tests)
}

func TestVMGlobalAssignStatements(t *testing.T) {
	tests := []vmTestCase{
		{"int one = 1; one;", 1},
		{"int one = 1; int two = 2; one + two;", 3},
		{"int one = 1; int two = one + one; one + two;", 3},
	}
	runVmTests(t, tests)
}

func TestVMStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"azula"`, "azula"},
		{`"az" + "ula"`, "azula"},
	}

	runVmTests(t, tests)
}

func TestVMArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, tests)
}

func TestVMHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[object.HashKey]int64{},
		},
		{
			"{1=>2, 2=>3}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			"{1 + 1=>2 * 2, 2 / 2=>3+2}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 1}).HashKey(): 5,
			},
		},
	}
	runVmTests(t, tests)
}

func TestVMIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"{1=>1, 3=>4}[3]", 4},
		{"[1, 2, 3][-1]", 3},
	}

	runVmTests(t, tests)
}

func TestVMFunctionCallNoArgs(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			func x(): int { return 5 + 10; };
			x();
			`,
			expected: 15,
		},
		{
			input: `
			func one(): int { return 1; };
			func two(): int { return 2; };
			one() + two();
			`,
			expected: 3,
		},
		{
			input: `
			func one(): int { return 1; return 2; };
			one();
			`,
			expected: 1,
		},
		{
			input: `
			func one(): int { return 1; };
			func two(): int { return one() + 1; };
			two();
			`,
			expected: 2,
		},
		{
			input: `
			func one(): int { };
			one();
			`,
			expected: vm.Null,
		},
	}
	runVmTests(t, tests)
}

func TestVMFunctionCallLocalBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			func x(): int { int y = 5 + 10; return y; };
			x();
			`,
			expected: 15,
		},
		{
			input: `
			func one(): int { int y = 1; return y; };
			func two(): int { int y = 2; return y; };
			one() + two();
			`,
			expected: 3,
		},
		{
			input: `
			int zero = 0;
			func one(): int { return zero + 1; };
			func two(): int { return zero + one() + 1; };
			two();
			`,
			expected: 2,
		},
	}
	runVmTests(t, tests)
}

func TestVMFunctionCallWithArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			func x(int y): int { return y; };
			x(15);
			`,
			expected: 15,
		},
		{
			input: `
			func sum(int y, int z): int { return y + z; };
			sum(10, 20);
			`,
			expected: 30,
		},
	}
	runVmTests(t, tests)
}

func TestVMCallingFunctionWithWrongArgs(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    "func x(): int { 5; }; x(5);",
			expected: "wrong number of arguments. want=0, got=1",
		},
	}
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := vm.New(comp.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error but resulted in none")
		}

		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error: want=%q, got=%q", tt.expected, err)
		}
	}
}

func TestVMRecursiveFibonacci(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			func fib(int x): int {
				if(x == 0) {
					return 0;
				}
				if(x == 1) {
					return 1;
				}
				return fib(x - 1) + fib(x - 2);
			};
			fib(15);
			`,
			expected: 610,
		},
	}

	runVmTests(t, tests)
}
