package builder

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// skipBuilder implements the SkipBuilder interface
type skipBuilder struct {
	skip int
	prev core.Buildable
}

// Limit adds a LIMIT clause
func (s *skipBuilder) Limit(count int) LimitBuilder {
	return &limitBuilder{
		limit: count,
		prev:  s,
	}
}

// Build builds this SKIP into a complete statement
func (s *skipBuilder) Build() (core.Statement, error) {
	// Simple implementation for now
	var cypher string

	// If this builder has a previous clause, we need to build that first
	if s.prev != nil {
		prevStmt, err := s.prev.Build()
		if err != nil {
			return nil, err
		}
		cypher = prevStmt.Cypher() + " "
	}

	// Add SKIP clause
	cypher += fmt.Sprintf("SKIP %d", s.skip)

	return core.NewStatement(cypher, nil), nil
}
