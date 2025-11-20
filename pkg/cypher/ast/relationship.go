package ast

import (
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

// relationshipPattern represents a relationship in a Cypher pattern
type relationshipPattern struct {
	direction  core.RelationshipDirection
	startNode  core.NodeExpression
	endNode    core.NodeExpression
	types      []string
	alias      string
	properties map[string]core.Expression
}

// Named sets the alias for this relationship pattern
func (r *relationshipPattern) Named(alias string) core.RelationshipPattern {
	clone := *r
	clone.alias = alias
	return &clone
}

// WithProperties adds properties to this relationship pattern
func (r *relationshipPattern) WithProperties(properties map[string]core.Expression) core.RelationshipPattern {
	clone := *r
	if clone.properties == nil {
		clone.properties = make(map[string]core.Expression)
	}
	for k, v := range properties {
		clone.properties[k] = v
	}
	return &clone
}

// WithProps adds properties with automatic conversion to expressions
func (r *relationshipPattern) WithProps(properties map[string]interface{}) core.RelationshipPattern {
	clone := *r
	if clone.properties == nil {
		clone.properties = make(map[string]core.Expression)
	}
	for k, v := range properties {
		clone.properties[k] = expr.LiteralFromValue(v)
	}
	return &clone
}

// Props is an alias for WithProps
func (r *relationshipPattern) Props(properties map[string]interface{}) core.RelationshipPattern {
	return r.WithProps(properties)
}

// Property returns a property access expression for this relationship
func (r *relationshipPattern) Property(propertyName string) core.PropertyExpression {
	return &propertyExpression{
		subject:      r,
		propertyName: propertyName,
	}
}

// Prop is an alias for Property
func (r *relationshipPattern) Prop(propertyName string) core.PropertyExpression {
	return r.Property(propertyName)
}

// StartNode returns the start node of this relationship
func (r *relationshipPattern) StartNode() core.NodeExpression {
	return r.startNode
}

// EndNode returns the end node of this relationship
func (r *relationshipPattern) EndNode() core.NodeExpression {
	return r.endNode
}

// Direction returns the direction of this relationship
func (r *relationshipPattern) Direction() core.RelationshipDirection {
	return r.direction
}

// Types returns the types of this relationship
func (r *relationshipPattern) Types() []string {
	return r.types
}

// Accept implements the Expression interface
func (r *relationshipPattern) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(r)
}

// SymbolicName returns the alias of this relationship pattern
func (r *relationshipPattern) SymbolicName() string {
	return r.alias
}

// And creates a logical AND with another expression
func (r *relationshipPattern) And(other core.Expression) core.Expression {
	return expr.And(r, other)
}

// Or creates a logical OR with another expression
func (r *relationshipPattern) Or(other core.Expression) core.Expression {
	return expr.Or(r, other)
}

// Not creates a logical NOT of this expression
func (r *relationshipPattern) Not() core.Expression {
	return expr.Not(r)
}

// String returns a string representation of this relationship pattern
func (r *relationshipPattern) String() string {
	var sb strings.Builder

	// Start with the relationship arrow
	switch r.direction {
	case core.OUTGOING:
		sb.WriteString("-[")
	case core.INCOMING:
		sb.WriteString("<-[")
	case core.BIDIRECTIONAL:
		sb.WriteString("-[")
	}

	if r.alias != "" {
		sb.WriteString(r.alias)
	}

	for _, typ := range r.types {
		sb.WriteString(":`")
		sb.WriteString(typ)
		sb.WriteString("`")
	}

	// Write properties if present
	if len(r.properties) > 0 {
		sb.WriteString(" {")
		first := true
		for k, v := range r.properties {
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

	sb.WriteString("]")

	// End with the appropriate arrow
	switch r.direction {
	case core.OUTGOING:
		sb.WriteString("->")
	case core.INCOMING:
		sb.WriteString("-")
	case core.BIDIRECTIONAL:
		sb.WriteString("-")
	}

	return sb.String()
}

// Expressions returns all expressions contained in this relationship pattern
func (r *relationshipPattern) Expressions() []core.Expression {
	result := make([]core.Expression, 0, len(r.properties)+2)

	// Add start and end nodes if they exist
	if r.startNode != nil {
		result = append(result, r.startNode)
	}
	if r.endNode != nil {
		result = append(result, r.endNode)
	}

	// Add properties
	for _, prop := range r.properties {
		result = append(result, prop)
	}

	return result
}
