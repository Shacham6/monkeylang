package ast

import (
	"fmt"
)

type ModifierFunc func(Node) (Node, error)

type modifiable interface {
	modify(ModifierFunc) error
}

func Modify(node Node, modifier ModifierFunc) (Node, error) {
	node.modify(modifier)
	return modifier(node)
}

func modifyIntoType[T Node](node Node, modifier ModifierFunc) (T, error) {
	modResult, err := Modify(node, modifier)
	if err != nil {
		return modResult.(T), err
	}
	converted, ok := modResult.(T)
	if !ok {
		return converted, fmt.Errorf("result of modifier of type %T not match expected type", modResult)
	}
	return converted, nil
}

func (p *Program) modify(modify ModifierFunc) error {
	for i, statement := range p.Statements {
		p.Statements[i], _ = modifyIntoType[Statement](statement, modify)
	}
	return nil
}

func (e *ExpressionStatement) modify(modify ModifierFunc) error {
	e.Expression, _ = modifyIntoType[Expression](e.Expression, modify)
	return nil
}

func (e *InfixExpression) modify(modify ModifierFunc) error {
	e.Left, _ = modifyIntoType[Expression](e.Left, modify)
	e.Right, _ = modifyIntoType[Expression](e.Right, modify)
	return nil
}

func (p *PrefixExpression) modify(modify ModifierFunc) error {
	p.Right, _ = modifyIntoType[Expression](p.Right, modify)
	return nil
}

func (a *ArrayLiteral) modify(modify ModifierFunc) error {
	for i, el := range a.Elements {
		a.Elements[i], _ = modifyIntoType[Expression](el, modify)
	}
	return nil
}

func (i *IndexExpression) modify(modify ModifierFunc) error {
	i.left, _ = modifyIntoType[Expression](i.Left(), modify)
	i.index, _ = modifyIntoType[Expression](i.Index(), modify)
	return nil
}

func (i *IfExpression) modify(modify ModifierFunc) error {
	i.condition, _ = modifyIntoType[Expression](i.condition, modify)
	i.consequence, _ = modifyIntoType[*BlockStatement](i.consequence, modify)
	i.alternative, _ = modifyIntoType[*IfExpressionAlternative](i.alternative, modify)
	return nil
}

func (i *IfExpressionAlternative) modify(modify ModifierFunc) error {
	if !i.ok {
		return nil
	}
	i.content, _ = modifyIntoType[*BlockStatement](i.content, modify)
	return nil
}

func (b *BlockStatement) modify(modify ModifierFunc) error {
	for i, el := range b.statements {
		b.statements[i], _ = modifyIntoType[Statement](el, modify)
	}
	return nil
}

func (r *ReturnStatement) modify(modify ModifierFunc) error {
	r.ReturnValue, _ = modifyIntoType[Expression](r.ReturnValue, modify)
	return nil
}

func (l *LetStatement) modify(modify ModifierFunc) error {
	l.Value, _ = modifyIntoType[Expression](l.Value, modify)
	return nil
}

func (f *FunctionLiteral) modify(modify ModifierFunc) error {
	f.body, _ = modifyIntoType[*BlockStatement](f.body, modify)
	for i, p := range f.parameters {
		f.parameters[i], _ = modifyIntoType[*Identifier](p, modify)
	}
	return nil
}

func (i *Identifier) modify(modify ModifierFunc) error { return nil }

func (i *IntegerLiteral) modify(modify ModifierFunc) error { return nil }

func (b *Boolean) modify(modify ModifierFunc) error { return nil }

func (s *StringLiteral) modify(modify ModifierFunc) error { return nil }

func (h *HashLiteral) modify(modify ModifierFunc) error {
	modifiedPairs := map[Expression]Expression{}
	for key, val := range h.pairs {
		modKey, _ := modifyIntoType[Expression](key, modify)
		modVal, _ := modifyIntoType[Expression](val, modify)
		modifiedPairs[modKey] = modVal
	}

	h.pairs = modifiedPairs

	return nil
}

func (c *CallExpression) modify(modify ModifierFunc) error {
	c.function, _ = modifyIntoType[Expression](c.function, modify)

	for i, arg := range c.arguments {
		c.arguments[i], _ = modifyIntoType[Expression](arg, modify)
	}

	return nil
}
