package ast

import (
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.NAME, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.NAME, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if got := program.String(); got != "let myVar = anotherVar;" {
		t.Errorf("program.String() is wrong, got=%q", got)
	}
}
