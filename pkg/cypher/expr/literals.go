// Package expr provides expression implementations for Cypher DSL
package expr

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// BooleanLiteral represents a boolean literal (true/false)
type BooleanLiteral struct {
	Value bool
}

// Accept implements the Expression interface
func (b *BooleanLiteral) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(b)
}

// String returns a string representation of this boolean literal
func (b *BooleanLiteral) String() string {
	return strconv.FormatBool(b.Value)
}

// And creates a logical AND with another expression
func (b *BooleanLiteral) And(other core.Expression) core.Expression {
	return And(b, other)
}

// Or creates a logical OR with another expression
func (b *BooleanLiteral) Or(other core.Expression) core.Expression {
	return Or(b, other)
}

// Not creates a logical NOT of this expression
func (b *BooleanLiteral) Not() core.Expression {
	return Not(b)
}

// IntegerLiteral represents an integer literal
type IntegerLiteral struct {
	Value int64
}

// Accept implements the Expression interface
func (i *IntegerLiteral) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(i)
}

// String returns a string representation of this integer literal
func (i *IntegerLiteral) String() string {
	return strconv.FormatInt(i.Value, 10)
}

// And creates a logical AND with another expression
func (i *IntegerLiteral) And(other core.Expression) core.Expression {
	return And(i, other)
}

// Or creates a logical OR with another expression
func (i *IntegerLiteral) Or(other core.Expression) core.Expression {
	return Or(i, other)
}

// Not creates a logical NOT of this expression
func (i *IntegerLiteral) Not() core.Expression {
	return Not(i)
}

// FloatLiteral represents a floating-point literal
type FloatLiteral struct {
	Value float64
}

// Accept implements the Expression interface
func (f *FloatLiteral) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(f)
}

// String returns a string representation of this float literal
func (f *FloatLiteral) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

// And creates a logical AND with another expression
func (f *FloatLiteral) And(other core.Expression) core.Expression {
	return And(f, other)
}

// Or creates a logical OR with another expression
func (f *FloatLiteral) Or(other core.Expression) core.Expression {
	return Or(f, other)
}

// Not creates a logical NOT of this expression
func (f *FloatLiteral) Not() core.Expression {
	return Not(f)
}

// NullLiteral represents a null value
type NullLiteral struct{}

// Accept implements the Expression interface
func (n *NullLiteral) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(n)
}

// String returns a string representation of this null literal
func (n *NullLiteral) String() string {
	return "NULL"
}

// And creates a logical AND with another expression
func (n *NullLiteral) And(other core.Expression) core.Expression {
	return And(n, other)
}

// Or creates a logical OR with another expression
func (n *NullLiteral) Or(other core.Expression) core.Expression {
	return Or(n, other)
}

// Not creates a logical NOT of this expression
func (n *NullLiteral) Not() core.Expression {
	return Not(n)
}

// StringLiteral represents a string literal value
type StringLiteral struct {
	Value string
}

// Accept implements the Expression interface
func (s *StringLiteral) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(s)
}

// String returns a string representation of this string literal
func (s *StringLiteral) String() string {
	escaped := strings.ReplaceAll(s.Value, "'", "\\'")
	return fmt.Sprintf("'%s'", escaped)
}

// And creates a logical AND with another expression
func (s *StringLiteral) And(other core.Expression) core.Expression {
	return And(s, other)
}

// Or creates a logical OR with another expression
func (s *StringLiteral) Or(other core.Expression) core.Expression {
	return Or(s, other)
}

// Not creates a logical NOT of this expression
func (s *StringLiteral) Not() core.Expression {
	return Not(s)
}

// ListExpression represents a list literal expression (e.g., [1, 2, 3])
type ListExpression struct {
	Elements []core.Expression
}

// Accept implements the Expression interface
func (l *ListExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(l)
}

// String returns a string representation of this list expression
func (l *ListExpression) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, element := range l.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(element.String())
	}
	sb.WriteString("]")
	return sb.String()
}

// And creates a logical AND with another expression
func (l *ListExpression) And(other core.Expression) core.Expression {
	return And(l, other)
}

// Or creates a logical OR with another expression
func (l *ListExpression) Or(other core.Expression) core.Expression {
	return Or(l, other)
}

// Not creates a logical NOT of this expression
func (l *ListExpression) Not() core.Expression {
	return Not(l)
}

// MapLiteralExpression represents a map literal expression (e.g., {name: 'John', age: 30})
type MapLiteralExpression struct {
	Entries map[string]core.Expression
}

// Accept implements the Expression interface
func (m *MapLiteralExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(m)
}

// String returns a string representation of this map literal expression
func (m *MapLiteralExpression) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	i := 0
	for key, value := range m.Entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(value.String())
		i++
	}
	sb.WriteString("}")
	return sb.String()
}

// And creates a logical AND with another expression
func (m *MapLiteralExpression) And(other core.Expression) core.Expression {
	return And(m, other)
}

// Or creates a logical OR with another expression
func (m *MapLiteralExpression) Or(other core.Expression) core.Expression {
	return Or(m, other)
}

// Not creates a logical NOT of this expression
func (m *MapLiteralExpression) Not() core.Expression {
	return Not(m)
}

// ParameterExpression represents a parameterized value (e.g., $name)
type ParameterExpression struct {
	Name  string
	Value any
}

// Accept implements the Expression interface
func (p *ParameterExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(p)
}

// String returns a string representation of this parameter expression
func (p *ParameterExpression) String() string {
	return "$" + p.Name
}

// And creates a logical AND with another expression
func (p *ParameterExpression) And(other core.Expression) core.Expression {
	return And(p, other)
}

// Or creates a logical OR with another expression
func (p *ParameterExpression) Or(other core.Expression) core.Expression {
	return Or(p, other)
}

// Not creates a logical NOT of this expression
func (p *ParameterExpression) Not() core.Expression {
	return Not(p)
}

// String creates a string literal
func String(value string) core.Expression {
	return &StringLiteral{Value: value}
}

// Integer creates an integer literal
func Integer(value int64) core.Expression {
	return &IntegerLiteral{Value: value}
}

// Float creates a float literal
func Float(value float64) core.Expression {
	return &FloatLiteral{Value: value}
}

// Boolean creates a boolean literal
func Boolean(value bool) core.Expression {
	return &BooleanLiteral{Value: value}
}

// Null creates a null literal
func Null() core.Expression {
	return &NullLiteral{}
}

// List creates a list expression
func List(elements ...core.Expression) core.Expression {
	return &ListExpression{Elements: elements}
}

// Map creates a map expression
func Map(entries map[string]core.Expression) core.Expression {
	return &MapLiteralExpression{Entries: entries}
}

// Param creates a parameter expression
func Param(name string, value any) core.Expression {
	return &ParameterExpression{
		Name:  name,
		Value: value,
	}
}

// LiteralFromValue converts a Go value to an Expression
func LiteralFromValue(value any) core.Expression {
	if value == nil {
		return Null()
	}

	switch v := value.(type) {
	case core.Expression:
		return v
	case string:
		return String(v)
	case int:
		return Integer(int64(v))
	case int64:
		return Integer(v)
	case float64:
		return Float(v)
	case bool:
		return Boolean(v)
	case []any:
		elements := make([]core.Expression, len(v))
		for i, elem := range v {
			elements[i] = LiteralFromValue(elem)
		}
		return List(elements...)
	case map[string]any:
		result := make(map[string]core.Expression)
		for k, val := range v {
			result[k] = LiteralFromValue(val)
		}
		return Map(result)
	default:
		// For unknown types, convert to string representation
		return String(fmt.Sprintf("%v", v))
	}
}
