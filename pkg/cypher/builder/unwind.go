package builder

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// unwindBuilder implements the UnwindBuilder interface
type unwindBuilder struct {
	expression core.Expression
	alias      string
	prev       core.Buildable
}

// Where adds a WHERE clause
func (u *unwindBuilder) Where(condition core.Expression) WhereBuilder {
	return &whereBuilder{
		condition: condition,
		prev:      u,
	}
}

// Match adds a MATCH clause
func (u *unwindBuilder) Match(pattern core.Expression) MatchBuilder {
	return &matchBuilder{
		pattern:  pattern,
		optional: false,
		prev:     u,
	}
}

// OptionalMatch adds an OPTIONAL MATCH clause
func (u *unwindBuilder) OptionalMatch(pattern core.Expression) MatchBuilder {
	return &matchBuilder{
		pattern:  pattern,
		optional: true,
		prev:     u,
	}
}

// With adds a WITH clause
func (u *unwindBuilder) With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
		prev:        u,
	}
}

// Returning adds a RETURN clause
func (u *unwindBuilder) Returning(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
		prev:        u,
	}
}

// Build builds this UNWIND into a complete statement
func (u *unwindBuilder) Build() (core.Statement, error) {
	// Simple implementation for now
	var cypher string

	// If this builder has a previous clause, we need to build that first
	if u.prev != nil {
		prevStmt, err := u.prev.Build()
		if err != nil {
			return nil, err
		}
		cypher = prevStmt.Cypher() + " "
	}

	// Add UNWIND keyword, expression and alias
	cypher += fmt.Sprintf("UNWIND %s AS %s", u.expression.String(), u.alias)

	return core.NewStatement(cypher, nil), nil
}
