package builder

import (
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/util"
)

// setBuilder implements the SetBuilder interface
type setBuilder struct {
	expressions []core.Expression
	prev        core.Buildable
}

// And adds another SET operation
func (s *setBuilder) And(expression core.Expression) SetBuilder {
	clone := *s
	clone.expressions = append(clone.expressions, expression)
	return &clone
}

// With adds a WITH clause
func (s *setBuilder) With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
		prev:        s,
	}
}

// Return adds a RETURN clause
func (s *setBuilder) Returning(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
		prev:        s,
	}
}

// Build builds this SET into a complete statement
func (s *setBuilder) Build() (core.Statement, error) {
	// If this builder has a previous clause, we need to build that first
	var prevStmt core.Statement
	var err error

	if s.prev != nil {
		prevStmt, err = s.prev.Build()
		if err != nil {
			return nil, err
		}
	}

	// Collect parameters
	paramsMap := make(map[string]any)

	// Build SET clause
	parts := []string{"SET"}

	exprStrings := make([]string, len(s.expressions))
	for i, expr := range s.expressions {
		exprStrings[i] = expr.String()

		// Collect parameters from the expression if it's a parameter expression
		util.ExtractParameters(expr, paramsMap)
	}

	parts = append(parts, strings.Join(exprStrings, ", "))

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
