package expr

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// OrderByExpression represents an expression with ordering direction
type OrderByExpression struct {
	Expression core.Expression
	Descending bool
}

// Accept implements the Expression interface
func (o *OrderByExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(o)
}

// String returns a string representation of this order by expression
func (o *OrderByExpression) String() string {
	if o.Descending {
		return fmt.Sprintf("%s DESC", o.Expression.String())
	}
	return fmt.Sprintf("%s ASC", o.Expression.String())
}

// And creates a logical AND with another expression
func (o *OrderByExpression) And(other core.Expression) core.Expression {
	return And(o, other)
}

// Or creates a logical OR with another expression
func (o *OrderByExpression) Or(other core.Expression) core.Expression {
	return Or(o, other)
}

// Not creates a logical NOT of this expression
func (o *OrderByExpression) Not() core.Expression {
	return Not(o)
}

// Desc creates a descending order by expression
func Desc(expression core.Expression) core.Expression {
	return &OrderByExpression{
		Expression: expression,
		Descending: true,
	}
}

// Asc creates a new ascending order by expression
func Asc(expr core.Expression) core.Expression {
	return &OrderByExpression{
		Expression: expr,
		Descending: false,
	}
}
