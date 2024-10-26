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
	"strings"
	"testing"
)

func parse(input string) (*ast.Program, error) {
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("got parsing errors: \n\t%s", strings.Join(p.Errors(), "\n\t"))
	}
	return prog, nil
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

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not *object.String. got = %T (%+V)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("result has wrong value. got = %q, want = %q", actual, expected)
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
	case int:
		if err := testIntegerObject(int64(expected), actual); err != nil {
			t.Fatalf("testIntegerObject failed: %s", err)
		}

	case int64:
		if err := testIntegerObject(expected, actual); err != nil {
			t.Fatalf("testIntegerObject failed: %s", err)
		}

	case bool:
		if err := testBooleanObject(expected, actual); err != nil {
			t.Fatalf("testBooleanObject failed: %s", err)
		}

	case string:
		if err := testStringObject(expected, actual); err != nil {
			t.Fatalf("testStringObject failed: %s", err)
		}

	case nil:
		if err := testNilObject(actual); err != nil {
			t.Fatalf("testNilObject failed: %s", err)
		}

	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Fatalf("actual is not *object.Array, got = %T", actual)
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong num of elements in array, want = %d, got = %d",
				len(expected),
				len(array.Elements))
		}

		for i := 0; i < max(len(array.Elements), len(expected)); i++ {
			actualVal, gotActual := getSafe(array.Elements, i)
			expectedVal, gotExpected := getSafe(expected, i)

			if gotActual != gotExpected {
				t.Errorf(
					"mismatched num of elements in expected and actual array,\n expected[%d] = %t, got[%d] = %t",
					i, gotExpected, i, gotActual)
				continue
			}

			testExpectedObject(t, expectedVal, actualVal)
		}

	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Fatalf("actual is not *object.Hash, got = %T", actual)
		}

		if len(hash.Pairs) != len(expected) {
			t.Errorf("lengths are not equal, got = %d, expected = %d", len(hash.Pairs), len(expected))
		}

		// merge the keys into a single token "map", it will act as an intersection set
		// of all the keys.
		allKeys := map[object.HashKey]struct{}{}
		for key := range expected {
			allKeys[key] = struct{}{}
		}
		for key := range hash.Pairs {
			allKeys[key] = struct{}{}
		}

		// perform the checks using ALL the found keys
		for key := range allKeys {
			expectedValue, gotExpected := expected[key]
			actualValue, gotActual := hash.Pairs[key]

			if gotExpected != gotActual {
				t.Errorf("got key %v in one set but not in the other, exists in expected = %t, exists in actual = %t",
					key, gotExpected, gotActual)
				continue
			}

			testExpectedObject(t, expectedValue, actualValue.Value)
		}

	default:
		t.Fatalf("expectation of type %T is not supported yet", expected)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func getSafe[T any](slice []T, n int) (result T, ok bool) {
	if n >= len(slice) {
		ok = false
		return
	}
	return slice[n], true
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
			program, err := parse(tt.input)
			if err != nil {
				t.Fatalf("failed parsing: %s", err)
			}

			comp := compiler.New()
			err = comp.Compile(program)
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
