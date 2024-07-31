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
	assertProgramNoErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	assertProgramStatements(t, program, 3)

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

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	assertProgramNoErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	assertProgramStatements(t, program, 3)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement is not *ast.ReturnStatement, got=%T", stmt)
			continue
		}
		if lit := returnStmt.TokenLiteral(); lit != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", lit)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertProgramNoErrors(t, p)
	assertProgramStatements(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got=%T", stmt.Expression)
	}

	if got := ident.Value; got != "foobar" {
		t.Errorf("ident.Value is not %s, got=%s", "foobar", got)
	}

	if got := ident.TokenLiteral(); got != "foobar" {
		t.Errorf("ident.TokenLiteral is not %s, got=%s", "foobar", got)
	}
}

// Assetion helpers

func assertProgramStatements(t *testing.T, program *ast.Program, want int) {
	t.Helper()
	if got := len(program.Statements); got != want {
		t.Fatalf("program has an unexpected # of statements: got=%d, want=%d", got, want)
	}
}

func assertLetStatement(t *testing.T, s ast.Statement, name string) {
	test.AssertEqual(t, s.TokenLiteral(), "let")

	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement, got=%T", s)
	}

	test.AssertEqual(t, letStatement.Name.Value, name)
	test.AssertEqual(t, letStatement.Name.TokenLiteral(), name)
}

func assertProgramNoErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors:", len(errors))
	var sep string
	for i, msg := range errors {
		if i == len(errors)-1 {
			sep = "└──"
		} else {
			sep = "├──"
		}
		t.Errorf("\t%s %s", sep, msg)
	}
	t.FailNow()
}
