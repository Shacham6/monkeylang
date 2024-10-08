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
		// vmtest.New("if (false) {false;}", nil),
	})
}
