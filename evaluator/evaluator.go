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

	case *ast.StringLiteral:
		return &object.String{Value: v.Value}

	case *ast.Identifier:
		return evalIdentifier(v, env)

	case *ast.HashLiteral:
		return evalHashLiteral(v, env)

	case *ast.IndexExpression:
		leftObj := Eval(v.Left(), env)
		if isError(leftObj) {
			return leftObj
		}
		indexObj := Eval(v.Index(), env)
		if isError(indexObj) {
			return indexObj
		}
		return evalIndexExpression(leftObj, indexObj)

	case *ast.ArrayLiteral:
		evaluatedElements := []object.Object{}
		for _, elNode := range v.Elements {
			elObject := Eval(elNode, env)
			if isError(elObject) {
				return elObject
			}
			evaluatedElements = append(evaluatedElements, elObject)
		}
		return &object.Array{Elements: evaluatedElements}

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
		// Handle the "quote" magic case
		if v.Function().TokenLiteral() == "quote" {
			// TODO(Jajo): Pay attention that this section will simply silently ignore
			// all arguments after the 1st whatever they may be, without even an argument.
			return quote(v.Arguments()[0])
		}

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
	return newError(fmt.Sprintf("Cannot handle node of type %T", node))
}

func evalHashLiteral(h *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := map[object.HashKey]object.HashPair{}

	for keyNode, valueNode := range h.Pairs() {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashed, err := key.HashKey()
		if err != nil {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalHashIndexExpression(left object.Object, index object.Object) object.Object {
	hashObject := left.(*object.Hash)

	key, err := index.HashKey()
	if err != nil {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key]
	if !ok {
		return &NULL
	}

	return pair.Value
}

func evalArrayIndexExpression(left object.Object, index object.Object) object.Object {
	arr := left.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arr.Elements) - 1)

	if idx < 0 || idx > max {
		return &NULL
	}
	return arr.Elements[idx]
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
	if val, ok := env.Get(i.Value); ok {
		return val
	}
	if builtin, ok := builtins[i.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", i.Value)
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
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("trying to call what is not a function: %s", fn.Type())
	}
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
	// case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:  // TODO(Jajo): clean
	case left.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(op, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
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

func evalStringInfixExpression(op string, left object.Object, right object.Object) object.Object {
	if op != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
	leftVal := left.(*object.String).Value
	switch right.Type() {
	case object.STRING_OBJ:
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal}
	case object.INTEGER_OBJ:
		rightVal := right.(*object.Integer).Value
		return &object.String{Value: fmt.Sprintf("%s%d", leftVal, rightVal)}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
