package cypher

import (
	"fmt"
)

// NodeType represents a pre-defined node type in the schema
type NodeType struct {
	Label string
	alias string
}

// Of creates a node pattern of this type
func (n *NodeType) Of() PatternElement {
	return Node(n.Label)
}

// Named creates a node pattern of this type with the given alias
func (n *NodeType) Named(alias string) PatternElement {
	// Store the alias with the node type, so when schemaPropertyExpression
	// is created, it can access this alias
	n.alias = alias
	return Node(n.Label).Named(alias)
}

// WithProperties creates a node pattern of this type with the given properties
func (n *NodeType) WithProperties(properties map[string]Expression) PatternElement {
	return Node(n.Label).(*nodePattern).WithProperties(properties)
}

// Property returns a property accessor for this node type
func (n *NodeType) Property(propertyName string) PropertyExpression {
	// For schema properties, we want to use the node's label/alias to reference
	// properties rather than creating a node pattern with an underscore
	return &schemaPropertyExpression{
		nodeType:     n,
		propertyName: propertyName,
		nodeAlias:    n.alias, // Pass the alias
	}
}

// RelationshipType represents a pre-defined relationship type in the schema
type RelationshipType struct {
	Type  string
	alias string
}

// Of creates a relationship pattern of this type
func (r *RelationshipType) Of() *relationshipBuilder {
	return Relationship(r.Type)
}

// Named creates a relationship pattern of this type with the given alias
func (r *RelationshipType) Named(alias string) *relationshipBuilder {
	r.alias = alias
	return Relationship(r.Type).Named(alias)
}

// Direction sets the direction of this relationship pattern
func (r *RelationshipType) Direction(dir Direction) *relationshipBuilder {
	return Relationship(r.Type).Direction(dir)
}

// WithProperties adds properties to this relationship pattern
func (r *RelationshipType) WithProperties(properties map[string]Expression) *relationshipBuilder {
	return Relationship(r.Type).WithProperties(properties)
}

// Property returns a property accessor for this relationship type
func (r *RelationshipType) Property(propertyName string) PropertyExpression {
	// For schema properties, we want to use the relationship's type/alias to reference
	// properties rather than creating a relationship pattern with an underscore
	return &schemaPropertyExpression{
		relType:      r,
		propertyName: propertyName,
		relAlias:     r.alias, // Pass the alias
	}
}

// NewNodeType creates a new node type in the schema
func NewNodeType(label string) *NodeType {
	return &NodeType{
		Label: label,
	}
}

// NewRelationshipType creates a new relationship type in the schema
func NewRelationshipType(relType string) *RelationshipType {
	return &RelationshipType{
		Type: relType,
	}
}

// Schema example:

// PersonNode is an example of a meta-model for a Person node
type PersonNode struct {
	*NodeType
}

// Person creates a new PersonNode
func Person() *PersonNode {
	return &PersonNode{
		NodeType: NewNodeType("Person"),
	}
}

// Name returns the name property of a Person node
func (p *PersonNode) Name() PropertyExpression {
	return p.Property("name")
}

// Age returns the age property of a Person node
func (p *PersonNode) Age() PropertyExpression {
	return p.Property("age")
}

// MovieNode is an example of a meta-model for a Movie node
type MovieNode struct {
	*NodeType
}

// Movie creates a new MovieNode
func Movie() *MovieNode {
	return &MovieNode{
		NodeType: NewNodeType("Movie"),
	}
}

// Title returns the title property of a Movie node
func (m *MovieNode) Title() PropertyExpression {
	return m.Property("title")
}

// Released returns the released property of a Movie node
func (m *MovieNode) Released() PropertyExpression {
	return m.Property("released")
}

// ActedInRelationship is an example of a meta-model for an ACTED_IN relationship
type ActedInRelationship struct {
	*RelationshipType
}

// ActedIn creates a new ActedInRelationship
func ActedIn() *ActedInRelationship {
	return &ActedInRelationship{
		RelationshipType: NewRelationshipType("ACTED_IN"),
	}
}

// Role returns the role property of an ACTED_IN relationship
func (a *ActedInRelationship) Role() PropertyExpression {
	return a.Property("role")
}

