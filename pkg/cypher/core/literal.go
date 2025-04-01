package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/safety"
)

// LiteralExpression represents a literal value in a Cypher query
type LiteralExpression struct {
	value any
}

// NewLiteral creates a new literal expression
func NewLiteral(value any) *LiteralExpression {
	// Check if this is a string literal and trigger warning if needed
	if strVal, ok := value.(string); ok {
		safety.WarnOnStringLiteral(strVal, "literal expression")
	}

	return &LiteralExpression{
		value: value,
	}
}

// Value returns the value of this literal
func (l *LiteralExpression) Value() any {
	return l.value
}

// String returns the string representation of this literal
func (l *LiteralExpression) String() string {
	if l.value == nil {
		return "NULL"
	}

	switch v := l.value.(type) {
	case string:
		// Escape quotes and wrap in quotes
		escaped := strings.ReplaceAll(v, "'", "\\'")
		return "'" + escaped + "'"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case []any:
		// Array literal
		var items []string
		for _, item := range v {
			itemLiteral := NewLiteral(item)
			items = append(items, itemLiteral.String())
		}
		return "[" + strings.Join(items, ", ") + "]"
	case map[string]any:
		// Map literal
		var pairs []string
		for key, value := range v {
			valueLiteral := NewLiteral(value)
			pairs = append(pairs, key+": "+valueLiteral.String())
		}
		return "{" + strings.Join(pairs, ", ") + "}"
	default:
		// Convert to string as a fallback
		return fmt.Sprintf("%v", v)
	}
}

// Accept applies a visitor to this expression
func (l *LiteralExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(l)
}

// And creates a logical AND with another expression
func (l *LiteralExpression) And(other Expression) Expression {
	// This will be handled by other expressions that embed LiteralExpression
	return nil
}

// Or creates a logical OR with another expression
func (l *LiteralExpression) Or(other Expression) Expression {
	// This will be handled by other expressions that embed LiteralExpression
	return nil
}

// Not creates a logical NOT of this expression
func (l *LiteralExpression) Not() Expression {
	// This will be handled by other expressions that embed LiteralExpression
	return nil
}
