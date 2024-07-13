package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	resultNode, err := ast.Modify(quoted, func(node ast.Node) (ast.Node, error) {
		if !isUnquotedCall(node) {
			return node, nil
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node, nil
		}

		if len(call.Arguments()) != 1 {
			return node, nil
		}

		evalRes := Eval(call.Arguments()[0], env)
		astNode, err := evalRes.Deval()
		return astNode, err
		// return convertObjectToASTNode(evalRes), nil
	})
	if err != nil {
		panic(err) // Don't know how to match this currently
	}

	return resultNode
}

func isUnquotedCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpression.Function().TokenLiteral() == "unquote"
}

// convertObjectToASTNode performs, for all intents and purposes, the opposite of `Eval`.
// func convertObjectToASTNode(obj object.Object) ast.Node {
// switch obj := obj.(type) {
// case *object.Integer:
// 	t := token.Token{
// 		Type:    token.INT,
// 		Literal: fmt.Sprintf("%d", obj.Value),
// 	}
// 	return ast.NewIntegerLiteral(t, obj.Value)
//
// case *object.Boolean:
// 	var tokenType token.TokenType
// 	if obj.Value {
// 		tokenType = token.TRUE
// 	} else {
// 		tokenType = token.FALSE
// 	}
//
// 	return ast.NewBoolean(
// 		token.Token{
// 			Type:    tokenType,
// 			Literal: fmt.Sprintf("%v", obj.Value),
// 		},
// 		obj.Value,
// 	)
//
// case *object.String:
// 	return ast.NewStringLiteral(
// 		token.Token{
// 			Type:    token.STRING,
// 			Literal: obj.Value,
// 		},
// 		obj.Value,
// 	)
// }
// panic(fmt.Sprintf("object of type %T not supported yet", obj))
// }
