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
	})
}
