package core

import (
	"fmt"
	"strings"
)

// ParameterExpression represents a parameter in a Cypher query
type ParameterExpression struct {
	name  string
	value any
}

// NewParameter creates a new parameter expression
func NewParameter(name string, value any) *ParameterExpression {
	return &ParameterExpression{
		name:  name,
		value: value,
	}
}

// Name returns the parameter name
func (p *ParameterExpression) Name() string {
	return p.name
}

// Value returns the parameter value
func (p *ParameterExpression) Value() any {
	return p.value
}

// String returns the string representation of this parameter
func (p *ParameterExpression) String() string {
	return "$" + p.name
}

// Accept implements the Expression interface
func (p *ParameterExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(p)
}

// And creates a logical AND with another expression
func (p *ParameterExpression) And(other Expression) Expression {
	// We need to import expr, so we'll use a type assertion to create an AND
	// This will be handled by other expressions that embed ParameterExpression
	return nil
}

// Or creates a logical OR with another expression
func (p *ParameterExpression) Or(other Expression) Expression {
	// We need to import expr, so we'll use a type assertion to create an OR
	// This will be handled by other expressions that embed ParameterExpression
	return nil
}

// Not creates a logical NOT of this expression
func (p *ParameterExpression) Not() Expression {
	// This will be handled by other expressions that embed ParameterExpression
	return nil
}

// Parameters maintains a map of parameters for a Cypher query
type Parameters struct {
	params       map[string]any
	paramCounter int
	paramPrefix  string
}

// NewParameters creates a new parameter container
func NewParameters() *Parameters {
	return &Parameters{
		params:       make(map[string]any),
		paramCounter: 0,
		paramPrefix:  "param",
	}
}

// Add adds a value as a parameter and returns the parameter expression
func (p *Parameters) Add(value any) *ParameterExpression {
	p.paramCounter++
	name := fmt.Sprintf("%s%d", p.paramPrefix, p.paramCounter)
	p.params[name] = value
	return NewParameter(name, value)
}

// AddNamed adds a value as a named parameter
func (p *Parameters) AddNamed(name string, value any) *ParameterExpression {
	// Clean the name to ensure it's valid
	name = cleanParameterName(name)
	p.params[name] = value
	return NewParameter(name, value)
}

// Get returns the map of all parameters
func (p *Parameters) Get() map[string]any {
	return p.params
}

// Merge merges another parameter container into this one
func (p *Parameters) Merge(other *Parameters) {
	if other == nil {
		return
	}
	for k, v := range other.params {
		p.params[k] = v
	}
}

// cleanParameterName ensures the parameter name is valid for Cypher
func cleanParameterName(name string) string {
	// Replace invalid characters with underscores
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")

	// Ensure it starts with a letter
	if len(name) > 0 && !isLetter(rune(name[0])) {
		name = "p_" + name
	}

	return name
}

// isLetter checks if a rune is a letter
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
