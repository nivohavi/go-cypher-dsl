package builder

import (
	"fmt"
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/util"
)

// withBuilder implements the WithBuilder interface
type withBuilder struct {
	expressions []core.Expression
	whereClause core.Expression
	orderBy     []core.Expression
	orderDir    string
	skipValue   int
	limitValue  int
	prev        core.Buildable
}

// Where adds a WHERE clause
func (w *withBuilder) Where(condition core.Expression) WithBuilder {
	clone := *w
	clone.whereClause = condition
	return &clone
}

// OrderBy adds an ORDER BY clause
func (w *withBuilder) OrderBy(expressions ...core.Expression) WithOrderable {
	clone := *w
	clone.orderBy = expressions
	return &clone
}

// Skip adds a SKIP clause
func (w *withBuilder) Skip(count int) WithBuilder {
	clone := *w
	clone.skipValue = count
	return &clone
}

// Limit adds a LIMIT clause
func (w *withBuilder) Limit(count int) WithBuilder {
	clone := *w
	clone.limitValue = count
	return &clone
}

// Asc specifies ascending order
func (w *withBuilder) Asc() WithBuilder {
	clone := *w
	clone.orderDir = "ASC"
	return &clone
}

// Desc specifies descending order
func (w *withBuilder) Desc() WithBuilder {
	clone := *w
	clone.orderDir = "DESC"
	return &clone
}

// Match adds a MATCH clause
func (w *withBuilder) Match(pattern core.Expression) MatchBuilder {
	return &matchBuilder{
		pattern:  pattern,
		optional: false,
		prev:     w,
	}
}

// OptionalMatch adds an OPTIONAL MATCH clause
func (w *withBuilder) OptionalMatch(pattern core.Expression) MatchBuilder {
	return &matchBuilder{
		pattern:  pattern,
		optional: true,
		prev:     w,
	}
}

// Create adds a CREATE clause
func (w *withBuilder) Create(pattern core.Expression) CreateBuilder {
	return &createBuilder{
		pattern: pattern,
		prev:    w,
	}
}

// Merge adds a MERGE clause
func (w *withBuilder) Merge(pattern core.Expression) MergeBuilder {
	return &mergeBuilder{
		pattern: pattern,
		prev:    w,
	}
}

// Returning adds a RETURN clause
func (w *withBuilder) Returning(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
		prev:        w,
	}
}

// Delete adds a DELETE clause
func (w *withBuilder) Delete(expressions ...core.Expression) DeleteBuilder {
	return &deleteBuilder{
		expressions: expressions,
		detach:      false,
		prev:        w,
	}
}

// DetachDelete adds a DETACH DELETE clause
func (w *withBuilder) DetachDelete(expressions ...core.Expression) DeleteBuilder {
	return &deleteBuilder{
		expressions: expressions,
		detach:      true,
		prev:        w,
	}
}

// Set adds a SET clause
func (w *withBuilder) Set(expression core.Expression) SetBuilder {
	return &setBuilder{
		expressions: []core.Expression{expression},
		prev:        w,
	}
}

// Remove adds a REMOVE clause
func (w *withBuilder) Remove(expression core.Expression) RemoveBuilder {
	return &removeBuilder{
		expressions: []core.Expression{expression},
		prev:        w,
	}
}

// Unwind adds an UNWIND clause
func (w *withBuilder) Unwind(expression core.Expression, alias string) UnwindBuilder {
	return &unwindBuilder{
		expression: expression,
		alias:      alias,
		prev:       w,
	}
}

// With adds another WITH clause
func (w *withBuilder) With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
		prev:        w,
	}
}

// Build builds this WITH into a complete statement
func (w *withBuilder) Build() (core.Statement, error) {
	// If this builder has a previous clause, we need to build that first
	var prevStmt core.Statement
	var err error

	if w.prev != nil {
		prevStmt, err = w.prev.Build()
		if err != nil {
			return nil, err
		}
	}

	// Collect parameters
	paramsMap := make(map[string]any)

	// Extract parameters from expressions
	for _, expr := range w.expressions {
		util.ExtractParameters(expr, paramsMap)
	}

	// Extract parameters from WHERE clause if present
	if w.whereClause != nil {
		util.ExtractParameters(w.whereClause, paramsMap)
	}

	// Extract parameters from ORDER BY expressions if present
	for _, expr := range w.orderBy {
		util.ExtractParameters(expr, paramsMap)
	}

	// Build WITH clause
	parts := []string{"WITH"}

	// Add expressions
	exprStrings := make([]string, len(w.expressions))
	for i, expr := range w.expressions {
		exprStrings[i] = expr.String()
	}
	parts = append(parts, strings.Join(exprStrings, ", "))

	// Add WHERE clause if present
	if w.whereClause != nil {
		parts = append(parts, "WHERE", w.whereClause.String())
	}

	// Add ORDER BY if present
	if len(w.orderBy) > 0 {
		orderExprs := make([]string, len(w.orderBy))
		for i, expr := range w.orderBy {
			orderExprs[i] = expr.String()
		}

		parts = append(parts, "ORDER BY "+strings.Join(orderExprs, ", "))

		if w.orderDir != "" {
			parts = append(parts, w.orderDir)
		}
	}

	// Add SKIP if present
	if w.skipValue > 0 {
		parts = append(parts, fmt.Sprintf("SKIP %d", w.skipValue))
	}

	// Add LIMIT if present
	if w.limitValue > 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", w.limitValue))
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
