package object

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
)

type deval interface {
	Deval() (ast.Node, error)
}

func newDevalForTypeNotSupportedError(t interface{}) error {
	return fmt.Errorf("object of type %T cannot be restored into a de-evaluated state", t)
}

func (b *Boolean) Deval() (ast.Node, error) {
	var tokenType token.TokenType
	if b.Value {
		tokenType = token.TRUE
	} else {
		tokenType = token.FALSE
	}

	return ast.NewBoolean(
		token.Token{
			Type:    tokenType,
			Literal: fmt.Sprintf("%v", b.Value),
		},
		b.Value,
	), nil
}

func (s *String) Deval() (ast.Node, error) {
	return ast.NewStringLiteral(
		token.Token{
			Type:    token.STRING,
			Literal: s.Value,
		},
		s.Value,
	), nil
}

func (a *Array) Deval() (ast.Node, error) {
	elements := []ast.Expression{}
	for _, i := range a.Elements {
		res, err := i.Deval()
		if err != nil {
			return nil, err
		}

		resExpression := res.(ast.Expression)

		elements = append(elements, resExpression)
	}

	return ast.NewArrayLiteral(token.New(token.LBRACKET, "["), elements), nil
}

func (n *Null) Deval() (ast.Node, error) {
	return ast.NewIdentifier(token.New(token.IDENT, "null"), "null"), nil
}

func (e *Error) Deval() (ast.Node, error) {
	return nil, newDevalForTypeNotSupportedError(e)
}

func (f *Function) Deval() (ast.Node, error) {
	return nil, newDevalForTypeNotSupportedError(f)
}

func (cf *CompiledFunction) Deval() (ast.Node, error) {
	return nil, newDevalForTypeNotSupportedError(cf)
}

func (r *ReturnValue) Deval() (ast.Node, error) {
	return nil, newDevalForTypeNotSupportedError(r)
}

func (b *Builtin) Deval() (ast.Node, error) {
	return nil, newDevalForTypeNotSupportedError(b)
}

func (h *Hash) Deval() (ast.Node, error) {
	return nil, newDevalForTypeNotSupportedError(h)
}

func (q *Quote) Deval() (ast.Node, error) {
	return q.Node, nil
}

func (i *Integer) Deval() (ast.Node, error) {
	t := token.Token{
		Type:    token.INT,
		Literal: fmt.Sprintf("%d", i.Value),
	}
	return ast.NewIntegerLiteral(t, i.Value), nil
}

func (m *Macro) Deval() (ast.Node, error) {
	return nil, newDevalForTypeNotSupportedError(m)
}
