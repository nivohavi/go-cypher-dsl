package expr

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// LogicalOperator represents a logical operator (AND, OR, XOR)
type LogicalOperator string

const (
	// AND represents logical AND
	AND LogicalOperator = "AND"
	// OR represents logical OR
	OR LogicalOperator = "OR"
	// XOR represents logical XOR
	XOR LogicalOperator = "XOR"
)

// LogicalExpression represents a logical expression (e.g., a AND b)
type LogicalExpression struct {
	left     core.Expression
	right    core.Expression
	operator string
}

// Left returns the left side of the logical expression
func (l *LogicalExpression) Left() core.Expression {
	return l.left
}

// Right returns the right side of the logical expression
func (l *LogicalExpression) Right() core.Expression {
	return l.right
}

// Accept implements the Expression interface
func (l *LogicalExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(l)
}

// String returns a string representation of this logical expression
func (l *LogicalExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", l.left.String(), l.operator, l.right.String())
}

// And creates a logical AND with another expression
func (l *LogicalExpression) And(other core.Expression) core.Expression {
	return And(l, other)
}

// Or creates a logical OR with another expression
func (l *LogicalExpression) Or(other core.Expression) core.Expression {
	return Or(l, other)
}

// Xor creates a logical XOR with another expression
func (l *LogicalExpression) Xor(other core.Expression) core.Expression {
	return Xor(l, other)
}

// Expressions returns all expressions contained in this logical expression
func (l *LogicalExpression) Expressions() []core.Expression {
	return []core.Expression{l.left, l.right}
}

// And creates a logical AND expression
func And(left, right core.Expression) core.Expression {
	return &LogicalExpression{
		left:     left,
		right:    right,
		operator: "AND",
	}
}

// Or creates a logical OR expression
func Or(left, right core.Expression) core.Expression {
	return &LogicalExpression{
		left:     left,
		right:    right,
		operator: string(OR),
	}
}

// Xor creates a logical XOR expression
func Xor(left, right core.Expression) core.Expression {
	return &LogicalExpression{
		left:     left,
		right:    right,
		operator: "XOR",
	}
}

// NotExpression represents a logical NOT expression (e.g., NOT a)
type NotExpression struct {
	expr core.Expression
}

// Accept implements the Expression interface
func (n *NotExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(n)
}

// String returns a string representation of this NOT expression
func (n *NotExpression) String() string {
	return fmt.Sprintf("NOT %s", n.expr.String())
}

// And creates a logical AND with another expression
func (n *NotExpression) And(other core.Expression) core.Expression {
	return And(n, other)
}

// Or creates a logical OR with another expression
func (n *NotExpression) Or(other core.Expression) core.Expression {
	return Or(n, other)
}

// Expressions returns the expressions in this NOT expression
func (n *NotExpression) Expressions() []core.Expression {
	return []core.Expression{n.expr}
}

// Not creates a logical NOT expression
func Not(expr core.Expression) core.Expression {
	return &NotExpression{
		expr: expr,
	}
}

// Not creates a logical NOT of this expression
func (l *LogicalExpression) Not() core.Expression {
	return Not(l)
}

// Not creates a logical NOT of this expression
func (n *NotExpression) Not() core.Expression {
	// Double negation cancels out
	return n.expr
}
