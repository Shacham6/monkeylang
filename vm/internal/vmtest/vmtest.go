// vmtest is an internal package designated for testing the vm.
//
// it will contain various utilities and common operations for that express purpose,
// and as such it is in the 'internal' package.
package vmtest

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
	"testing"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not *object.Integer. got = %T (%+V)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("result has wrong value. got = %d, want = %d", result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not *object.Boolean. got = %T (%+V)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("result has wrong value. got = %t, want = %t", actual, expected)
	}

	return nil
}

func testNilObject(actual object.Object) error {
	_, ok := actual.(*object.Null)
	if ok {
		return nil
	}

	return fmt.Errorf("object is not null. got = %T", actual)
}

func testExpectedObject(t *testing.T, expected any, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int64:
		if err := testIntegerObject(expected, actual); err != nil {
			t.Fatalf("testIntegerObject failed: %s", err)
		}

	case bool:
		if err := testBooleanObject(expected, actual); err != nil {
			t.Fatalf("testBooleanObject failed: %s", err)
		}

	case nil:
		if err := testNilObject(actual); err != nil {
			t.Fatalf("testNilObject failed: %s", err)
		}
	}
}

type VmTestCase struct {
	input    string
	expected any
}

func New(input string, expected any) VmTestCase {
	return VmTestCase{input, expected}
}

func RunVmTests(t *testing.T, tests []VmTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
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
		})
	}
}
