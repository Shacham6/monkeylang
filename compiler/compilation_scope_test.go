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

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scope index wrong, got = %d, want = %d", compiler.scopeIndex, 0)
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
