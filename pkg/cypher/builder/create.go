package builder

import (
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/util"
)

// createBuilder implements the CreateBuilder interface
type createBuilder struct {
	pattern core.Expression
	prev    core.Buildable
}

// Create adds another CREATE clause
func (c *createBuilder) Create(pattern core.Expression) CreateBuilder {
	return &createBuilder{
		pattern: pattern,
		prev:    c,
	}
}

// With adds a WITH clause
func (c *createBuilder) With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
		prev:        c,
	}
}

// Returning adds a RETURN clause
func (c *createBuilder) Returning(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
		prev:        c,
	}
}

// Set adds a SET clause
func (c *createBuilder) Set(expression core.Expression) SetBuilder {
	return &setBuilder{
		expressions: []core.Expression{expression},
		prev:        c,
	}
}

// Build builds this CREATE into a complete statement
func (c *createBuilder) Build() (core.Statement, error) {
	// If this builder has a previous clause, we need to build that first
	var prevStmt core.Statement
	var err error

	if c.prev != nil {
		prevStmt, err = c.prev.Build()
		if err != nil {
			return nil, err
		}
	}

	// Collect parameters
	paramsMap := make(map[string]any)

	// Extract parameters from pattern
	util.ExtractParameters(c.pattern, paramsMap)

	// Build CREATE clause
	parts := []string{"CREATE", c.pattern.String()}

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
