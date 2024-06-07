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
	case *ast.PrefixExpression:
		return evalPrefixExpression(v)
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
	panic("We don't support identifiers yet")
}

func evalPrefixExpression(p *ast.PrefixExpression) object.Object {
	right := Eval(p.Right)
	return resolvePrefixResult(p.Operator, right)
}

func resolvePrefixResult(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	}

	panic(fmt.Sprintf("Operator %s not supported yet", op))
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case &TRUE:
		return &FALSE
	case &FALSE:
		return &TRUE
	case &NULL:
		return &TRUE
	default:
		return &FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		panic("We don't support minus prefix operators on non numbers currently")
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
