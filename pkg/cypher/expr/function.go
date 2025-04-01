package expr

import (
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// FunctionExpression represents a function call in Cypher (e.g., count(*))
type FunctionExpression struct {
	Name      string
	Arguments []core.Expression
}

// Accept implements the Expression interface
func (f *FunctionExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(f)
}

// String returns a string representation of this function expression
func (f *FunctionExpression) String() string {
	var sb strings.Builder
	sb.WriteString(f.Name)
	sb.WriteString("(")
	for i, arg := range f.Arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.String())
	}
	sb.WriteString(")")
	return sb.String()
}

// And creates a logical AND with another expression
func (f *FunctionExpression) And(other core.Expression) core.Expression {
	return And(f, other)
}

// Or creates a logical OR with another expression
func (f *FunctionExpression) Or(other core.Expression) core.Expression {
	return Or(f, other)
}

// Not creates a logical NOT of this expression
func (f *FunctionExpression) Not() core.Expression {
	return Not(f)
}

// Function creates a new function expression
func Function(name string, args ...core.Expression) core.Expression {
	return &FunctionExpression{
		Name:      name,
		Arguments: args,
	}
}

// Count creates a COUNT function expression
func Count(expr core.Expression) core.Expression {
	return Function("count", expr)
}

// CountStar creates a COUNT(*) function expression
func CountStar() core.Expression {
	return &FunctionExpression{
		Name:      "count",
		Arguments: []core.Expression{&Literal{Value: "*"}},
	}
}

// Sum creates a SUM function expression
func Sum(expr core.Expression) core.Expression {
	return Function("sum", expr)
}

// Avg creates an AVG function expression
func Avg(expr core.Expression) core.Expression {
	return Function("avg", expr)
}

// Min creates a MIN function expression
func Min(expr core.Expression) core.Expression {
	return Function("min", expr)
}

// Max creates a MAX function expression
func Max(expr core.Expression) core.Expression {
	return Function("max", expr)
}
