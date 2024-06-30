package object

import (
	"fmt"
	"hash/fnv"
)

type HashKey struct {
	objectType ObjectType
	value      uint64
}

func NewHashKey(objectType ObjectType, value uint64) HashKey {
	return HashKey{objectType, value}
}

func ZeroHashKey() HashKey {
	return HashKey{}
}

func (h HashKey) Value() uint64 {
	return h.value
}

func (h HashKey) ObjectType() ObjectType {
	return h.objectType
}

func (h HashKey) Eq(o HashKey) bool {
	return h.value == o.value && h.objectType == o.objectType
}

//
// Implementations for the various objects starting now
//

func (b *Boolean) HashKey() (HashKey, *Error) {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return NewHashKey(b.Type(), value), nil
}

func (i *Integer) HashKey() (HashKey, *Error) {
	return NewHashKey(i.Type(), uint64(i.Value)), nil
}

func (s *String) HashKey() (HashKey, *Error) {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return NewHashKey(s.Type(), h.Sum64()), nil
}

func newTypeNotHashableError(obj Object) *Error {
	return &Error{
		Message: fmt.Sprintf("object of type '%T' is not hashable", obj.Type()),
	}
}

func (e *Error) HashKey() (HashKey, *Error) {
	return ZeroHashKey(), newTypeNotHashableError(e)
}

func (a *Array) HashKey() (HashKey, *Error) {
	return ZeroHashKey(), newTypeNotHashableError(a)
}

func (n *Null) HashKey() (HashKey, *Error) {
	return ZeroHashKey(), newTypeNotHashableError(n)
}

func (f *Function) HashKey() (HashKey, *Error) {
	return ZeroHashKey(), newTypeNotHashableError(f)
}

func (r *ReturnValue) HashKey() (HashKey, *Error) {
	return ZeroHashKey(), newTypeNotHashableError(r)
}

func (b *Builtin) HashKey() (HashKey, *Error) {
	return ZeroHashKey(), newTypeNotHashableError(b)
}

func (h *Hash) HashKey() (HashKey, *Error) {
	return ZeroHashKey(), newTypeNotHashableError(h)
}
