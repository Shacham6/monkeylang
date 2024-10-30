package compiler

import (
	"monkey/code"
	"testing"
)

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got = %d, want = %d", compiler.scopeIndex, 0)
	}
	globalSymbolTable := compiler.symbolTable

	compiler.emit(code.OpMul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex wrong. got = %d, want = %d", compiler.scopeIndex, 1)
	}

	compiler.emit(code.OpSub)

	if len(compiler.scope().Instructions) != 1 {
		t.Errorf("instructions length wrong, got = %d", len(compiler.scopes[compiler.scopeIndex].Instructions))
	}

	last := compiler.scope().LastInstruction
	if last.Opcode != code.OpSub {
		t.Errorf("LastInstruction op code wrong, got = %d, want = %d", last.Opcode, code.OpSub)
	}

	parTable, hasParTable := compiler.symbolTable.parent()
	if !hasParTable {
		t.Errorf("hasParTable is %v, expected %v", hasParTable, true)
	}

	if parTable != globalSymbolTable {
		t.Errorf("compiler did not enclose symbol table")
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scope index wrong, got = %d, want = %d", compiler.scopeIndex, 0)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("returned symbol table is not the previously enclosed symbol table")
	}

	parTable, hasParTable = compiler.symbolTable.parent()
	if hasParTable {
		t.Errorf("hasParTable should be none")
	}

	if parTable != nil {
		t.Errorf("parTable is not nil when it should be nil")
	}

	compiler.emit(code.OpAdd)

	if len(compiler.scope().Instructions) != 2 {
		t.Errorf("instructions length wrong. got = %d, want = %d", len(compiler.scope().Instructions), 2)
	}

	previous := compiler.scope().PrevInstruction
	if previous.Opcode != code.OpMul {
		t.Errorf("PrevInstruction.Opcode wrong, got = %d, want = %d", previous.Opcode, code.OpMul)
	}
}
