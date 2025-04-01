package core

// RelationshipDirection represents the direction of a relationship
type RelationshipDirection int

const (
	// OUTGOING represents an outgoing relationship (->)
	OUTGOING RelationshipDirection = iota
	// INCOMING represents an incoming relationship (<-)
	INCOMING
	// BIDIRECTIONAL represents a bidirectional relationship (--)
	BIDIRECTIONAL
)

// NodeExpression represents a node in a Cypher query
type NodeExpression interface {
	PatternElement
	PropertyContainer
	// Named sets the alias for this node pattern
	Named(alias string) NodeExpression
	// WithLabels adds labels to this node pattern
	WithLabels(labels ...string) NodeExpression
	// WithProperties adds properties to this node pattern
	WithProperties(properties map[string]Expression) NodeExpression
	// WithProps adds properties with automatic conversion to expressions
	WithProps(properties map[string]interface{}) Expression
	// Props is an alias for WithProps
	Props(properties map[string]interface{}) Expression
	// RelationshipTo creates a relationship from this node to another
	RelationshipTo(other NodeExpression, types ...string) RelationshipPattern
	// RelationshipFrom creates a relationship from another node to this one
	RelationshipFrom(other NodeExpression, types ...string) RelationshipPattern
	// RelationshipBetween creates an undirected relationship between this node and another
	RelationshipBetween(other NodeExpression, types ...string) RelationshipPattern
	// SymbolicName returns the alias of this node pattern
	SymbolicName() string
}

// RelationshipPattern represents a relationship in a Cypher query
type RelationshipPattern interface {
	PatternElement
	PropertyContainer
	// Named sets the alias for this relationship pattern
	Named(alias string) RelationshipPattern
	// WithProperties adds properties to this relationship pattern
	WithProperties(properties map[string]Expression) RelationshipPattern
	// WithProps adds properties with automatic conversion to expressions
	WithProps(properties map[string]interface{}) RelationshipPattern
	// Props is an alias for WithProps
	Props(properties map[string]interface{}) RelationshipPattern
	// StartNode returns the start node of this relationship
	StartNode() NodeExpression
	// EndNode returns the end node of this relationship
	EndNode() NodeExpression
	// Direction returns the direction of this relationship
	Direction() RelationshipDirection
	// Types returns the types of this relationship
	Types() []string
	// SymbolicName returns the alias of this relationship pattern
	SymbolicName() string
}