// schemaPropertyExpression represents a property access in the context of a schema type
type schemaPropertyExpression struct {
	nodeType     *NodeType
	relType      *RelationshipType
	propertyName string
	// Track the aliases for use in rendering
	nodeAlias string
	relAlias  string
}

// PropertyName returns the name of the property being accessed
func (p *schemaPropertyExpression) PropertyName() string {
	return p.propertyName
}

// Subject returns the expression on which the property is accessed
func (p *schemaPropertyExpression) Subject() Expression {
	if p.nodeType != nil {
		return p.nodeType.Of()
	}
	if p.relType != nil {
		return p.relType.Of().Build()
	}
	return nil
}

// String returns a string representation of this property expression
func (p *schemaPropertyExpression) String() string {
	if p.nodeType != nil {
		return fmt.Sprintf("%s.%s", p.nodeType.Label, p.propertyName)
	}
	if p.relType != nil {
		return fmt.Sprintf("%s.%s", p.relType.Type, p.propertyName)
	}
	return fmt.Sprintf("unknown.%s", p.propertyName)
}

// Accept implements the Expression interface
func (p *schemaPropertyExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(p)
}

// Eq creates an equality comparison with another expression
func (p *schemaPropertyExpression) Eq(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    autoExpression(other),
		operator: "=",
	}
}

// Ne creates an inequality comparison with another expression
func (p *schemaPropertyExpression) Ne(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    autoExpression(other),
		operator: "<>",
	}
}

// Gt creates a greater-than comparison with another expression
func (p *schemaPropertyExpression) Gt(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    autoExpression(other),
		operator: ">",
	}
}

// Gte creates a greater-than-or-equal comparison with another expression
func (p *schemaPropertyExpression) Gte(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    autoExpression(other),
		operator: ">=",
	}
}

// Lt creates a less-than comparison with another expression
func (p *schemaPropertyExpression) Lt(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    autoExpression(other),
		operator: "<",
	}
}

// Lte creates a less-than-or-equal comparison with another expression
func (p *schemaPropertyExpression) Lte(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    autoExpression(other),
		operator: "<=",
	}
}

// Contains creates a CONTAINS operation with a string
func (p *schemaPropertyExpression) Contains(value interface{}) BooleanExpression {
	return &stringOperation{
		left:     p,
		right:    autoExpression(value),
		operator: "CONTAINS",
	}
}

// StartsWith creates a STARTS WITH operation with a string
func (p *schemaPropertyExpression) StartsWith(value interface{}) BooleanExpression {
	return &stringOperation{
		left:     p,
		right:    autoExpression(value),
		operator: "STARTS WITH",
	}
}

// EndsWith creates a ENDS WITH operation with a string
func (p *schemaPropertyExpression) EndsWith(value interface{}) BooleanExpression {
	return &stringOperation{
		left:     p,
		right:    autoExpression(value),
		operator: "ENDS WITH",
	}
}

// In creates an IN operation with a list
func (p *schemaPropertyExpression) In(list interface{}) BooleanExpression {
	return &comparisonExpression{
		left:     p,
		right:    autoExpression(list),
		operator: "IN",
	}
}

// IsNull creates an IS NULL check
func (p *schemaPropertyExpression) IsNull() BooleanExpression {
	return &nullCheck{
		expression: p,
		isNull:     true,
	}
}

// IsNotNull creates an IS NOT NULL check
func (p *schemaPropertyExpression) IsNotNull() BooleanExpression {
	return &nullCheck{
		expression: p,
		isNull:     false,
	}
}

// Aliases for improved readability
func (p *schemaPropertyExpression) Equals(other interface{}) ComparisonExpression {
	return p.Eq(other)
}

func (p *schemaPropertyExpression) NotEquals(other interface{}) ComparisonExpression {
	return p.Ne(other)
}

func (p *schemaPropertyExpression) GreaterThan(other interface{}) ComparisonExpression {
	return p.Gt(other)
}

func (p *schemaPropertyExpression) GreaterThanOrEqual(other interface{}) ComparisonExpression {
	return p.Gte(other)
}

func (p *schemaPropertyExpression) LessThan(other interface{}) ComparisonExpression {
	return p.Lt(other)
}

func (p *schemaPropertyExpression) LessThanOrEqual(other interface{}) ComparisonExpression {
	return p.Lte(other)
}
