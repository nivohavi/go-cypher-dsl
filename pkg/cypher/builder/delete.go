package builder

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// deleteBuilder implements the DeleteBuilder interface
type deleteBuilder struct {
	expressions []core.Expression
	detach      bool
	prev        core.Buildable
}

// With adds a WITH clause
func (d *deleteBuilder) With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
		prev:        d,
	}
}

// Returning adds a RETURN clause
func (d *deleteBuilder) Returning(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
		prev:        d,
	}
}

// Build builds this DELETE into a complete statement
func (d *deleteBuilder) Build() (core.Statement, error) {
	// Simple implementation for now
	var cypher string

	// If this builder has a previous clause, we need to build that first
	if d.prev != nil {
		prevStmt, err := d.prev.Build()
		if err != nil {
			return nil, err
		}
		cypher = prevStmt.Cypher() + " "
	}

	// Add DELETE or DETACH DELETE keyword
	if d.detach {
		cypher += "DETACH DELETE "
	} else {
		cypher += "DELETE "
	}

	// Add expressions
	for i, expr := range d.expressions {
		if i > 0 {
			cypher += ", "
		}
		cypher += expr.String()
	}

	return core.NewStatement(cypher, nil), nil
}
