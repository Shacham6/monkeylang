// This file contains within it few panics. Some of them are due to pending the error handling - and a
// a real legitimate few are legitimate (like the one in `Eval`). That's the explanation.

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

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node) object.Object {
	switch v := node.(type) {
	case *ast.Program:
		return evalProgram(v)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(v)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(v)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(v.Value())
	case *ast.Identifier:
		return evalIdentifier(v)
	case *ast.PrefixExpression:
		right := Eval(v.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(v.Operator, right)
	case *ast.InfixExpression:
		left := Eval(v.Left)
		if isError(left) {
			return left
		}

		right := Eval(v.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpression(v.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(v)
	case *ast.IfExpression:
		return evalIfExpression(v)
	case *ast.ReturnStatement:
		val := Eval(v.ReturnValue)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	}
	panic(fmt.Sprintf("Cannot handle node of type %T", node))
}

func evalProgram(p *ast.Program) object.Object {
	var result object.Object

	for _, statement := range p.Statements {
		result = Eval(statement)

		switch v := result.(type) {
		case *object.ReturnValue:
			return v.Value
		case *object.Error:
			return v
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements() {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}

	}

	return result
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition())
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence())
	} else if ie.Alternative().Ok() {
		return Eval(ie.Alternative().Content())
	} else {
		return &NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case &NULL:
		return false
	case &TRUE:
		return true
	case &FALSE:
		return false
	default:
		return true
	}
}

func evalExpressionStatement(es *ast.ExpressionStatement) object.Object {
	return Eval(es.Expression)
}

func evalIntegerLiteral(il *ast.IntegerLiteral) object.Object {
	return &object.Integer{Value: il.Value}
}

func nativeBoolToBooleanObject(b bool) object.Object {
	if b {
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

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", op, right.Type())
	}

	// panic(fmt.Sprintf("Operator %s not supported yet", op))
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	case op == "==":
		return nativeBoolToBooleanObject(left == right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
	// panic(fmt.Sprintf("Operator '%s' not supported between %s and %s", op, left.Type(), right.Type()))
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	// panic("Operator '%s' not supported between integers")
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
		return newError("unknown operator: -%s", right.Type())
		// panic("We don't support minus prefix operators on non numbers currently")
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
