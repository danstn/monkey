package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
}

// New creates a new parser given an initialised lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// read two tokens, so currToken and peekToken are both set
	p.advance()
	p.advance()

	return p
}

func (p *Parser) advance() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		if stmt := p.parseStatement(); stmt != nil {
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
	default:
		return nil
	}
}

// let x = 5;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	// ensure next token is identifier and advance
	if !p.advanceIfPeekIs(token.IDENT) {
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

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) advanceIfPeekIs(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.advance()
		return true
	} else {
		return false
	}
}
