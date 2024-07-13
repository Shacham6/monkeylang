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
