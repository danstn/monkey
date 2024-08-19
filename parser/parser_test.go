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

	assertParserNoErrors(t, p)
	assertProgramNotNil(t, program)
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
	assertParserNoErrors(t, p)
	assertProgramNotNil(t, program)
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

	assertParserNoErrors(t, p)
	assertProgramNotNil(t, program)
	assertProgramStatements(t, program, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParserNoErrors(t, p)
	assertProgramNotNil(t, program)
	assertProgramStatements(t, program, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression is not &ast.IntegerLiteral, got=%T", stmt.Expression)
	}

	test.AssertEqual(t, literal.Value, 5)
	test.AssertEqual(t, literal.TokenLiteral(), "5")
}

func TestBoolLiteralExpression(t *testing.T) {
	input := "true; false;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParserNoErrors(t, p)
	assertProgramNotNil(t, program)
	assertProgramStatements(t, program, 2)

	stmt1 := assertExpressionStatement(t, program.Statements[0])
	lit1, ok := stmt1.Expression.(*ast.BoolLiteral)
	if !ok {
		t.Fatalf("expression is not &ast.BoolLiteral, got=%T", stmt1.Expression)
	}
	test.AssertEqual(t, lit1.Value, true)
	test.AssertEqual(t, lit1.TokenLiteral(), "true")

	stmt2 := assertExpressionStatement(t, program.Statements[1])
	lit2, ok := stmt2.Expression.(*ast.BoolLiteral)
	if !ok {
		t.Fatalf("expression is not &ast.BoolLiteral, got=%T", stmt1.Expression)
	}
	test.AssertEqual(t, lit2.Value, false)
	test.AssertEqual(t, lit2.TokenLiteral(), "false")
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertParserNoErrors(t, p)
		assertProgramNotNil(t, program)
		assertProgramStatements(t, program, 1)

		stmt := assertExpressionStatement(t, program.Statements[0])

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not at.PrefixExpression, got=%T", stmt.Expression)
		}
		if op := exp.Operator; op != tt.operator {
			t.Fatalf("exp.Operator is not '%s', got=%s", tt.operator, op)
		}
		assertIntegerLiteral(t, exp.Right, tt.integerValue)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		p := New(lexer.New(tt.input))
		program := p.ParseProgram()
		assertParserNoErrors(t, p)
		assertProgramNotNil(t, program)
		assertProgramStatements(t, program, 1)

		stmt := assertExpressionStatement(t, program.Statements[0])

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression, got %T", stmt.Expression)
		}

		assertInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c ",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (1 + 3) + 4",
			"((1 + (1 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertParserNoErrors(t, p)

		if got := program.String(); got != tt.want {
			t.Errorf("#%d want=%q, got=%q", i, tt.want, got)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParserNoErrors(t, p)
	assertProgramNotNil(t, program)
	assertProgramStatements(t, program, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp is not ast.IfExpression, got %T", stmt.Expression)
	}

	assertInfixExpression(t, exp.Condition, "x", "<", "y")
	if n := len(exp.Consequence.Statements); n != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", n)
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence statement is not ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	assertIdentifier(t, consequence.Expression, "x")

	test.AssertEqual(t, exp.Alternative, nil)

}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParserNoErrors(t, p)
	assertProgramNotNil(t, program)
	assertProgramStatements(t, program, 1)

	stmt := assertExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp is not ast.IfExpression, got %T", stmt.Expression)
	}

	assertInfixExpression(t, exp.Condition, "x", "<", "y")
	if n := len(exp.Consequence.Statements); n != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", n)
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence statement is not ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	assertIdentifier(t, consequence.Expression, "x")

	test.AssertNotEqual(t, exp.Alternative, nil)

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative statement is not ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}

	assertIdentifier(t, alternative.Expression, "y")
}
