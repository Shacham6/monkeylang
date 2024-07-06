package ast

type IfExpressionAlternative struct {
	ok      bool
	content *BlockStatement
}

func NewIfExpressionAlternative(content *BlockStatement) *IfExpressionAlternative {
	return &IfExpressionAlternative{content != nil, content}
}

func (i *IfExpressionAlternative) Ok() bool {
	return i.ok
}

func (i *IfExpressionAlternative) Content() *BlockStatement {
	return i.content
}

func (i *IfExpressionAlternative) TokenLiteral() string {
	return i.content.TokenLiteral()
}

func (i *IfExpressionAlternative) String() string {
	return i.content.String()
}
