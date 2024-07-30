package ast

import "monkey/token"

// NODE
// -----------------------------------------------------------------------------

type Node interface {
	// TokenLiteral is used for debugging and testing.
	TokenLiteral() string
}

// Statement is an identifier and an expression. For example:
//
//	let x = 5;
//	let y = add(2, 2) * 5 / 10;
//	return 5;
type Statement interface {
	Node
	statementNode()
}

// Expression produces a value. For example:
//
//	5
//	add(5, 5) * 5 / 10
//	fn(x, y) { return x }
type Expression interface {
	Node
	expressionNode()
}

// IMPLEMENTATION
// -----------------------------------------------------------------------------

// Program is the root node of every AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }