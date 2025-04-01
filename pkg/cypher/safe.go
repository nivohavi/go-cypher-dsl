package cypher

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/validation"
)

// Safe returns a parameterized version of a string value to prevent Cypher injection
// This is a convenience function to encourage using parameters instead of string literals
func Safe(value string) core.Expression {
	params := core.NewParameters()
	return params.Add(value)
}

// SafeMap converts a map of string keys and any values to safe parameter expressions
// This is useful for property maps in nodes and relationships
func SafeMap(props map[string]any) map[string]core.Expression {
	if props == nil {
		return nil
	}

	params := core.NewParameters()
	result := make(map[string]core.Expression, len(props))

	for k, v := range props {
		result[k] = params.Add(v)
	}

	return result
}

// SafeNodeProps creates a node with properties safely parameterized
func SafeNodeProps(alias string, properties map[string]any, labels ...string) core.NodeExpression {
	node := Node(labels...).Named(alias)
	return node.WithProperties(SafeMap(properties))
}

// SafeRelProps creates a relationship with properties safely parameterized
func SafeRelProps(from, to core.NodeExpression, relType string, properties map[string]any) core.PatternElement {
	rel := from.RelationshipTo(to, relType)
	return rel.WithProperties(SafeMap(properties))
}

// ValidateStatement validates a statement using the default validator
// Returns true if valid, false and a list of errors otherwise
func ValidateStatement(stmt core.Statement) (bool, []error) {
	if stmt == nil {
		return false, []error{
			fmt.Errorf("statement is nil"),
		}
	}

	// Try to use the statement's own validator if available
	validatable, ok := stmt.(interface {
		Validate() (bool, []error)
	})

	if ok {
		return validatable.Validate()
	}

	// Fall back to creating a new validator
	validator := validation.DefaultValidator()
	return validator.IsValid(stmt.Cypher()), validator.ValidateQuery(stmt.Cypher())
}
