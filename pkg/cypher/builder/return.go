package builder

import (
	"fmt"
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// returnBuilder implements the ReturnBuilder interface
type returnBuilder struct {
	expressions    []core.Expression
	orderBy        []core.Expression
	orderDir       string
	skipValue      int
	limitValue     int
	distinct       bool
	returnAll      bool
	returnAllProps bool
	prev           core.Buildable
}

// OrderBy adds an ORDER BY clause
func (r *returnBuilder) OrderBy(expressions ...core.Expression) ReturnOrderable {
	clone := *r
	clone.orderBy = expressions
	return &clone
}

// Skip adds a SKIP clause
func (r *returnBuilder) Skip(count int) ReturnBuilder {
	clone := *r
	clone.skipValue = count
	return &clone
}

// Limit adds a LIMIT clause
func (r *returnBuilder) Limit(count int) ReturnBuilder {
	clone := *r
	clone.limitValue = count
	return &clone
}

// Asc specifies ascending order
func (r *returnBuilder) Asc() ReturnBuilder {
	clone := *r
	clone.orderDir = "ASC"
	return &clone
}

// Desc specifies descending order
func (r *returnBuilder) Desc() ReturnBuilder {
	clone := *r
	clone.orderDir = "DESC"
	return &clone
}

// Build builds this RETURN into a complete statement
func (r *returnBuilder) Build() (core.Statement, error) {
	// If this builder has a previous clause, we need to build that first
	var prevStmt core.Statement
	var err error

	if r.prev != nil {
		prevStmt, err = r.prev.Build()
		if err != nil {
			return nil, err
		}
	}

	// Collect all parameters from expressions
	paramsMap := make(map[string]any)

	// Extract parameters from expressions
	for _, expr := range r.expressions {
		extractParameters(expr, paramsMap)
	}

	// Extract parameters from ORDER BY expressions if present
	for _, expr := range r.orderBy {
		extractParameters(expr, paramsMap)
	}

	// Build RETURN clause
	var parts []string

	if r.distinct {
		parts = append(parts, "RETURN DISTINCT")
	} else if r.returnAll {
		parts = append(parts, "RETURN *")
	} else {
		parts = append(parts, "RETURN")
	}

	if !r.returnAll {
		exprs := make([]string, len(r.expressions))
		for i, expr := range r.expressions {
			exprs[i] = expr.String()
		}
		parts = append(parts, strings.Join(exprs, ", "))
	}

	// Add ORDER BY if present
	if len(r.orderBy) > 0 {
		orderExprs := make([]string, len(r.orderBy))
		for i, expr := range r.orderBy {
			orderExprs[i] = expr.String()
		}

		parts = append(parts, "ORDER BY "+strings.Join(orderExprs, ", "))

		if r.orderDir != "" {
			parts = append(parts, r.orderDir)
		}
	}

	// Add SKIP if present
	if r.skipValue > 0 {
		parts = append(parts, fmt.Sprintf("SKIP %d", r.skipValue))
	}

	// Add LIMIT if present
	if r.limitValue > 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", r.limitValue))
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
