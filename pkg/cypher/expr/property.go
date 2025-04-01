package expr

import (
	"fmt"
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// PropertyExpression represents a property access expression (e.g., n.name)
type PropertyExpression struct {
	Subject      core.Expression
	PropertyName string
	Chain        []string // Additional property accesses for deep property paths
}

// Accept implements the Expression interface
func (p *PropertyExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(p)
}

// String returns a string representation of this property expression
func (p *PropertyExpression) String() string {
	var subjectStr string
	if named, ok := p.Subject.(core.NamedExpression); ok && named.SymbolicName() != "" {
		subjectStr = named.SymbolicName()
	} else {
		subjectStr = p.Subject.String()
	}

	if len(p.Chain) == 0 {
		return fmt.Sprintf("%s.%s", subjectStr, p.PropertyName)
	}

	// Handle deep property path
	allProps := append([]string{p.PropertyName}, p.Chain...)
	return fmt.Sprintf("%s.%s", subjectStr, strings.Join(allProps, "."))
}

// Eq creates an equals comparison with the given value
func (p *PropertyExpression) Eq(value any) core.Expression {
	return Equals(p, LiteralFromValue(value))
}

// Gt creates a greater-than comparison with the given value
func (p *PropertyExpression) Gt(value any) core.Expression {
	return GreaterThan(p, LiteralFromValue(value))
}

// Lt creates a less-than comparison with the given value
func (p *PropertyExpression) Lt(value any) core.Expression {
	return LessThan(p, LiteralFromValue(value))
}

// Gte creates a greater-than-or-equal comparison with the given value
func (p *PropertyExpression) Gte(value any) core.Expression {
	return GreaterThanEqual(p, LiteralFromValue(value))
}

// Lte creates a less-than-or-equal comparison with the given value
func (p *PropertyExpression) Lte(value any) core.Expression {
	return LessThanEqual(p, LiteralFromValue(value))
}

// IsNull creates a null check
func (p *PropertyExpression) IsNull() core.Expression {
	return IsNull(p)
}

// IsNotNull creates a not-null check
func (p *PropertyExpression) IsNotNull() core.Expression {
	return IsNotNull(p)
}

// In creates an IN comparison with the given values
func (p *PropertyExpression) In(values ...any) core.Expression {
	return In(p, values...)
}

// StartsWith creates a STARTS WITH comparison
func (p *PropertyExpression) StartsWith(value string) core.Expression {
	return StartsWith(p, value)
}

// EndsWith creates an ENDS WITH comparison
func (p *PropertyExpression) EndsWith(value string) core.Expression {
	return EndsWith(p, value)
}

// Contains creates a CONTAINS comparison
func (p *PropertyExpression) Contains(value string) core.Expression {
	return Contains(p, value)
}

// RegularExpression creates a =~ comparison with a regular expression
func (p *PropertyExpression) RegularExpression(pattern string) core.Expression {
	return RegularExpression(p, pattern)
}

// Property creates a property access expression
func Property(entity string, property string, additionalProperties ...string) core.Expression {
	return &PropertyExpression{
		Subject:      String(entity),
		PropertyName: property,
		Chain:        additionalProperties,
	}
}

// And creates a logical AND with another expression
func (p *PropertyExpression) And(other core.Expression) core.Expression {
	return And(p, other)
}

// Or creates a logical OR with another expression
func (p *PropertyExpression) Or(other core.Expression) core.Expression {
	return Or(p, other)
}

// Not creates a logical NOT of this expression
func (p *PropertyExpression) Not() core.Expression {
	return Not(p)
}

// NewProperty creates a new property expression
func NewProperty(subject core.Expression, propertyName string) *PropertyExpression {
	return &PropertyExpression{
		Subject:      subject,
		PropertyName: propertyName,
		Chain:        []string{},
	}
}
