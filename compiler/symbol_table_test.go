package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": NewSymbol("a", GlobalScope, 0),
		"b": NewSymbol("b", GlobalScope, 1),
		"c": NewSymbol("c", LocalScope, 0),
		"d": NewSymbol("d", LocalScope, 1),
		"e": NewSymbol("e", LocalScope, 0),
		"f": NewSymbol("f", LocalScope, 1),
	}

	global := NewSymbolTable()

	a := global.Define("a")
	expectSymbol(t, "a", a, expected["a"])

	b := global.Define("b")
	expectSymbol(t, "b", b, expected["b"])

	firstLocal := global.SpawnScoped()

	c := firstLocal.Define("c")
	expectSymbol(t, "c", c, expected["c"])

	d := firstLocal.Define("d")
	expectSymbol(t, "d", d, expected["d"])

	secondLocal := firstLocal.SpawnScoped()

	e := secondLocal.Define("e")
	expectSymbol(t, "e", e, expected["e"])

	f := secondLocal.Define("f")
	expectSymbol(t, "f", f, expected["f"])
}

func expectSymbol(t *testing.T, name string, got Symbol, expected Symbol) {
	t.Helper()

	if got == expected {
		return
	}

	t.Errorf("expected %s = %+v, got = %+v", name, expected, got)
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		NewSymbol("a", GlobalScope, 0),
		NewSymbol("b", GlobalScope, 1),
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s is not defined", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got %+v",
				sym.Name, sym, result)
		}
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := global.SpawnScoped()
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		NewSymbol("a", GlobalScope, 0),
		NewSymbol("b", GlobalScope, 1),
		NewSymbol("c", LocalScope, 0),
		NewSymbol("d", LocalScope, 1),
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got = %+v", sym.Name, sym, result)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := global.SpawnScoped()
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := firstLocal.SpawnScoped()
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				NewSymbol("a", GlobalScope, 0),
				NewSymbol("b", GlobalScope, 1),
				NewSymbol("c", LocalScope, 0),
				NewSymbol("d", LocalScope, 1),
			},
		},
		{
			secondLocal,
			[]Symbol{
				NewSymbol("a", GlobalScope, 0),
				NewSymbol("b", GlobalScope, 1),
				NewSymbol("e", LocalScope, 0),
				NewSymbol("f", LocalScope, 1),
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got = %+v", sym.Name, sym, result)
			}
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := NewSymbolTable()
	firstLocal := global.SpawnScoped()
	secondLocal := firstLocal.SpawnScoped()

	expected := []Symbol{
		NewSymbol("a", BuiltinScope, 0),
		NewSymbol("c", BuiltinScope, 1),
		NewSymbol("e", BuiltinScope, 2),
		NewSymbol("f", BuiltinScope, 3),
	}

	for i, v := range expected {
		global.DefineBuiltin(i, v.Name)
	}

	for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			result, ok := table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s is not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
			}
		}
	}
}
