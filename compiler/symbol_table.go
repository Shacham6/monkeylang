package compiler

type SymbolScope string

const (
	LocalScope   SymbolScope = "LOCAL"
	GlobalScope  SymbolScope = "GLOBAL"
	BuiltinScope SymbolScope = "BUILTIN"
	FreeScope    SymbolScope = "FREE"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

func NewSymbol(name string, scope SymbolScope, index int) Symbol {
	return Symbol{name, scope, index}
}

type SymbolTable struct {
	store       map[string]Symbol
	FreeSymbols []Symbol

	numDefintions int
	parent_       *SymbolTable
	isEnclosed    bool
}

func NewSymbolTable() *SymbolTable {
	return newEnclosedSymbolTable(nil)
}

func newEnclosedSymbolTable(parent *SymbolTable) *SymbolTable {
	store := map[string]Symbol{}
	freeSymbols := []Symbol{}
	numDefinitions := 0
	return &SymbolTable{store, freeSymbols, numDefinitions, parent, parent != nil}
}

func (s *SymbolTable) Define(name string) Symbol {
	var scope SymbolScope
	_, hasParent := s.parent()
	if hasParent {
		scope = LocalScope
	} else {
		scope = GlobalScope
	}

	symbol := NewSymbol(name, scope, s.numDefintions)
	s.store[name] = symbol
	s.numDefintions++
	return symbol
}

func (s *SymbolTable) DefineBuiltin(idx int, name string) Symbol {
	symbol := NewSymbol(name, BuiltinScope, idx)
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)
	symbol := Symbol{
		Name:  original.Name,
		Index: len(s.FreeSymbols) - 1,
		Scope: FreeScope,
	}
	s.store[original.Name] = symbol
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
	if ok {
		return sym, ok
	}

	parent, ok := s.parent()
	if !ok {
		var sy Symbol
		return sy, false
	}

	sym, ok = parent.Resolve(name)
	if !ok {
		return sym, ok
	}

	if sym.Scope == GlobalScope || sym.Scope == BuiltinScope {
		return sym, ok
	}

	// If were here must mean the scope is either LOCAL or FREE
	free := s.DefineFree(sym)
	return free, true
}

func (s *SymbolTable) parent() (*SymbolTable, bool) {
	return s.parent_, s.parent_ != nil
}

func (s *SymbolTable) SpawnScoped() *SymbolTable {
	return newEnclosedSymbolTable(s)
}
