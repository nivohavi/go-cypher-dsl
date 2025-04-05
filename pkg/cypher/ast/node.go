// Package ast provides abstract syntax tree implementations for Cypher DSL
package ast

import (
	"fmt"
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

// nodePattern represents a node pattern in Cypher (e.g., (n:Person))
type nodePattern struct {
	labels     []string
	alias      string
	properties map[string]core.Expression
}

// Node creates a new node pattern with the given labels
func Node(labels ...string) core.NodeExpression {
	return &nodePattern{
		labels:     labels,
		properties: make(map[string]core.Expression),
	}
}

// Named sets the alias for this node pattern
func (n *nodePattern) Named(alias string) core.NodeExpression {
	clone := *n
	clone.alias = alias
	return &clone
}

// WithLabels adds labels to this node pattern
func (n *nodePattern) WithLabels(labels ...string) core.NodeExpression {
	clone := *n
	clone.labels = append(clone.labels, labels...)
	return &clone
}

// WithProperties adds properties to this node pattern
func (n *nodePattern) WithProperties(properties map[string]core.Expression) core.NodeExpression {
	clone := *n
	for k, v := range properties {
		clone.properties[k] = v
	}
	return &clone
}

// WithProps adds properties with automatic conversion to expressions
func (n *nodePattern) WithProps(properties map[string]interface{}) core.NodeExpression {
	clone := *n
	for k, v := range properties {
		switch val := v.(type) {
		case core.Expression:
			clone.properties[k] = val
		default:
			clone.properties[k] = expr.LiteralFromValue(val)
		}
	}
	return &clone
}

// Props is an alias for WithProps
func (n *nodePattern) Props(properties map[string]interface{}) core.Expression {
	return n.WithProps(properties)
}

// Property returns a property access expression for this node
func (n *nodePattern) Property(propertyName string) core.PropertyExpression {
	return &propertyExpression{
		subject:      n,
		propertyName: propertyName,
	}
}

// RelationshipTo creates a relationship from this node to another
func (n *nodePattern) RelationshipTo(other core.NodeExpression, types ...string) core.RelationshipPattern {
	if otherNode, ok := other.(*nodePattern); ok {
		rel := &relationshipPattern{
			startNode:  n,
			endNode:    otherNode,
			types:      types,
			direction:  core.OUTGOING,
			properties: make(map[string]core.Expression),
		}
		return rel
	}
	// Handle case when other is not a *nodePattern
	panic("other must be created with Node() function")
}

// RelationshipFrom creates a relationship from another node to this one
func (n *nodePattern) RelationshipFrom(other core.NodeExpression, types ...string) core.RelationshipPattern {
	if otherNode, ok := other.(*nodePattern); ok {
		rel := &relationshipPattern{
			startNode:  otherNode,
			endNode:    n,
			types:      types,
			direction:  core.INCOMING,
			properties: make(map[string]core.Expression),
		}
		return rel
	}
	// Handle case when other is not a *nodePattern
	panic("other must be created with Node() function")
}

// RelationshipBetween creates an undirected relationship between this node and another
func (n *nodePattern) RelationshipBetween(other core.NodeExpression, types ...string) core.RelationshipPattern {
	if otherNode, ok := other.(*nodePattern); ok {
		rel := &relationshipPattern{
			startNode:  n,
			endNode:    otherNode,
			types:      types,
			direction:  core.BIDIRECTIONAL,
			properties: make(map[string]core.Expression),
		}
		return rel
	}
	// Handle case when other is not a *nodePattern
	panic("other must be created with Node() function")
}

// SymbolicName returns the alias of this node pattern
func (n *nodePattern) SymbolicName() string {
	return n.alias
}

// Accept implements the Expression interface
func (n *nodePattern) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(n)
}

