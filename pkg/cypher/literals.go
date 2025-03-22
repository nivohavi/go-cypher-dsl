package cypher

import (
	"fmt"
	"reflect"
	"strings"
)

// Assert interface implementation
var (
	_ Expression = (*literal)(nil)
	_ Expression = (*parameter)(nil)
)

// literal represents a literal value in Cypher
type literal struct {
	value interface{}
}

// Accept implements the Expression interface
func (l *literal) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(l)
}

// String returns a string representation of the literal
func (l *literal) String() string {
	switch v := l.value.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "\\'"))
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Value returns the underlying value of the literal
func (l *literal) Value() interface{} {
	return l.value
}

// Literal creates a new literal expression
func Literal(value interface{}) Expression {
	return &literal{value: value}
}

// Null creates a null literal expression
func Null() Expression {
	return &literal{value: nil}
}

// parameter represents a parameter in a Cypher query
type parameter struct {
	name  string
	value interface{}
}

// Accept implements the Expression interface
func (p *parameter) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(p)
}

// String returns a string representation of the parameter
func (p *parameter) String() string {
	return fmt.Sprintf("$%s", p.name)
}

// Name returns the parameter name
func (p *parameter) Name() string {
	return p.name
}

// Value returns the parameter value
func (p *parameter) Value() interface{} {
	return p.value
}

// HasValue checks if the parameter has a value
func (p *parameter) HasValue() bool {
	return p.value != nil
}

// Parameter creates a new parameter expression
func Parameter(name string) Expression {
	return &parameter{name: name}
}

// ParamWithValue creates a new parameter expression with a value
func ParamWithValue(name string, value interface{}) Expression {
	return &parameter{name: name, value: value}
}

// nullExpression represents a NULL check operation
type nullExpression struct {
	expression Expression
	isNull     bool
}

// Accept implements the Expression interface
func (n *nullExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(n)
}

// String returns a string representation of this NULL check
func (n *nullExpression) String() string {
	if n.isNull {
		return fmt.Sprintf("%s IS NULL", n.expression.String())
	}
	return fmt.Sprintf("%s IS NOT NULL", n.expression.String())
}

// And combines this expression with another using AND
func (n *nullExpression) And(other BooleanExpression) BooleanExpression {
	return &binaryExpression{
		left:     n,
		right:    other,
		operator: "AND",
	}
}

// Or combines this expression with another using OR
func (n *nullExpression) Or(other BooleanExpression) BooleanExpression {
	return &binaryExpression{
		left:     n,
		right:    other,
		operator: "OR",
	}
}

// Not negates this boolean expression
func (n *nullExpression) Not() BooleanExpression {
	return &notExpression{expression: n}
}

// IsBooleanExpression marks this as a boolean expression
func (n *nullExpression) IsBooleanExpression() {}

// inExpression represents an IN operation in Cypher
type inExpression struct {
	expression Expression
	values     []Expression
}

// Accept implements the Expression interface
func (i *inExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(i)
}

// String returns a string representation of this IN operation
func (i *inExpression) String() string {
	var valuesStrs []string
	for _, v := range i.values {
		valuesStrs = append(valuesStrs, v.String())
	}
	return fmt.Sprintf("%s IN [%s]", i.expression.String(), strings.Join(valuesStrs, ", "))
}

// And combines this expression with another using AND
func (i *inExpression) And(other BooleanExpression) BooleanExpression {
	return &binaryExpression{
		left:     i,
		right:    other,
		operator: "AND",
	}
}

// Or combines this expression with another using OR
func (i *inExpression) Or(other BooleanExpression) BooleanExpression {
	return &binaryExpression{
		left:     i,
		right:    other,
		operator: "OR",
	}
}

// Not negates this boolean expression
func (i *inExpression) Not() BooleanExpression {
	return &notExpression{expression: i}
}

// IsBooleanExpression marks this as a boolean expression
func (i *inExpression) IsBooleanExpression() {}

// binaryExpression represents a binary operation in Cypher
type binaryExpression struct {
	left     Expression
	right    Expression
	operator string
}

// Accept implements the Expression interface
func (b *binaryExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(b)
}

// String returns a string representation of the binary expression
func (b *binaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", b.left, b.operator, b.right)
}

// And creates a boolean AND expression
func (b *binaryExpression) And(other BooleanExpression) BooleanExpression {
	return &binaryExpression{
		left:     b,
		right:    other,
		operator: "AND",
	}
}

// Or creates a boolean OR expression
func (b *binaryExpression) Or(other BooleanExpression) BooleanExpression {
	return &binaryExpression{
		left:     b,
		right:    other,
		operator: "OR",
	}
}

// Not creates a boolean NOT expression
func (b *binaryExpression) Not() BooleanExpression {
	return &notExpression{expression: b}
}

// IsBooleanExpression marks this as a boolean expression
func (b *binaryExpression) IsBooleanExpression() {}

// ListOf creates a list from the provided expressions
func ListOf(elements ...Expression) Expression {
	return &listExpression{elements: elements}
}

// List creates a list from Go values, automatically converting to expressions
func List(values ...interface{}) Expression {
	exprs := make([]Expression, len(values))
	for i, v := range values {
		exprs[i] = toExpression(v)
	}
	return ListOf(exprs...)
}

// listExpression represents a list in Cypher
type listExpression struct {
	elements []Expression
}

// Accept implements the Expression interface
func (l *listExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(l)
}

// String returns a string representation of the list
func (l *listExpression) String() string {
	elemStrs := make([]string, len(l.elements))
	for i, elem := range l.elements {
		elemStrs[i] = elem.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(elemStrs, ", "))
}

// Elements returns the expressions in this list
func (l *listExpression) Elements() []Expression {
	return l.elements
}

// Map creates a new map expression
func Map(entries map[string]Expression) Expression {
	return &mapExpression{entries: entries}
}

// MapOf creates a map from Go values, automatically converting to expressions
func MapOf(entries map[string]interface{}) Expression {
	exprEntries := make(map[string]Expression)
	for k, v := range entries {
		exprEntries[k] = toExpression(v)
	}
	return Map(exprEntries)
}

// mapExpression represents a map in Cypher
type mapExpression struct {
	entries map[string]Expression
}

// Accept implements the Expression interface
func (m *mapExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(m)
}

// String returns a string representation of the map
func (m *mapExpression) String() string {
	entryStrs := make([]string, 0, len(m.entries))
	for k, v := range m.entries {
		entryStrs = append(entryStrs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(entryStrs, ", "))
}

// Entries returns the entries in this map
func (m *mapExpression) Entries() map[string]Expression {
	return m.entries
}

// toExpression converts an arbitrary value to an Expression
func ToExpression(value interface{}) Expression {
	if value == nil {
		return Literal(nil)
	}

	// If it's already an Expression, return it
	if expr, ok := value.(Expression); ok {
		return expr
	}

	// If it's a built-in type, convert to literal
	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return Literal(value)
	}

	// For other types, just use their String() representation
	return Literal(fmt.Sprintf("%v", value))
}

// Kept for backward compatibility - use ToExpression instead
func toExpression(value interface{}) Expression {
	return ToExpression(value)
}
