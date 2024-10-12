package compiler

type SymbolScope string

const (
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
}

func NewSymbolTable() *SymbolTable {
	s := map[string]Symbol{}
	return &SymbolTable{s, 0}
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := NewSymbol(name, GlobalScope, s.numDefintions)
	s.store[name] = symbol
	s.numDefintions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	return symbol, ok
}
