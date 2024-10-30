package compiler

type SymbolScope string

const (
	LocalScope  SymbolScope = "LOCAL"
	GlobalScope SymbolScope = "GLOBAL"
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
	store         map[string]Symbol
	numDefintions int
	parent_       *SymbolTable
	isEnclosed    bool
}

func NewSymbolTable() *SymbolTable {
	return newEnclosedSymbolTable(nil)
}

func newEnclosedSymbolTable(parent *SymbolTable) *SymbolTable {
	store := map[string]Symbol{}
	numDefinitions := 0
	return &SymbolTable{store, numDefinitions, parent, parent != nil}
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

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	if ok {
		return symbol, ok
	}

	if parent, hasParent := s.parent(); hasParent {
		return parent.forwardedResolve(name)
	}

	var sy Symbol
	return sy, false
}

func (s *SymbolTable) forwardedResolve(name string) (Symbol, bool) {
	parent, ok := s.parent()
	if !ok {
		s, ok := s.store[name]
		return s, ok
	}
	return parent.forwardedResolve(name)
}

func (s *SymbolTable) parent() (*SymbolTable, bool) {
	return s.parent_, s.parent_ != nil
}

func (s *SymbolTable) SpawnScoped() *SymbolTable {
	return newEnclosedSymbolTable(s)
}
