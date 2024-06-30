package testutils

import "testing"

func CheckIsA[T any](t *testing.T, val any, failMsg string) *T {
	res, ok := val.(*T)
	if !ok {
		t.Fatalf("%s, got = %T (%+v)", failMsg, val, val)
	}
	return res
}
