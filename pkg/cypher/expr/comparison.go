package expr

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// ComparisonOperator represents the type of comparison
type ComparisonOperator string

const (
	// EQ represents equals (=)
	EQ ComparisonOperator = "="
	// NE represents not equals (<>)
	NE ComparisonOperator = "<>"
	// GT represents greater than (>)
	GT ComparisonOperator = ">"
	// LT represents less than (<)
	LT ComparisonOperator = "<"
	// GTE represents greater than or equal (>=)
	GTE ComparisonOperator = ">="
	// LTE represents less than or equal (<=)
	LTE ComparisonOperator = "<="
)

// ComparisonExpression represents a comparison between two expressions
type ComparisonExpression struct {
	left     core.Expression
	right    core.Expression
	operator string
}

// Left returns the left side of the comparison
func (c *ComparisonExpression) Left() core.Expression {
	return c.left
}

// Right returns the right side of the comparison
func (c *ComparisonExpression) Right() core.Expression {
	return c.right
}

// Accept implements the Expression interface
func (c *ComparisonExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(c)
}

// String returns a string representation of this comparison
func (c *ComparisonExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", c.left.String(), c.operator, c.right.String())
}

// And creates a logical AND with another expression
func (c *ComparisonExpression) And(other core.Expression) core.Expression {
	return And(c, other)
}

// Or creates a logical OR with another expression
func (c *ComparisonExpression) Or(other core.Expression) core.Expression {
	return Or(c, other)
}

// Not creates a logical NOT of this expression
func (c *ComparisonExpression) Not() core.Expression {
	return Not(c)
}

// Xor creates a logical XOR with another expression
func (c *ComparisonExpression) Xor(other core.Expression) core.Expression {
	return Xor(c, other)
}

// Equals creates an equality comparison
func Equals(left, right core.Expression) core.Expression {
	return &ComparisonExpression{
		left:     left,
		right:    right,
		operator: string(EQ),
	}
}

// NotEquals creates a not-equal comparison
func NotEquals(left, right core.Expression) core.Expression {
	return &ComparisonExpression{
		left:     left,
		right:    right,
		operator: string(NE),
	}
}

// GreaterThan creates a greater-than comparison
func GreaterThan(left, right core.Expression) core.Expression {
	return &ComparisonExpression{
		left:     left,
		right:    right,
		operator: string(GT),
	}
}

// LessThan creates a less-than comparison
func LessThan(left, right core.Expression) core.Expression {
	return &ComparisonExpression{
		left:     left,
		right:    right,
		operator: string(LT),
	}
}

// GreaterThanEqual creates a greater-than-or-equal comparison
func GreaterThanEqual(left, right core.Expression) core.Expression {
	return &ComparisonExpression{
		left:     left,
		right:    right,
		operator: string(GTE),
	}
}

// LessThanEqual creates a less-than-or-equal comparison
func LessThanEqual(left, right core.Expression) core.Expression {
	return &ComparisonExpression{
		left:     left,
		right:    right,
		operator: string(LTE),
	}
}

// IsNull creates a null check
func IsNull(expr core.Expression) core.Expression {
	return &ComparisonExpression{
		left:     expr,
		right:    &NullLiteral{},
		operator: "IS",
	}
}

// IsNotNull creates a not-null check
func IsNotNull(expr core.Expression) core.Expression {
	return &ComparisonExpression{
		left:     expr,
		right:    &NullLiteral{},
		operator: "IS NOT",
	}
}

// In creates an IN comparison
func In(expr core.Expression, values ...any) core.Expression {
	var elements []core.Expression
	for _, v := range values {
		if expr, ok := v.(core.Expression); ok {
			elements = append(elements, expr)
		} else {
			elements = append(elements, LiteralFromValue(v))
		}
	}
	list := &ListExpression{Elements: elements}
	return &ComparisonExpression{
		left:     expr,
		right:    list,
		operator: "IN",
	}
}

// Contains creates a CONTAINS comparison
func Contains(expr core.Expression, value string) core.Expression {
	return &ComparisonExpression{
		left:     expr,
		right:    String(value),
		operator: "CONTAINS",
	}
}

// StartsWith creates a STARTS WITH comparison
func StartsWith(expr core.Expression, value string) core.Expression {
	return &ComparisonExpression{
		left:     expr,
		right:    String(value),
		operator: "STARTS WITH",
	}
}

// EndsWith creates an ENDS WITH comparison
func EndsWith(expr core.Expression, value string) core.Expression {
	return &ComparisonExpression{
		left:     expr,
		right:    String(value),
		operator: "ENDS WITH",
	}
}

// RegularExpression creates a =~ comparison with a regular expression
func RegularExpression(expr core.Expression, pattern string) core.Expression {
	return &ComparisonExpression{
		left:     expr,
		right:    String(pattern),
		operator: "=~",
	}
}
