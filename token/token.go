package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = TokenType("ILLEGAL")
	EOF     = TokenType("EOF")

	// Identifiers + literals
	IDENT = TokenType("IDENT")
	INT   = TokenType("INT")

	// Operators
	ASSIGN  = TokenType("ASSIGN")
	PLUS    = TokenType("PLUS")
	MINUS   = TokenType("MINUS")
	BANG    = TokenType("BANG")
	ASTERIX = TokenType("ASTERIX")
	SLASH   = TokenType("SLASH")
	LT      = TokenType("<")
	GT      = TokenType(">")

	// Delimiters
	COMMA     = TokenType("COMMA")
	SEMICOLON = TokenType("SEMICOLON")

	LPAREN = TokenType("LPAREN")
	RPAREN = TokenType("RPAREN")
	LBRACE = TokenType("LBRACE")
	RBRACE = TokenType("RBRACE")

	// Keywords
	FUNCTION = TokenType("FUNCTION")
	LET      = TokenType("LET")
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(rawString string) TokenType {
	if keyword, ok := keywords[rawString]; ok {
		return keyword
	}

	return IDENT
}
