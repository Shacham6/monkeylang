package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t *Token) Eq(o *Token) bool {
	return t.Type == o.Type && t.Literal == o.Literal
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN  = "ASSIGN"
	PLUS    = "PLUS"
	MINUS   = "MINUS"
	BANG    = "BANG"
	ASTERIX = "ASTERIX"
	SLASH   = "SLASH"
	LT      = "LT"
	GT      = "GT"
	EQ      = "EQ"
	NOT_EQ  = "NOT_EQ"

	// Delimiters
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN = "LPAREN"
	RPAREN = "RPAREN"
	LBRACE = "LBRACE"
	RBRACE = "RBRACE"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(rawString string) TokenType {
	if keyword, ok := keywords[rawString]; ok {
		return keyword
	}

	return IDENT
}
