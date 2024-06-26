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

func Eval(node ast.Node, env *object.Environment) object.Object {
	// make sure all branches return a value
	switch v := node.(type) {
	case *ast.Program:
		return evalProgram(v, env)

	case *ast.ExpressionStatement:
		return evalExpressionStatement(v, env)

	case *ast.IntegerLiteral:
		return evalIntegerLiteral(v)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(v.Value())

	case *ast.Identifier:
		return evalIdentifier(v, env)

	case *ast.LetStatement:
		val := Eval(v.Value, env)
		if isError(val) {
			return val
		}
		env.Set(v.Name.Value, val)
		return &NULL

	case *ast.FunctionLiteral:
		params := v.Parameters()
		body := v.Body()
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(v.Function(), env)
		if isError(function) {
			return function
		}

		args := []object.Object{}
		for _, a := range v.Arguments() {
			res := Eval(a, env)
			if isError(res) {
				return res
			}
			args = append(args, res)
		}
		return applyFunction(function, args)

	case *ast.PrefixExpression:
		right := Eval(v.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(v.Operator, right)

	case *ast.InfixExpression:
		left := Eval(v.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(v.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(v.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(v, env)

	case *ast.IfExpression:
		return evalIfExpression(v, env)

	case *ast.ReturnStatement:
		val := Eval(v.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}
	panic(fmt.Sprintf("Cannot handle node of type %T", node))
}

func evalProgram(p *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range p.Statements {
		result = Eval(statement, env)

		switch v := result.(type) {
		case *object.ReturnValue:
			return v.Value
		case *object.Error:
			return v
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements() {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}

	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition(), env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence(), env)
	} else if ie.Alternative().Ok() {
		return Eval(ie.Alternative().Content(), env)
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

func evalExpressionStatement(es *ast.ExpressionStatement, env *object.Environment) object.Object {
	return Eval(es.Expression, env)
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

func evalIdentifier(i *ast.Identifier, env *object.Environment) object.Object {
	if i.Value == "null" {
		return &NULL
	}
	val, ok := env.Get(i.Value)
	if !ok {
		return newError("identifier not found: %s", i.Value)
	}
	return val
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

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("trying to call what is not a function: %s", fn.Type())
	}
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := fn.Env.NewScoped()

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	returnValue, ok := obj.(*object.ReturnValue)
	if ok {
		return returnValue.Value
	}
	return obj
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
