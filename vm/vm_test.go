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
	})
}
