package testutils

import "testing"

func CheckIsA[T any](t *testing.T, val any, failMsg string) *T {
	res, ok := val.(*T)
	if !ok {
		t.Fatalf("%s, got = %T", failMsg, val)
	}
	return res
}
