package evaluator_test

import (
	"monkey/ast"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func testParseProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Logf("got unexpected parsing %d errors", len(p.Errors()))
		for _, parsingErr := range p.Errors() {
			t.Errorf("PARSING ERROR | %s", parsingErr)
		}

		t.FailNow()
	}

	return program
}

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1
	let function = fn(x, y) {x + y};
	let mymacro = macro(x, y) {x + y; };
	`

	env := object.NewEnvironment()
	program := testParseProgram(t, input)

	evaluator.DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("wrong number of statements, got = %d, want = %d", len(program.Statements), 2)
	}

	_, ok := env.Get("number")
	if ok {
		t.Fatal("number should not be defined")
	}

	_, ok = env.Get("function")
	if ok {
		t.Fatal("function should not be defined")
	}

	obj, ok := env.Get("mymacro")
	if !ok {
		t.Fatal("mymacro not in environment")
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("object is not macro, got = %T", obj)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("wrong number of macro parameters, got = %d, want = %d", len(macro.Parameters), 2)
	}

	if macro.Parameters[0].String() != "x" {
		t.Errorf("parameter is not 'x', got = %q", macro.Parameters[0])
	}

	if macro.Parameters[1].String() != "y" {
		t.Errorf("parameter is not 'y', got = %q", macro.Parameters[1])
	}

	expectedBody := "(block (expr (infix x + y)))"
	if macro.Body.String() != expectedBody {
		t.Errorf("macro body is not %s, got = %s", expectedBody, macro.Body.String())
	}
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`
			let infixExpression = macro() { quote(1 + 2); };
			infixExpression();
			`,
			`(1 + 2)`,
		},
		{
			`
			let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); }
			reverse(2 + 2, 10 - 5);
			`,
			`(10 - 5) - (2 + 2)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expected := testParseProgram(t, tt.expected)
			program := testParseProgram(t, tt.input)

			env := object.NewEnvironment()
			evaluator.DefineMacros(program, env)
			expanded := evaluator.ExpandMacros(program, env)

			if expanded.String() != expected.String() {
				t.Errorf("expanded not equal to expected, got = %q, want = %q", expanded.String(), expected.String())
			}
		})
	}
}
