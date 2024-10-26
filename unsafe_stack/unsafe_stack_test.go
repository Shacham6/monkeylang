package unsafestack

import (
	"fmt"
	"testing"
)

func assertStackSize[T any](stack UnsafeSizedStack[T], expectedSize int) error {
	if stack.Size() != expectedSize {
		return fmt.Errorf("Size() wrong. got = %d, want = %d", stack.Size(), expectedSize)
	}
	return nil
}

func TestSimple(t *testing.T) {
	stack := Make[int](10)

	if err := assertStackSize(stack, 0); err != nil {
		t.Errorf("assertStackSize: %s", err)
	}

	stack.Push(1)
	if err := assertStackSize(stack, 1); err != nil {
		t.Errorf("assertStackSize: %s", err)
	}

	if actual := stack.Current(); actual != 1 {
		t.Errorf("actual is wrong, got = %d, want = %d", actual, 1)
	}

	if actual := stack.Pop(); actual != 1 {
		t.Errorf("actual is wrong, got = %d, want = %d", actual, 1)
	}

	if err := assertStackSize(stack, 0); err != nil {
		t.Errorf("assertStackSize: %s", err)
	}
}
