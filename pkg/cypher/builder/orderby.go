package builder

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// orderByBuilder implements the OrderByBuilder interface
type orderByBuilder struct {
	expressions []core.Expression
	direction   string // "ASC" or "DESC"
	skipValue   int
	limitValue  int
	prev        core.Buildable
}

// Asc specifies ascending order
func (o *orderByBuilder) Asc() OrderByBuilder {
	clone := *o
	clone.direction = "ASC"
	return &clone
}

// Desc specifies descending order
func (o *orderByBuilder) Desc() OrderByBuilder {
	clone := *o
	clone.direction = "DESC"
	return &clone
}

// Skip adds a SKIP clause
func (o *orderByBuilder) Skip(count int) OrderByBuilder {
	clone := *o
	clone.skipValue = count
	return &clone
}

// Limit adds a LIMIT clause
func (o *orderByBuilder) Limit(count int) OrderByBuilder {
	clone := *o
	clone.limitValue = count
	return &clone
}

// Build builds this ORDER BY into a complete statement
func (o *orderByBuilder) Build() (core.Statement, error) {
	// Simple implementation for now
	var cypher string

	// If this builder has a previous clause, we need to build that first
	if o.prev != nil {
		prevStmt, err := o.prev.Build()
		if err != nil {
			return nil, err
		}
		cypher = prevStmt.Cypher() + " "
	}

	// Add ORDER BY clause
	cypher += "ORDER BY "

	// Add expressions with direction
	for i, expr := range o.expressions {
		if i > 0 {
			cypher += ", "
		}
		cypher += expr.String()
		if o.direction != "" {
			cypher += " " + o.direction
		}
	}

	// Add SKIP clause if needed
	if o.skipValue > 0 {
		cypher += fmt.Sprintf(" SKIP %d", o.skipValue)
	}

	// Add LIMIT clause if needed
	if o.limitValue > 0 {
		cypher += fmt.Sprintf(" LIMIT %d", o.limitValue)
	}

	return core.NewStatement(cypher, nil), nil
}
