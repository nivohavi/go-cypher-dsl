package expr

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// Literal represents a literal value in Cypher
type Literal struct {
	Value any
}

// Accept implements the visitor pattern for the Expression interface
func (l *Literal) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(l)
}

// String returns a string representation of this literal
func (l *Literal) String() string {
	return formatValue(l.Value)
}

// And creates a logical AND with another expression
func (l *Literal) And(other core.Expression) core.Expression {
	return And(l, other)
}

// Or creates a logical OR with another expression
func (l *Literal) Or(other core.Expression) core.Expression {
	return Or(l, other)
}

// Not creates a logical NOT of this expression
func (l *Literal) Not() core.Expression {
	return Not(l)
}

// Var represents a variable (symbolic name) in Cypher
type Var struct {
	Name string
}

// Accept implements the visitor pattern for the Expression interface
func (v *Var) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(v)
}

// String returns the name of this variable
func (v *Var) String() string {
	return v.Name
}

// And creates a logical AND with another expression
func (v *Var) And(other core.Expression) core.Expression {
	return And(v, other)
}

// Or creates a logical OR with another expression
func (v *Var) Or(other core.Expression) core.Expression {
	return Or(v, other)
}

// Not creates a logical NOT of this expression
func (v *Var) Not() core.Expression {
	return Not(v)
}

// SymbolicName returns the symbolic name of this variable
func (v *Var) SymbolicName() string {
	return v.Name
}

// Property returns a property access expression
func (v *Var) Property(propertyName string) core.PropertyExpression {
	return NewProperty(v, propertyName)
}

// BaseExpression provides a default implementation of the Expression interface
type BaseExpression struct{}

// Accept implements the visitor pattern for the Expression interface
func (b *BaseExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(b)
}

// String returns a default string representation
func (b *BaseExpression) String() string {
	return ""
}

// And creates a logical AND with another expression
func (b *BaseExpression) And(other core.Expression) core.Expression {
	return And(b, other)
}

// Or creates a logical OR with another expression
func (b *BaseExpression) Or(other core.Expression) core.Expression {
	return Or(b, other)
}

// Not creates a logical NOT of this expression
func (b *BaseExpression) Not() core.Expression {
	return Not(b)
}
