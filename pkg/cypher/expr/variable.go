package expr

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// VariableExpression represents a variable reference in Cypher
type VariableExpression struct {
	name string
}

// NewVariableExpression creates a new variable expression
func NewVariableExpression(name string) *VariableExpression {
	return &VariableExpression{
		name: name,
	}
}

// Accept implements the Expression interface
func (v *VariableExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(v)
}

// String returns the string representation of the variable
func (v *VariableExpression) String() string {
	return v.name
}

// Name returns the name of the variable
func (v *VariableExpression) Name() string {
	return v.name
}

// And creates a logical AND with another expression
func (v *VariableExpression) And(other core.Expression) core.Expression {
	return And(v, other)
}

// Or creates a logical OR with another expression
func (v *VariableExpression) Or(other core.Expression) core.Expression {
	return Or(v, other)
}

// Not creates a logical NOT of this expression
func (v *VariableExpression) Not() core.Expression {
	return Not(v)
}
