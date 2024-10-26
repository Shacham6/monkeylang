package unsafestack

import (
	"fmt"
	"testing"
)

func assertStackSize[T any](t *testing.T, stack UnsafeSizedStack[T], expectedSize int) error {
	if stack.Size() != expectedSize {
		return fmt.Errorf("Size() wrong. got = %d, want = %d", stack.Size(), expectedSize)
	}
	return nil
}

func TestSimple(t *testing.T) {
	stack := Make[int](10)

	if err := assertStackSize(t, stack, 0); err != nil {
		t.Errorf("assertStackSize: %s", err)
	}

	stack.Push(1)
	if err := assertStackSize(t, stack, 1); err != nil {
		t.Errorf("assertStackSize: %s", err)
	}

	val := stack.Pop()
	if val != 1 {
		t.Errorf("val is wrong, got = %d, want = %d", val, 1)
	}

	if err := assertStackSize(t, stack, 0); err != nil {
		t.Errorf("assertStackSize: %s", err)
	}
}
