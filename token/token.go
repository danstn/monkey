package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
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
	MINUS  = "-"
	BANG   = "!"
	STAR   = "*"
	SLASH  = "/"
	GT     = ">"
	LT     = "<"
	EQ     = "=="
	NEQ    = "!="

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
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

func New(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
