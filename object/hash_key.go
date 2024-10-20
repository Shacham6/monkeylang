package object

import (
	"fmt"
	"hash/fnv"
)

type hashKey interface {
	// Does not explicitly return the Hash object because some objects are not hashable,
	// and so we want to maintain the ability to return an error instead.
	HashKey() (HashKey, error)
}

type HashKey struct {
	objectType ObjectType
	value      uint64
}

func NewHashKey(objectType ObjectType, value uint64) HashKey {
	return HashKey{objectType, value}
}

func ZeroHashKey() HashKey {
	return HashKey{} //nolint:exhaustruct
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

func (b *Boolean) HashKey() (HashKey, error) {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return NewHashKey(b.Type(), value), nil
}

func (i *Integer) HashKey() (HashKey, error) {
	return NewHashKey(i.Type(), uint64(i.Value)), nil
}

func (s *String) HashKey() (HashKey, error) {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return NewHashKey(s.Type(), h.Sum64()), nil
}

func newTypeNotHashableError(obj Object) error {
	return fmt.Errorf("object of type '%T' is not hashable", obj.Type())
}

func (e *Error) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(e)
}

func (a *Array) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(a)
}

func (n *Null) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(n)
}

func (f *Function) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(f)
}

func (f *CompiledFunction) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(f)
}

func (r *ReturnValue) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(r)
}

func (b *Builtin) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(b)
}

func (h *Hash) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(h)
}

func (q *Quote) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(q)
}

func (m *Macro) HashKey() (HashKey, error) {
	return ZeroHashKey(), newTypeNotHashableError(m)
}
