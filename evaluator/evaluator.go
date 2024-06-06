package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = object.Null{}
	TRUE  = object.Boolean{Value: true}
	FALSE = object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch v := node.(type) {
	case *ast.Program:
		return evalProgram(v)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(v)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(v)
	case *ast.Boolean:
		return evalBooleanLiteral(v)
	case *ast.Identifier:
		return evalIdentifier(v)
	}
	panic(fmt.Sprintf("Cannot handle node of type %T", node))
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

func evalIntegerLiteral(il *ast.IntegerLiteral) object.Object {
	return &object.Integer{Value: il.Value}
}

func evalBooleanLiteral(b *ast.Boolean) object.Object {
	if b.Value() {
		return &TRUE
	}
	return &FALSE
}

func evalIdentifier(i *ast.Identifier) object.Object {
	if i.Value == "null" {
		return &NULL
	}
	panic("don't support identifiers yet")
}
