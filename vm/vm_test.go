package vm_test

import (
	"monkey/vm/internal/vmtest"
	"testing"
)

func TestIntegerArithmetic(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("1", 1),
		vmtest.New("2", 2),
		vmtest.New("1 + 2", 3),
		vmtest.New("1 - 2", -1),
		vmtest.New("1 * 2", 2),
		vmtest.New("4 / 2", 2),
		vmtest.New("((1 + 2 - 1) * 2) / 2", 2),
		vmtest.New("4 / 2 * 2 + 2 - 2", 4),
		vmtest.New("-5", -5),
		vmtest.New("-10", -10),
		vmtest.New("-50 + 100 + -50", 0),
		vmtest.New("(5 + 10 * 2 + 15 / 3) * 2 + -10", 50),
	})
}

func TestBooleanExpressions(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("true", true),
		vmtest.New("false", false),
		vmtest.New("1 < 2", true),
		vmtest.New("1 > 2", false),
		vmtest.New("1 < 1", false),
		vmtest.New("1 > 1", false),
		vmtest.New("1 == 1", true),
		vmtest.New("1 != 1", false),
		vmtest.New("1 == 2", false),
		vmtest.New("1 != 2", true),
		vmtest.New("true == true", true),
		vmtest.New("false == true", false),
		vmtest.New("false == false", true),
		vmtest.New("true == false", false),
		vmtest.New("true != false", true),
		vmtest.New("false != true", true),
		vmtest.New("(1 > 2) == false", true),
		vmtest.New("(1 < 2) == false", false),
		vmtest.New("(1 > 2) == true", false),
		vmtest.New("(1 < 2) == true", true),
		vmtest.New("!true", false),
		vmtest.New("!false", true),
		vmtest.New("!5", false),
		vmtest.New("!!true", true),
		vmtest.New("!!false", false),
		vmtest.New("!!5", true),
		vmtest.New("!null", true),
		vmtest.New("!!null", false),
	})
}

func TestNil(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("null", nil),
	})
}

func TestConditionals(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("if (true) { 10 }", 10),
		vmtest.New("if (true) { 10 } else { 20 }", 10),
		vmtest.New("if (false) { 10 } else { 20 }", 20),
		vmtest.New("if (true) { 5 + 5 } else { 20 }", 10),
		vmtest.New("if (false) { 10 } else { 10 + 10 }", 20),
		vmtest.New("if (1) {10}", 10),
		vmtest.New("if (1 - 1) {10} else {20}", 20),
		vmtest.New("if (1 < 2) {10} else {20}", 10),
		vmtest.New("if (1 > 2) {10} else {20}", 20),
		vmtest.New("if (true) {10}; 20", 20),
		vmtest.New("if (false) {false;}", nil),
		vmtest.New("if (null) {10} else {20}", 20),
		vmtest.New("if (null == null) {10}", 10),
	})
}

func TestGlobalLetStatements(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("let a = 1; a", 1),
		vmtest.New("let a = 1; let b = 2; a + b", 3),
		vmtest.New("let a = 1; let b = a + a; a + b", 3),
	})
}

func TestStringExpressions(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New(`"lol"`, "lol"),
		vmtest.New(`"mon" + "key"`, "monkey"),
		vmtest.New(`"mon" + "key" + "banana"`, "monkeybanana"),
	})
}

func TestArrayExpressions(t *testing.T) {
	vmtest.RunVmTests(t, []vmtest.VmTestCase{
		vmtest.New("[]", []int{}),
		vmtest.New("[1, 2, 3]", []int{1, 2, 3}),
		vmtest.New("[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}),
	})
}
