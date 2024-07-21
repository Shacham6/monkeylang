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
		res, err := modifyIntoType[Statement](statement, modify)
		if err != nil {
			return err
		}
		p.Statements[i] = res
	}
	return nil
}

func (e *ExpressionStatement) modify(modify ModifierFunc) error {
	res, err := modifyIntoType[Expression](e.Expression, modify)
	if err != nil {
		return err
	}

	e.Expression = res
	return nil
}

func (e *InfixExpression) modify(modify ModifierFunc) error {
	leftRes, err := modifyIntoType[Expression](e.Left, modify)
	if err != nil {
		return err
	}

	rightRes, err := modifyIntoType[Expression](e.Right, modify)
	if err != nil {
		return err
	}

	e.Left = leftRes
	e.Right = rightRes

	return nil
}

func (p *PrefixExpression) modify(modify ModifierFunc) error {
	rightRes, err := modifyIntoType[Expression](p.Right, modify)
	if err != nil {
		return err
	}

	p.Right = rightRes
	return nil
}

func (a *ArrayLiteral) modify(modify ModifierFunc) error {
	for i, el := range a.Elements {
		elRes, err := modifyIntoType[Expression](el, modify)
		if err != nil {
			return err
		}
		a.Elements[i] = elRes
	}
	return nil
}

func (i *IndexExpression) modify(modify ModifierFunc) error {
	leftRes, err := modifyIntoType[Expression](i.Left(), modify)
	if err != nil {
		return err
	}
	i.left = leftRes

	indexRes, err := modifyIntoType[Expression](i.Index(), modify)
	if err != nil {
		return err
	}
	i.index = indexRes

	return nil
}

func (i *IfExpression) modify(modify ModifierFunc) error {
	var err error
	i.condition, err = modifyIntoType[Expression](i.condition, modify)
	if err != nil {
		return err
	}
	i.consequence, err = modifyIntoType[*BlockStatement](i.consequence, modify)
	if err != nil {
		return err
	}
	i.alternative, err = modifyIntoType[*IfExpressionAlternative](i.alternative, modify)
	if err != nil {
		return err
	}
	return nil
}

func (i *IfExpressionAlternative) modify(modify ModifierFunc) error {
	if !i.ok {
		return nil
	}

	res, err := modifyIntoType[*BlockStatement](i.content, modify)
	if err != nil {
		return err
	}

	i.content = res
	return nil
}

func (b *BlockStatement) modify(modify ModifierFunc) error {
	for i, el := range b.statements {
		statementRes, err := modifyIntoType[Statement](el, modify)
		if err != nil {
			return err
		}
		b.statements[i] = statementRes
	}
	return nil
}

func (r *ReturnStatement) modify(modify ModifierFunc) error {
	returnRes, err := modifyIntoType[Expression](r.ReturnValue, modify)
	if err != nil {
		return err
	}
	r.ReturnValue = returnRes
	return nil
}

func (l *LetStatement) modify(modify ModifierFunc) error {
	letRes, err := modifyIntoType[Expression](l.Value, modify)
	if err != nil {
		return err
	}
	l.Value = letRes
	return nil
}

func (f *FunctionLiteral) modify(modify ModifierFunc) error {
	bodyRes, err := modifyIntoType[*BlockStatement](f.body, modify)
	if err != nil {
		return err
	}
	f.body = bodyRes

	for i, p := range f.parameters {
		paramRes, err := modifyIntoType[*Identifier](p, modify)
		if err != nil {
			return err
		}
		f.parameters[i] = paramRes
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
		modKey, err := modifyIntoType[Expression](key, modify)
		if err != nil {
			return err
		}

		modVal, err := modifyIntoType[Expression](val, modify)
		if err != nil {
			return err
		}
		modifiedPairs[modKey] = modVal
	}

	h.pairs = modifiedPairs

	return nil
}

func (c *CallExpression) modify(modify ModifierFunc) error {
	c.function, _ = modifyIntoType[Expression](c.function, modify)

	for i, arg := range c.arguments {
		argRes, err := modifyIntoType[Expression](arg, modify)
		if err != nil {
			return err
		}
		c.arguments[i] = argRes
	}

	return nil
}

func (m *MacroLiteral) modify(modify ModifierFunc) error {
	// Logically I don't think that we dig into "macro" literals specifically because modify is a macro
	// concept to begin with. But on the other hand no reason these should be mutually exclusive...
	// I don't know.
	return nil
}
