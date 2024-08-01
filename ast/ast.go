package ast

import (
	"bytes"
	"monkey/token"
)

// Node
// -----------------------------------------------------------------------------

type Node interface {
	// TokenLiteral is used for debugging and testing.
	TokenLiteral() string
	// String will allow us to print AST notes for debugging
	String() string
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

// Program
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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Let Statement
// -----------------------------------------------------------------------------

type LetStatement struct {
	Token token.Token // token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// Identifier
// -----------------------------------------------------------------------------

type Identifier struct {
	Token token.Token // token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Return Statement
// -----------------------------------------------------------------------------

type ReturnStatement struct {
	Token       token.Token // token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// Expression Statement
// -----------------------------------------------------------------------------

// ExpressionStatement is a wrapper for expressions so that they could be added
// to the Statements slice in Program
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Integer Literal Expression
// -----------------------------------------------------------------------------

type IntegerLiteral struct {
	Token token.Token // 5
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal } // "5"
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
