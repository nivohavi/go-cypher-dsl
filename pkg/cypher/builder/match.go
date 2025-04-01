package builder

import (
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/util"
)

// matchBuilder implements the MatchBuilder interface
type matchBuilder struct {
	pattern     core.Expression
	optional    bool
	whereClause core.Expression
	prev        core.Buildable
}

// Where adds a WHERE clause to this MATCH
func (m *matchBuilder) Where(condition core.Expression) MatchBuilder {
	clone := *m
	clone.whereClause = condition
	return &clone
}

// OptionalMatch adds an OPTIONAL MATCH clause
func (m *matchBuilder) OptionalMatch(pattern core.Expression) MatchBuilder {
	return &matchBuilder{
		pattern:  pattern,
		optional: true,
		prev:     m,
	}
}

// Match adds a MATCH clause
func (m *matchBuilder) Match(pattern core.Expression) MatchBuilder {
	return &matchBuilder{
		pattern:  pattern,
		optional: false,
		prev:     m,
	}
}

// With adds a WITH clause
func (m *matchBuilder) With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
		prev:        m,
	}
}

// Returning adds a RETURN clause
func (m *matchBuilder) Returning(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
		prev:        m,
	}
}

// Create adds a CREATE clause
func (m *matchBuilder) Create(pattern core.Expression) CreateBuilder {
	return &createBuilder{
		pattern: pattern,
		prev:    m,
	}
}

// Merge adds a MERGE clause
func (m *matchBuilder) Merge(pattern core.Expression) MergeBuilder {
	return &mergeBuilder{
		pattern: pattern,
		prev:    m,
	}
}

// Delete adds a DELETE clause
func (m *matchBuilder) Delete(expressions ...core.Expression) DeleteBuilder {
	return &deleteBuilder{
		expressions: expressions,
		detach:      false,
		prev:        m,
	}
}

// DetachDelete adds a DETACH DELETE clause
func (m *matchBuilder) DetachDelete(expressions ...core.Expression) DeleteBuilder {
	return &deleteBuilder{
		expressions: expressions,
		detach:      true,
		prev:        m,
	}
}

// Set adds a SET clause
func (m *matchBuilder) Set(expression core.Expression) SetBuilder {
	return &setBuilder{
		expressions: []core.Expression{expression},
		prev:        m,
	}
}

// Remove adds a REMOVE clause
func (m *matchBuilder) Remove(expression core.Expression) RemoveBuilder {
	return &removeBuilder{
		expressions: []core.Expression{expression},
		prev:        m,
	}
}

// Unwind adds an UNWIND clause
func (m *matchBuilder) Unwind(expression core.Expression, alias string) UnwindBuilder {
	return &unwindBuilder{
		expression: expression,
		alias:      alias,
		prev:       m,
	}
}

// Build builds this MATCH into a complete statement
func (m *matchBuilder) Build() (core.Statement, error) {
	// If this builder has a previous clause, we need to build that first
	var prevStmt core.Statement
	var err error

	if m.prev != nil {
		prevStmt, err = m.prev.Build()
		if err != nil {
			return nil, err
		}
	}

	// Collect parameters
	paramsMap := make(map[string]any)

	// Collect parameters from pattern
	util.ExtractParameters(m.pattern, paramsMap)

	// Collect parameters from where clause
	if m.whereClause != nil {
		util.ExtractParameters(m.whereClause, paramsMap)
	}

	// Build MATCH clause
	var parts []string

	if m.optional {
		parts = append(parts, "OPTIONAL MATCH")
	} else {
		parts = append(parts, "MATCH")
	}

	parts = append(parts, m.pattern.String())

	// Add WHERE clause if present
	if m.whereClause != nil {
		parts = append(parts, "WHERE", m.whereClause.String())
	}

	// Create the query string
	query := strings.Join(parts, " ")

	// Merge with previous statement if any
	if prevStmt != nil {
		// Merge previous parameters
		prevParams := prevStmt.Params()
		if prevParams != nil {
			for k, v := range prevParams {
				paramsMap[k] = v
			}
		}

		return core.NewStatement(prevStmt.Cypher()+" "+query, paramsMap), nil
	}

	// Create a new statement
	return core.NewStatement(query, paramsMap), nil
}

// Helper function to extract parameters from expressions recursively (deprecated, use util.ExtractParameters instead)
func extractParameters(expr core.Expression, paramsMap map[string]any) {
	util.ExtractParameters(expr, paramsMap)
}
