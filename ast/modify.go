package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}

	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)

	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *ArrayLiteral:
		for i, el := range node.Elements {
			node.Elements[i], _ = Modify(el, modifier).(Expression)
		}

	case *IndexExpression:
		node.left, _ = Modify(node.Left(), modifier).(Expression)
		node.index, _ = Modify(node.Index(), modifier).(Expression)

	case *IfExpression:
		node.condition, _ = Modify(node.condition, modifier).(Expression)
		node.consequence, _ = Modify(node.consequence, modifier).(*BlockStatement)
		node.alternative, _ = Modify(node.alternative, modifier).(*IfExpressionAlternative)

	case *IfExpressionAlternative:
		if !node.ok {
			return modifier(node)
		}
		node.content, _ = Modify(node.content, modifier).(*BlockStatement)

	case *BlockStatement:
		for i, el := range node.statements {
			node.statements[i], _ = Modify(el, modifier).(Statement)
		}

	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)

	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)

	case *FunctionLiteral:
		node.body, _ = Modify(node.body, modifier).(*BlockStatement)
		for i, p := range node.parameters {
			node.parameters[i], _ = Modify(p, modifier).(*Identifier)
		}
	}

	return modifier(node)
}
