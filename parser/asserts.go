package parser

import (
	"fmt"
	"testing"

	"monkey/ast"
	"monkey/test"
)

// Assertion helpers
// -----------------------------------------------------------------------------

func assertExpressionStatement(t *testing.T, statement ast.Statement) *ast.ExpressionStatement {
	t.Helper()
	stmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not ast.ExpressionStatement, got %T", statement)
	}
	return stmt
}

func assertIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	t.Helper()
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral, got=%T", il)
		return
	}

	test.AssertEqual(t, integ.Value, value)
	test.AssertEqual(t, integ.TokenLiteral(), fmt.Sprintf("%d", value))
}

func assertProgramNotNil(t *testing.T, program *ast.Program) {
	t.Helper()
	if program == nil {
		t.Fatalf("program is nil")
	}
}

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

func assertParserNoErrors(t *testing.T, p *Parser) {
	t.Helper()
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

func assertIdentifier(t *testing.T, exp ast.Expression, want string) {
	t.Helper()
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not ast.Identifier, got=%T", exp)
	}

	if ident.Value != want {
		t.Fatalf("ident.Value is not %s, got=%s", want, ident.Value)
	}

	if ident.TokenLiteral() != want {
		t.Fatalf("ident.TokenLiteral is not %s, got=%s", want, ident.TokenLiteral())
	}
}

func assertLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	t.Helper()
	switch v := expected.(type) {
	case int:
		assertIntegerLiteral(t, exp, int64(v))
	case int64:
		assertIntegerLiteral(t, exp, v)
	case string:
		assertIdentifier(t, exp, v)
	default:
		t.Fatalf("type of exp not handled, got=%T", exp)
	}
}

func assertInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
	}

	assertLiteralExpression(t, opExp.Left, left)
	if opExp.Operator != operator {
		t.Errorf("exp operator is not %s, got=%q", operator, opExp.Operator)
	}
	assertLiteralExpression(t, opExp.Right, right)
}
