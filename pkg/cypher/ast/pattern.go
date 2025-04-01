package ast

import (
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

// PatternExpression represents a complete pattern in Cypher (e.g., (a)-[r]->(b))
type PatternExpression struct {
	elements []core.PatternElement
	alias    string
}

// Accept implements the Expression interface
func (p *PatternExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(p)
}

// String returns a string representation of the pattern
func (p *PatternExpression) String() string {
	var builder strings.Builder

	if p.alias != "" {
		builder.WriteString(p.alias)
		builder.WriteString(" = ")
	}

	// For simple path patterns (node-relationship-node), we need special formatting
	if len(p.elements) == 3 {
		node1, ok1 := p.elements[0].(*nodePattern)
		rel, ok2 := p.elements[1].(*relationshipPattern)
		node2, ok3 := p.elements[2].(*nodePattern)

		if ok1 && ok2 && ok3 {
			// Format as a path pattern: (n1)-[r]->(n2)
			builder.WriteString(node1.String())
			builder.WriteString(rel.String())
			builder.WriteString(node2.String())
			return builder.String()
		}
	}

	// Default handling for other patterns
	for _, element := range p.elements {
		builder.WriteString(element.String())
	}

	return builder.String()
}

// SymbolicName returns the alias of this pattern
func (p *PatternExpression) SymbolicName() string {
	return p.alias
}

// Expressions returns all the expressions contained in this pattern
func (p *PatternExpression) Expressions() []core.Expression {
	result := make([]core.Expression, len(p.elements))
	for i, elem := range p.elements {
		result[i] = elem
	}
	return result
}

// And creates a logical AND with another expression
func (p *PatternExpression) And(other core.Expression) core.Expression {
	return expr.And(p, other)
}

// Or creates a logical OR with another expression
func (p *PatternExpression) Or(other core.Expression) core.Expression {
	return expr.Or(p, other)
}

// Not creates a logical NOT of this expression
func (p *PatternExpression) Not() core.Expression {
	return expr.Not(p)
}

// Pattern creates a new pattern from the given elements
func Pattern(elements ...core.PatternElement) core.Expression {
	return &PatternExpression{
		elements: elements,
	}
}

// Path is an alias for Pattern
func Path(elements ...core.PatternElement) core.Expression {
	return Pattern(elements...)
}

// RelationshipChain represents a chain of relationships
type RelationshipChain struct {
	startNode     core.NodeExpression
	relationships []core.RelationshipPattern
	endNodes      []core.NodeExpression
}

// Accept implements the Expression interface
func (r *RelationshipChain) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(r)
}

// String returns a string representation of this relationship chain
func (r *RelationshipChain) String() string {
	var sb strings.Builder

	// Add start node
	sb.WriteString(r.startNode.String())

	// Add relationships and end nodes
	for i, rel := range r.relationships {
		sb.WriteString(rel.String())
		sb.WriteString(r.endNodes[i].String())
	}

	return sb.String()
}

// And creates a logical AND with another expression
func (r *RelationshipChain) And(other core.Expression) core.Expression {
	return expr.And(r, other)
}

// Or creates a logical OR with another expression
func (r *RelationshipChain) Or(other core.Expression) core.Expression {
	return expr.Or(r, other)
}

// Not creates a logical NOT of this expression
func (r *RelationshipChain) Not() core.Expression {
	return expr.Not(r)
}

// Chain creates a new relationship chain
func Chain(startNode core.NodeExpression, relationships ...core.RelationshipPattern) core.Expression {
	if len(relationships) == 0 {
		return startNode
	}

	endNodes := make([]core.NodeExpression, len(relationships))
	for i, rel := range relationships {
		endNodes[i] = rel.EndNode()
	}

	return &RelationshipChain{
		startNode:     startNode,
		relationships: relationships,
		endNodes:      endNodes,
	}
}

// RelateNodes creates a simple relationship between two nodes
func RelateNodes(fromNode core.NodeExpression, relType string, toNode core.NodeExpression) core.Expression {
	rel := fromNode.RelationshipTo(toNode, relType)
	return Pattern(fromNode, rel, toNode)
}

// RelateBidirectionally creates a bidirectional relationship between two nodes
func RelateBidirectionally(node1 core.NodeExpression, relType string, node2 core.NodeExpression) core.Expression {
	rel := node1.RelationshipBetween(node2, relType)
	return Pattern(node1, rel, node2)
}
