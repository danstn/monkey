package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		let five = 5;
		let ten = 10;

		let add = fn(x, y) {
			x + y;
		};

		let res = add(five, ten);
		
		!-/*5;
		5 < 10 > 5;

		if (5 < 10) { return true; } else { return false; }

		10 == 10;
		10 != 9;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"}, // let five = 5;
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"}, // let ten = 10;
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"}, // let add = fn(x, y) { x + y; };
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"}, // let res = add(five, ten);
		{token.IDENT, "res"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"}, // !-/*5;
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.STAR, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"}, // 5 < 10 > 5;
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"}, // if (5 < 10) { return true; } else { return false; }
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"}, // 10 == 10;
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"}, // 10 != 9;
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong token type. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
