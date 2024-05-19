package optional_test

import (
	"monkey/optional"
	"testing"
)

func TestOkWhenNotNil(t *testing.T) {
	num := 10
	op := optional.Make(&num)
	if !op.Ok() {
		t.Fatalf("op.Ok() got = %v, expect = %v", op.Ok(), true)
	}
}
