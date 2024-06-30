package object_test

import (
	"monkey/object"
	"testing"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &object.String{Value: "hello world"}
	hello2 := &object.String{Value: "hello world"}

	diff1 := &object.String{Value: "praise the sun"}
	diff2 := &object.String{Value: "praise the sun"}

	if ensureHashSuccess(t, hello1) != ensureHashSuccess(t, hello2) {
		t.Errorf("strings with same content and different hash keys")
	}

	if ensureHashSuccess(t, diff1) != ensureHashSuccess(t, diff2) {
		t.Errorf("strings with same content and different hash keys")
	}

	if ensureHashSuccess(t, hello1) == ensureHashSuccess(t, diff1) {
		t.Errorf("strings with different content have the same hash keys")
	}
}

func ensureHashSuccess(t *testing.T, o object.Object) object.HashKey {
	hashKey, err := o.HashKey()
	if err != nil {
		t.Fatalf("obj (type '%T') failed hashing", o)
	}

	return hashKey
}
