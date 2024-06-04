package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch v := node.(type) {
	case *ast.Program:
		return evalProgram(v)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(v)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: v.Value}
	}
	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)
	}

	return result
}

func evalExpressionStatement(es *ast.ExpressionStatement) object.Object {
	return Eval(es.Expression)
}
