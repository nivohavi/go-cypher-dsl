package builder

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// mergeBuilder implements the MergeBuilder interface
type mergeBuilder struct {
	pattern       core.Expression
	onCreateExprs []core.Expression
	onMatchExprs  []core.Expression
	prev          core.Buildable
}

// OnCreate adds an ON CREATE SET clause
func (m *mergeBuilder) OnCreate(expression core.Expression) MergeBuilder {
	clone := *m
	clone.onCreateExprs = append(clone.onCreateExprs, expression)
	return &clone
}

// OnMatch adds an ON MATCH SET clause
func (m *mergeBuilder) OnMatch(expression core.Expression) MergeBuilder {
	clone := *m
	clone.onMatchExprs = append(clone.onMatchExprs, expression)
	return &clone
}

// With adds a WITH clause
func (m *mergeBuilder) With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
		prev:        m,
	}
}

// Returning adds a RETURN clause
func (m *mergeBuilder) Returning(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
		prev:        m,
	}
}

// Build builds this MERGE into a complete statement
func (m *mergeBuilder) Build() (core.Statement, error) {
	// Simple implementation for now
	var cypher string

	// If this builder has a previous clause, we need to build that first
	if m.prev != nil {
		prevStmt, err := m.prev.Build()
		if err != nil {
			return nil, err
		}
		cypher = prevStmt.Cypher() + " "
	}

	// Add MERGE keyword and pattern
	cypher += "MERGE " + m.pattern.String()

	// Add ON CREATE SET clause if present
	if len(m.onCreateExprs) > 0 {
		cypher += " ON CREATE SET "
		for i, expr := range m.onCreateExprs {
			if i > 0 {
				cypher += ", "
			}
			cypher += expr.String()
		}
	}

	// Add ON MATCH SET clause if present
	if len(m.onMatchExprs) > 0 {
		cypher += " ON MATCH SET "
		for i, expr := range m.onMatchExprs {
			if i > 0 {
				cypher += ", "
			}
			cypher += expr.String()
		}
	}

	return core.NewStatement(cypher, nil), nil
}
