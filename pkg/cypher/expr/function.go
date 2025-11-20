package expr

import (
	"fmt"
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

// As creates an alias for this function expression
func (f *FunctionExpression) As(alias string) core.Expression {
	return As(f, alias)
}

// DistinctExpression represents an expression wrapped with DISTINCT (e.g., DISTINCT n)
type DistinctExpression struct {
	Expression core.Expression
}

// Accept implements the Expression interface
func (d *DistinctExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(d)
}

// String returns a string representation of this distinct expression
func (d *DistinctExpression) String() string {
	return "DISTINCT " + d.Expression.String()
}

// And creates a logical AND with another expression
func (d *DistinctExpression) And(other core.Expression) core.Expression {
	return And(d, other)
}

// Or creates a logical OR with another expression
func (d *DistinctExpression) Or(other core.Expression) core.Expression {
	return Or(d, other)
}

// Not creates a logical NOT of this expression
func (d *DistinctExpression) Not() core.Expression {
	return Not(d)
}

// Distinct wraps an expression with DISTINCT keyword
func Distinct(expr core.Expression) core.Expression {
	return &DistinctExpression{Expression: expr}
}

// Distinct wraps this function expression's first argument with DISTINCT
// This is useful for count(DISTINCT expr) patterns
func (f *FunctionExpression) Distinct() core.Expression {
	if len(f.Arguments) == 0 {
		return f
	}
	// Create a new function with the first argument wrapped in DISTINCT
	args := make([]core.Expression, len(f.Arguments))
	args[0] = Distinct(f.Arguments[0])
	for i := 1; i < len(f.Arguments); i++ {
		args[i] = f.Arguments[i]
	}
	return &FunctionExpression{
		Name:      f.Name,
		Arguments: args,
	}
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

// Collect creates a COLLECT function expression
func Collect(expr core.Expression) core.Expression {
	return Function("collect", expr)
}

// BinaryExpression represents a binary operation (e.g., a + b)
type BinaryExpression struct {
	Left     core.Expression
	Right    core.Expression
	Operator string
}

// Accept implements the Expression interface
func (b *BinaryExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(b)
}

// String returns a string representation of this binary expression
func (b *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Operator, b.Right.String())
}

// And creates a logical AND with another expression
func (b *BinaryExpression) And(other core.Expression) core.Expression {
	return And(b, other)
}

// Or creates a logical OR with another expression
func (b *BinaryExpression) Or(other core.Expression) core.Expression {
	return Or(b, other)
}

// Not creates a logical NOT of this expression
func (b *BinaryExpression) Not() core.Expression {
	return Not(b)
}

// Concat concatenates multiple string expressions using the + operator
// This chains expressions: expr1 + expr2 + expr3 + ...
func Concat(expressions ...core.Expression) core.Expression {
	if len(expressions) == 0 {
		return String("")
	}
	if len(expressions) == 1 {
		return expressions[0]
	}
	// Chain expressions with + operator
	result := &BinaryExpression{
		Left:     expressions[0],
		Right:    expressions[1],
		Operator: "+",
	}
	for i := 2; i < len(expressions); i++ {
		result = &BinaryExpression{
			Left:     result,
			Right:    expressions[i],
			Operator: "+",
		}
	}
	return result
}

// Substring creates a SUBSTRING function expression
// substring(expr, start [, length])
func Substring(expr core.Expression, start core.Expression, length ...core.Expression) core.Expression {
	args := []core.Expression{expr, start}
	if len(length) > 0 {
		args = append(args, length[0])
	}
	return Function("substring", args...)
}

// Replace creates a REPLACE function expression
// replace(expr, search, replace)
func Replace(expr, search, replace core.Expression) core.Expression {
	return Function("replace", expr, search, replace)
}

// Split creates a SPLIT function expression
// split(expr, delimiter)
func Split(expr, delimiter core.Expression) core.Expression {
	return Function("split", expr, delimiter)
}

// ToLower creates a toLower function expression
func ToLower(expr core.Expression) core.Expression {
	return Function("toLower", expr)
}

// ToUpper creates a toUpper function expression
func ToUpper(expr core.Expression) core.Expression {
	return Function("toUpper", expr)
}

// Trim creates a TRIM function expression
func Trim(expr core.Expression) core.Expression {
	return Function("trim", expr)
}

// LTrim creates a lTrim function expression
func LTrim(expr core.Expression) core.Expression {
	return Function("lTrim", expr)
}

// RTrim creates a rTrim function expression
func RTrim(expr core.Expression) core.Expression {
	return Function("rTrim", expr)
}

// RawCypherExpression represents a raw Cypher string that will be inserted as-is
// WARNING: Use with caution to avoid Cypher injection vulnerabilities
type RawCypherExpression struct {
	Cypher string
}

// Accept implements the Expression interface
func (r *RawCypherExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(r)
}

// String returns the raw Cypher string
func (r *RawCypherExpression) String() string {
	return r.Cypher
}

// And creates a logical AND with another expression
func (r *RawCypherExpression) And(other core.Expression) core.Expression {
	return And(r, other)
}

// Or creates a logical OR with another expression
func (r *RawCypherExpression) Or(other core.Expression) core.Expression {
	return Or(r, other)
}

// Not creates a logical NOT of this expression
func (r *RawCypherExpression) Not() core.Expression {
	return Not(r)
}

// RawCypher creates a raw Cypher expression that will be inserted as-is into the query
// WARNING: Use with caution to avoid Cypher injection vulnerabilities.
// Only use this when the DSL doesn't support a specific Cypher feature.
func RawCypher(cypher string) core.Expression {
	return &RawCypherExpression{Cypher: cypher}
}
