package ast

import (
	"bytes"
	"monkey/token"
	"strings"
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

// Boolean Literal Expression
// -----------------------------------------------------------------------------

type BoolLiteral struct {
	Token token.Token // true
	Value bool
}

func (bl *BoolLiteral) expressionNode()      {}
func (bl *BoolLiteral) TokenLiteral() string { return bl.Token.Literal } // "true"
func (bl *BoolLiteral) String() string       { return bl.Token.Literal }

// Prefix Expression
// -----------------------------------------------------------------------------

type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// Infix Expression
// -----------------------------------------------------------------------------

type InfixExpression struct {
	Token    token.Token // the prefix token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// Block Statement
// -----------------------------------------------------------------------------

type BlockStatement struct {
	Token      token.Token // The '{' token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// If Expression
// -----------------------------------------------------------------------------

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// Function Literal Expression
// -----------------------------------------------------------------------------

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	// fn
	out.WriteString("fn")

	// (x, y, z)
	var params []string
	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}
	out.WriteString(token.LPAREN)
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(token.RPAREN)

	// { .. }
	out.WriteString(fl.Body.String())

	return out.String()
}

// Call Expression
// -----------------------------------------------------------------------------

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString(token.LPAREN)
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(token.RPAREN)

	return out.String()
}
