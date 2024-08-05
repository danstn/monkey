package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
	"strings"
)

// operator precedence
const (
	_ int = iota
	LOWEST
	EQUALS  // ==
	LTGT    // < or >
	SUM     // +
	PRODUCT // *
	PREFIX  // -X or !X
	CALL    // someFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:    EQUALS,
	token.NEQ:   EQUALS,
	token.LT:    LTGT,
	token.GT:    LTGT,
	token.PLUS:  SUM,
	token.MINUS: SUM,
	token.STAR:  PRODUCT,
	token.SLASH: PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	// the argument is "left side" of the infix operator being parsed
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors   []string // TODO: extend to add row/col
	progress []string // literal progress of what is being parsed at the moment

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// New creates a new parser given an initialised lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.NAME, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// read two tokens, so currToken and peekToken are both set
	p.advance()
	p.advance()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Progress() string {
	return strings.Join(p.progress, " ")
}

func (p *Parser) advance() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
	p.progress = append(p.progress, p.currToken.Literal)
}

func (p *Parser) flushProgress() {
	p.progress = []string{}
}

// Top Level Parsers
// -----------------------------------------------------------------------------

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			p.flushProgress()
			program.Statements = append(program.Statements, stmt)
		}
		p.advance()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Specific Parsers
// -----------------------------------------------------------------------------

// let x = 5;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	// ensure next token is identifier and advance
	if !p.advanceIfPeekIs(token.NAME) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	// ensure next token is '=' and advance
	if !p.advanceIfPeekIs(token.ASSIGN) {
		return nil
	}

	// TODO: implement expression parsing
	for !p.currTokenIs(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

// return 5;
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.advance()

	// TODO: implement expression parsing
	for !p.currTokenIs(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// we want expression statements to have optional semicolons, which makes
	// it easier to type in REPL
	if p.peekTokenIs(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseExpression(_ int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.advance()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// Helpers
// -----------------------------------------------------------------------------

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be '%s', got '%s' parsing: '%s ...'", t, p.peekToken.Type, p.Progress())
	p.errors = append(p.errors, msg)
}

func (p *Parser) advanceIfPeekIs(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.advance()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if precedence, ok := precedences[p.currToken.Type]; ok {
		return precedence
	}
	return LOWEST
}
