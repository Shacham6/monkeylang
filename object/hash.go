package object

import (
	"fmt"
	"strings"
)

type Hash struct {
	Pairs map[HashKey]HashPair
}

type HashPair struct {
	Key   Object
	Value Object
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}

func (h *Hash) Inspect() string {
	pairs := []string{}
	for _, value := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", value.Key.Inspect(), value.Value.Inspect()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}
