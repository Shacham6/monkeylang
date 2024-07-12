package ast

type ModifierFunc func(Node) Node

type modifiable interface {
	modify(ModifierFunc)
}

func Modify(node Node, modifier ModifierFunc) Node {
	node.modify(modifier)
	return modifier(node)
}

func (p *Program) modify(modify ModifierFunc) {
	for i, statement := range p.Statements {
		p.Statements[i], _ = Modify(statement, modify).(Statement)
	}
}

func (e *ExpressionStatement) modify(modify ModifierFunc) {
	e.Expression, _ = Modify(e.Expression, modify).(Expression)
}

func (e *InfixExpression) modify(modify ModifierFunc) {
	e.Left, _ = Modify(e.Left, modify).(Expression)
	e.Right, _ = Modify(e.Right, modify).(Expression)
}

func (p *PrefixExpression) modify(modify ModifierFunc) {
	p.Right, _ = Modify(p.Right, modify).(Expression)
}

func (a *ArrayLiteral) modify(modify ModifierFunc) {
	for i, el := range a.Elements {
		a.Elements[i], _ = Modify(el, modify).(Expression)
	}
}

func (i *IndexExpression) modify(modify ModifierFunc) {
	i.left, _ = Modify(i.Left(), modify).(Expression)
	i.index, _ = Modify(i.Index(), modify).(Expression)
}

func (i *IfExpression) modify(modify ModifierFunc) {
	i.condition, _ = Modify(i.condition, modify).(Expression)
	i.consequence, _ = Modify(i.consequence, modify).(*BlockStatement)
	i.alternative, _ = Modify(i.alternative, modify).(*IfExpressionAlternative)
}

func (i *IfExpressionAlternative) modify(modify ModifierFunc) {
	if !i.ok {
		return
	}
	i.content, _ = Modify(i.content, modify).(*BlockStatement)
}

func (b *BlockStatement) modify(modify ModifierFunc) {
	for i, el := range b.statements {
		b.statements[i], _ = Modify(el, modify).(Statement)
	}
}

func (r *ReturnStatement) modify(modify ModifierFunc) {
	r.ReturnValue, _ = Modify(r.ReturnValue, modify).(Expression)
}

func (l *LetStatement) modify(modify ModifierFunc) {
	l.Value, _ = Modify(l.Value, modify).(Expression)
}

func (f *FunctionLiteral) modify(modify ModifierFunc) {
	f.body, _ = Modify(f.body, modify).(*BlockStatement)
	for i, p := range f.parameters {
		f.parameters[i], _ = Modify(p, modify).(*Identifier)
	}
}

func (i *Identifier) modify(modify ModifierFunc) {}

func (i *IntegerLiteral) modify(modify ModifierFunc) {}

func (b *Boolean) modify(modify ModifierFunc) {}

func (s *StringLiteral) modify(modify ModifierFunc) {}

func (h *HashLiteral) modify(modify ModifierFunc) {}

func (c *CallExpression) modify(modify ModifierFunc) {
	c.function, _ = Modify(c.function, modify).(Expression)

	for i, arg := range c.arguments {
		c.arguments[i], _ = Modify(arg, modify).(Expression)
	}
}
