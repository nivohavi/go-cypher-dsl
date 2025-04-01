package builder

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// removeBuilder implements the RemoveBuilder interface
type removeBuilder struct {
	expressions []core.Expression
	prev        core.Buildable
}

// And adds another REMOVE operation
func (r *removeBuilder) And(expression core.Expression) RemoveBuilder {
	clone := *r
	clone.expressions = append(clone.expressions, expression)
	return &clone
}

// With adds a WITH clause
func (r *removeBuilder) With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
		prev:        r,
	}
}

// Returning adds a RETURN clause
func (r *removeBuilder) Returning(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
		prev:        r,
	}
}

// Build builds this REMOVE into a complete statement
func (r *removeBuilder) Build() (core.Statement, error) {
	// Simple implementation for now
	var cypher string

	// If this builder has a previous clause, we need to build that first
	if r.prev != nil {
		prevStmt, err := r.prev.Build()
		if err != nil {
			return nil, err
		}
		cypher = prevStmt.Cypher() + " "
	}

	// Add REMOVE keyword
	cypher += "REMOVE "

	// Add expressions
	for i, expr := range r.expressions {
		if i > 0 {
			cypher += ", "
		}
		cypher += expr.String()
	}

	return core.NewStatement(cypher, nil), nil
}
