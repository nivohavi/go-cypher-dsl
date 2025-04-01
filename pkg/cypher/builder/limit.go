package builder

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// limitBuilder implements the LimitBuilder interface
type limitBuilder struct {
	limit int
	prev  core.Buildable
}

// Build builds this LIMIT into a complete statement
func (l *limitBuilder) Build() (core.Statement, error) {
	// Simple implementation for now
	var cypher string

	// If this builder has a previous clause, we need to build that first
	if l.prev != nil {
		prevStmt, err := l.prev.Build()
		if err != nil {
			return nil, err
		}
		cypher = prevStmt.Cypher() + " "
	}

	// Add LIMIT clause
	cypher += fmt.Sprintf("LIMIT %d", l.limit)

	return core.NewStatement(cypher, nil), nil
}