// String returns a string representation of this node pattern
func (n *nodePattern) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	if n.alias != "" {
		sb.WriteString(n.alias)
	}

	// Write labels
	for _, label := range n.labels {
		sb.WriteString(":")
		sb.WriteString(label)
	}

	// Write properties if present
	if len(n.properties) > 0 {
		sb.WriteString(" {")
		first := true
		for k, v := range n.properties {
			if !first {
				sb.WriteString(", ")
			}
			first = false
			sb.WriteString(k)
			sb.WriteString(": ")
			sb.WriteString(v.String())
		}
		sb.WriteString("}")
	}

	sb.WriteString(")")
	return sb.String()
}

// And creates a logical AND with another expression
func (n *nodePattern) And(other core.Expression) core.Expression {
	return expr.And(n, other)
}

// Or creates a logical OR with another expression
func (n *nodePattern) Or(other core.Expression) core.Expression {
	return expr.Or(n, other)
}

// Not creates a logical NOT of this expression
func (n *nodePattern) Not() core.Expression {
	return expr.Not(n)
}

// Expressions returns all expressions contained in this node pattern
func (n *nodePattern) Expressions() []core.Expression {
	result := make([]core.Expression, 0, len(n.properties))
	for _, prop := range n.properties {
		result = append(result, prop)
	}
	return result
}

// propertyExpression represents a property access expression (e.g., n.name)
type propertyExpression struct {
	subject      core.Expression
	propertyName string
}

// Accept implements the Expression interface
func (p *propertyExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(p)
}

// String returns a string representation of this property expression
func (p *propertyExpression) String() string {
	if namedExpr, ok := p.subject.(core.NamedExpression); ok && namedExpr.SymbolicName() != "" {
		return fmt.Sprintf("%s.%s", namedExpr.SymbolicName(), p.propertyName)
	}
	return fmt.Sprintf("%s.%s", p.subject.String(), p.propertyName)
}

// Eq creates an equals comparison with the given value
func (p *propertyExpression) Eq(value any) core.Expression {
	return expr.Equals(p, expr.LiteralFromValue(value))
}

// Gt creates a greater-than comparison with the given value
func (p *propertyExpression) Gt(value any) core.Expression {
	return expr.GreaterThan(p, expr.LiteralFromValue(value))
}

// Lt creates a less-than comparison with the given value
func (p *propertyExpression) Lt(value any) core.Expression {
	return expr.LessThan(p, expr.LiteralFromValue(value))
}

// Gte creates a greater-than-or-equal comparison with the given value
func (p *propertyExpression) Gte(value any) core.Expression {
	return expr.GreaterThanEqual(p, expr.LiteralFromValue(value))
}

// Lte creates a less-than-or-equal comparison with the given value
func (p *propertyExpression) Lte(value any) core.Expression {
	return expr.LessThanEqual(p, expr.LiteralFromValue(value))
}

// IsNull creates a null check
func (p *propertyExpression) IsNull() core.Expression {
	return expr.IsNull(p)
}

// IsNotNull creates a not-null check
func (p *propertyExpression) IsNotNull() core.Expression {
	return expr.Not(expr.IsNull(p))
}

// In creates an IN comparison with the given values
func (p *propertyExpression) In(values ...any) core.Expression {
	return expr.In(p, values...)
}

// StartsWith creates a STARTS WITH comparison
func (p *propertyExpression) StartsWith(value string) core.Expression {
	return expr.StartsWith(p, value)
}

// EndsWith creates an ENDS WITH comparison
func (p *propertyExpression) EndsWith(value string) core.Expression {
	return expr.EndsWith(p, value)
}

// Contains creates a CONTAINS comparison
func (p *propertyExpression) Contains(value string) core.Expression {
	return expr.Contains(p, value)
}

// RegularExpression creates a =~ comparison with a regular expression
func (p *propertyExpression) RegularExpression(pattern string) core.Expression {
	return expr.RegularExpression(p, pattern)
}

// And creates a logical AND with another expression
func (p *propertyExpression) And(other core.Expression) core.Expression {
	return expr.And(p, other)
}

// Or creates a logical OR with another expression
func (p *propertyExpression) Or(other core.Expression) core.Expression {
	return expr.Or(p, other)
}

// Not creates a logical NOT of this expression
func (p *propertyExpression) Not() core.Expression {
	return expr.Not(p)
}
