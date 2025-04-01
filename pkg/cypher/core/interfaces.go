// Package core provides the foundational interfaces for the Cypher DSL
package core

// Expression is the base interface for all expressions in Cypher
type Expression interface {
	// Accept implements the visitor pattern for traversing expressions
	Accept(visitor ExpressionVisitor) any
	// String returns a string representation of this expression
	String() string
	// And creates a logical AND with another expression
	And(other Expression) Expression
	// Or creates a logical OR with another expression
	Or(other Expression) Expression
	// Not creates a logical NOT of this expression
	Not() Expression
}

// ExpressionVisitor is used to visit expressions in the visitor pattern
type ExpressionVisitor interface {
	// Visit processes the provided element and returns a result
	Visit(element any) any
}

// ExpressionContainer is an interface for objects that contain expressions
type ExpressionContainer interface {
	// Expressions returns all contained expressions
	Expressions() []Expression
}

// PropertyContainer represents an element that can have properties
type PropertyContainer interface {
	Expression
	// Property returns a property access expression
	Property(propertyName string) PropertyExpression
}

// PropertyExpression represents a property access expression
type PropertyExpression interface {
	Expression
	// Eq creates an equals comparison with the given value
	Eq(value any) Expression
	// Gt creates a greater-than comparison with the given value
	Gt(value any) Expression
	// Lt creates a less-than comparison with the given value
	Lt(value any) Expression
	// Gte creates a greater-than-or-equal comparison with the given value
	Gte(value any) Expression
	// Lte creates a less-than-or-equal comparison with the given value
	Lte(value any) Expression
	// IsNull creates a null check
	IsNull() Expression
	// IsNotNull creates a not-null check
	IsNotNull() Expression
	// In creates an IN comparison with the given values
	In(values ...any) Expression
	// StartsWith creates a STARTS WITH comparison
	StartsWith(value string) Expression
	// EndsWith creates an ENDS WITH comparison
	EndsWith(value string) Expression
	// Contains creates a CONTAINS comparison
	Contains(value string) Expression
	// RegularExpression creates a =~ comparison with a regular expression
	RegularExpression(pattern string) Expression
}

// PatternElement represents an element in a Cypher pattern
type PatternElement interface {
	Expression
}

// NamedExpression represents an expression with a symbolic name (alias)
type NamedExpression interface {
	Expression
	SymbolicName() string
}

// BooleanExpression represents an expression that can be used in boolean contexts
type BooleanExpression interface {
	Expression
	// And creates a logical AND with another expression
	And(other Expression) Expression
	// Or creates a logical OR with another expression
	Or(other Expression) Expression
	// Xor creates a logical XOR with another expression
	Xor(other Expression) Expression
}

// Statement represents a complete Cypher statement
type Statement interface {
	// Cypher returns the Cypher query string
	Cypher() string
	// Params returns the parameters for this statement
	Params() map[string]any
	// Accept applies a visitor to this statement
	Accept(visitor StatementVisitor) any
}

// StatementVisitor is used to visit statements in the visitor pattern
type StatementVisitor interface {
	// Visit processes the provided element and returns a result
	Visit(element any) any
}

// Buildable represents an object that can be built into a complete statement
type Buildable interface {
	// Build builds this object into a complete statement
	Build() (Statement, error)
}
