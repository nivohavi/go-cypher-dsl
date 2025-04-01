package expr

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// AliasExpression represents an expression with an alias (e.g., n AS person)
type AliasExpression struct {
	Expression core.Expression
	Alias      string
}

// Accept implements the Expression interface
func (a *AliasExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(a)
}

// String returns a string representation of this alias expression
func (a *AliasExpression) String() string {
	return fmt.Sprintf("%s AS %s", a.Expression.String(), a.Alias)
}

// And creates a logical AND with another expression
func (a *AliasExpression) And(other core.Expression) core.Expression {
	return And(a, other)
}

// Or creates a logical OR with another expression
func (a *AliasExpression) Or(other core.Expression) core.Expression {
	return Or(a, other)
}

// Not creates a logical NOT of this expression
func (a *AliasExpression) Not() core.Expression {
	return Not(a)
}

// As creates an alias expression
func As(expression core.Expression, alias string) core.Expression {
	return &AliasExpression{
		Expression: expression,
		Alias:      alias,
	}
}
