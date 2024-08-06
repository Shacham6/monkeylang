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

func ExpandMacros(program *ast.Program, env *object.Environment) ast.Node {
	node, err := ast.Modify(program, func(node ast.Node) (ast.Node, error) {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node, nil
		}

		macro, ok := isMacroCall(callExpression, env)
		if !ok {
			return node, nil
		}

		args := quoteArgs(callExpression)
		evalEnv := extendMacroEnv(macro, args)

		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros") // TODO: Change into proper error
		}

		return quote.Node, nil
	})
	if err != nil {
		panic(err)
	}

	return node
}

func isMacroCall(callExpression *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	ident, ok := callExpression.Function().(*ast.Identifier)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(ident.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

func quoteArgs(callExpression *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}
	for _, a := range callExpression.Arguments() {
		args = append(args, &object.Quote{Node: a})
	}
	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := macro.Env.NewScoped()

	for paramIdx, param := range macro.Parameters {
		extended.Set(param.Value, args[paramIdx])
	}

	return extended
}
