package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// Token types
const (
	ILLEGAL = "ILLEGAL" // token/char we don't know about
	EOF     = "EOF"

	// identifiers + literals
	IDENT = "IDENT" // add, foo, x, y ...
	INT   = "INT"   // 1234567890

	// operators
	ASSIGN = "="
	PLUS   = "+"

	// delimeters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

func New(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
