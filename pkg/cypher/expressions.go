package cypher

// Expression is the base interface for all Cypher expressions in the AST.
// Any element that can be rendered to Cypher syntax should implement this.
type Expression interface {
	// Accept applies the ExpressionVisitor to this expression
	Accept(visitor ExpressionVisitor) any
	// String returns a string representation of the expression (for debugging)
	String() string
}

// ExpressionVisitor is the visitor interface for traversing the expression AST
type ExpressionVisitor interface {
	// Visit processes an expression and returns a result
	Visit(element any) any
}

// ExpressionContainer is an interface for objects that contain expressions
type ExpressionContainer interface {
	// Expressions returns all contained expressions
	Expressions() []Expression
}

// NamedExpression is an expression that can be referenced by name
type NamedExpression interface {
	Expression
	// SymbolicName returns the symbolic name of this expression
	SymbolicName() string
}

// PropertyContainer is an interface for expressions that can have properties
type PropertyContainer interface {
	Expression
	// Property returns a property access expression for this container
	Property(propertyName string) PropertyExpression

	// Helper methods for common property access patterns
	// These are convenience methods that map to the underlying Property() method
	Prop(name string) PropertyExpression                    // Alias for Property
	Props(properties map[string]interface{}) Expression     // Set multiple properties at once
	WithProps(properties map[string]interface{}) Expression // Alias for Props
}

// NodeExpression represents a node in the graph with access to common operations
type NodeExpression interface {
	PropertyContainer
	NamedExpression

	// Node relationship methods
	RelationshipTo(other NodeExpression, types ...string) RelationshipPattern
	RelationshipFrom(other NodeExpression, types ...string) RelationshipPattern
	RelationshipBetween(other NodeExpression, types ...string) RelationshipPattern

	// Identity methods
	Named(name string) NodeExpression
	WithLabels(labels ...string) NodeExpression
	WithProperties(properties map[string]Expression) NodeExpression
}

// Direction represents the direction of a relationship
type Direction int

const (
	// OUTGOING represents an outgoing relationship (->)
	OUTGOING Direction = iota
	// INCOMING represents an incoming relationship (<-)
	INCOMING
	// BIDIRECTIONAL represents a bidirectional relationship (--)
	BIDIRECTIONAL
)

// RelationshipPattern represents a relationship between nodes
type RelationshipPattern interface {
	PropertyContainer
	NamedExpression

	// Identity methods
	Named(name string) RelationshipPattern
	WithType(relType string) RelationshipPattern
	WithTypes(types ...string) RelationshipPattern
	WithProperties(properties map[string]Expression) RelationshipPattern

	// Relationship direction
	Direction() Direction
	SetDirection(direction Direction) RelationshipPattern

	// Access to nodes
	StartNode() NodeExpression
	EndNode() NodeExpression

	// Pattern continuation
	RelationshipTo(node NodeExpression, types ...string) PatternElement
	RelationshipFrom(node NodeExpression, types ...string) PatternElement
}

// PropertyExpression represents a property access operation
type PropertyExpression interface {
	Expression
	// PropertyName returns the name of the property being accessed
	PropertyName() string
	// Subject returns the expression on which the property is accessed
	Subject() Expression

	// Comparison operations
	Eq(other interface{}) ComparisonExpression
	Ne(other interface{}) ComparisonExpression
	Gt(other interface{}) ComparisonExpression
	Gte(other interface{}) ComparisonExpression
	Lt(other interface{}) ComparisonExpression
	Lte(other interface{}) ComparisonExpression

	// String operations
	Contains(value interface{}) BooleanExpression
	StartsWith(value interface{}) BooleanExpression
	EndsWith(value interface{}) BooleanExpression

	// Collection operations
	In(list interface{}) BooleanExpression

	// Null checks
	IsNull() BooleanExpression
	IsNotNull() BooleanExpression

	// Aliases for improved readability
	Equals(other interface{}) ComparisonExpression
	NotEquals(other interface{}) ComparisonExpression
	GreaterThan(other interface{}) ComparisonExpression
	GreaterThanOrEqual(other interface{}) ComparisonExpression
	LessThan(other interface{}) ComparisonExpression
	LessThanOrEqual(other interface{}) ComparisonExpression
}

// BooleanExpression represents an expression that evaluates to a boolean
type BooleanExpression interface {
	Expression
	// Logical operations
	And(other BooleanExpression) BooleanExpression
	Or(other BooleanExpression) BooleanExpression
	Not() BooleanExpression

	// Marker method
	IsBooleanExpression()
}

// ComparisonExpression represents a comparison between two expressions
type ComparisonExpression interface {
	BooleanExpression
	// Left returns the left side of the comparison
	Left() Expression
	// Right returns the right side of the comparison
	Right() Expression
	// Operator returns the comparison operator
	Operator() string
}

// LiteralExpression represents a literal value in Cypher
type LiteralExpression interface {
	Expression
	// Value returns the Go value of this literal
	Value() interface{}
}

// ParameterExpression represents a parameter in a Cypher query
type ParameterExpression interface {
	Expression
	// Name returns the parameter name
	Name() string
	// Value returns the parameter value, if set
	Value() interface{}
	// HasValue checks if the parameter has a value assigned
	HasValue() bool
}

// ListExpression represents a list of expressions
type ListExpression interface {
	Expression
	// Elements returns the expressions in this list
	Elements() []Expression
}

// MapExpression represents a map of key-value pairs
type MapExpression interface {
	Expression
	// Entries returns the entries in this map
	Entries() map[string]Expression
}

// FunctionExpression represents a function call in Cypher
type FunctionExpression interface {
	Expression
	// FunctionName returns the name of the function
	FunctionName() string
	// Arguments returns the arguments passed to the function
	Arguments() []Expression
}

// AliasedExpression represents an expression with an alias
type AliasedExpression interface {
	Expression
	// Inner returns the underlying expression
	Inner() Expression
	// Alias returns the alias for this expression
	Alias() string
}

// Statement represents a Cypher statement (e.g., a query)
type Statement interface {
	// Cypher returns the Cypher query string
	Cypher() string
	// Params returns the parameters for this statement
	Params() map[string]interface{}
	// Accept accepts a visitor for this statement
	Accept(visitor StatementVisitor) any
}

// StatementVisitor is the visitor interface for visiting statements
type StatementVisitor interface {
	// Visit visits an element and returns any value
	Visit(element any) any
}

// PatternElement is the base interface for all pattern elements (nodes, relationships, paths)
type PatternElement interface {
	Expression
}

// Buildable is an interface for objects that can build a Statement
type Buildable interface {
	// Build constructs the final Statement
	Build() (Statement, error)
}
