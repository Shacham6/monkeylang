package evaluator

import (
	"monkey/object"
)

var builtins = func() map[string]*object.Builtin {
	var bm = map[string]*object.Builtin{}
	for _, bb := range object.Builtins {
		bm[bb.Name] = bb.Builtin
	}
	return bm
}()
