package cypher

import (
	"fmt"
	"strings"
)

// FunctionExpression represents a function call in Cypher
type functionExpression struct {
	name      string
	arguments []Expression
}

// Accept implements the Expression interface
func (f *functionExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(f)
}

// String returns a string representation of the function call
func (f *functionExpression) String() string {
	args := make([]string, len(f.arguments))
	for i, arg := range f.arguments {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", f.name, strings.Join(args, ", "))
}

// FunctionName returns the name of the function
func (f *functionExpression) FunctionName() string {
	return f.name
}

// Arguments returns the arguments of the function
func (f *functionExpression) Arguments() []Expression {
	return f.arguments
}

// Function creates a new function expression
func Function(name string, args ...Expression) Expression {
	return &functionExpression{
		name:      name,
		arguments: args,
	}
}

// existsExpression represents an EXISTS operation in Cypher
type existsExpression struct {
	expression Expression
}

// Accept implements the Expression interface
func (e *existsExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(e)
}

// String returns a string representation of the EXISTS expression
func (e *existsExpression) String() string {
	return fmt.Sprintf("EXISTS(%s)", e.expression)
}

// Exists creates an EXISTS expression
func Exists(expression Expression) BooleanExpression {
	return &existsExpression{expression: expression}
}

// And creates a boolean AND expression
func (e *existsExpression) And(other BooleanExpression) BooleanExpression {
	return &binaryExpression{
		left:     e,
		right:    other,
		operator: "AND",
	}
}

// Or creates a boolean OR expression
func (e *existsExpression) Or(other BooleanExpression) BooleanExpression {
	return &binaryExpression{
		left:     e,
		right:    other,
		operator: "OR",
	}
}

// Not negates this boolean expression
func (e *existsExpression) Not() BooleanExpression {
	return &notExpression{expression: e}
}

// IsBooleanExpression marks this as a boolean expression
func (e *existsExpression) IsBooleanExpression() {}

// Common Cypher function helpers

// Count creates a COUNT function call
func Count(expression Expression) Expression {
	return Function("count", expression)
}

// CountDistinct creates a COUNT(DISTINCT ...) function call
func CountDistinct(expression Expression) Expression {
	// TODO: Implement DISTINCT properly
	return Function("count", expression)
}

// Collect creates a COLLECT function call
func Collect(expression Expression) Expression {
	return Function("collect", expression)
}

// Max creates a MAX function call
func Max(expression Expression) Expression {
	return Function("max", expression)
}

// Min creates a MIN function call
func Min(expression Expression) Expression {
	return Function("min", expression)
}

// Avg creates an AVG function call
func Avg(expression Expression) Expression {
	return Function("avg", expression)
}

// Sum creates a SUM function call
func Sum(expression Expression) Expression {
	return Function("sum", expression)
}

// Size creates a SIZE function call
func Size(expression Expression) Expression {
	return Function("size", expression)
}

// Type creates a TYPE function call
func Type(expression Expression) Expression {
	return Function("type", expression)
}

// Id creates an ID function call
func Id(expression Expression) Expression {
	return Function("id", expression)
}

// Trim creates a TRIM function call
func Trim(expression Expression) Expression {
	return Function("trim", expression)
}

// ToLower creates a toLower function call
func ToLower(expression Expression) Expression {
	return Function("toLower", expression)
}

// ToUpper creates a toUpper function call
func ToUpper(expression Expression) Expression {
	return Function("toUpper", expression)
}

// Substring creates a substring function call
func Substring(expression Expression, start, length Expression) Expression {
	return Function("substring", expression, start, length)
}

// Replace creates a replace function call
func Replace(original, search, replace Expression) Expression {
	return Function("replace", original, search, replace)
}

// Split creates a split function call
func Split(expression, delimiter Expression) Expression {
	return Function("split", expression, delimiter)
}

// Coalesce creates a COALESCE function call
func Coalesce(expressions ...Expression) Expression {
	return Function("coalesce", expressions...)
}
