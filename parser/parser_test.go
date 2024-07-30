package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/test"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foo = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		assertLetStatement(t, statement, tt.expectedIdentifier)
	}
}

func assertLetStatement(t *testing.T, s ast.Statement, name string) {
	t.Helper()

	test.AssertEqual(t, s.TokenLiteral(), "let")

	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement, got=%T", s)
	}

	test.AssertEqual(t, letStatement.Name.Value, name)
	test.AssertEqual(t, letStatement.Name.TokenLiteral(), name)
}
