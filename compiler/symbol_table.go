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

type SymbolTable struct {
	store         map[string]Symbol
	numDefintions int
}

func NewSymbolTable() *SymbolTable {
	s := map[string]Symbol{}
	return &SymbolTable{s, 0}
}
