package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	// This implementation only walks and finds top level macros, it does not find
	// inner macros.
	for i, statement := range program.Statements {
		if !isMacroDefinition(statement) {
			continue
		}
		addMacro(statement, env)
		definitions = append(definitions, i)
	}

	for i := len(definitions) - 1; i >= 0; i-- {
		definitionIndex := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

func isMacroDefinition(node ast.Statement) bool {
	letStatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}

	// This does not seem recursive... Or maybe it is? idk
	_, ok = letStatement.Value.(*ast.MacroLiteral)

	return ok
}

func addMacro(node ast.Statement, env *object.Environment) {
	letStatement, _ := node.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Parameters: macroLiteral.Parameters(),
		Env:        env,
		Body:       macroLiteral.Body(),
	}

	env.Set(letStatement.Name.Value, macro)
}
