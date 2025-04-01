package builder

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

// whereBuilder implements the WhereBuilder interface
type whereBuilder struct {
	condition core.Expression
	prev      core.Buildable
}

// AndWhere adds another condition with AND
func (w *whereBuilder) AndWhere(condition core.Expression) WhereBuilder {
	return &whereBuilder{
		condition: expr.And(w.condition, condition),
		prev:      w.prev,
	}
}

// OrWhere adds another condition with OR
func (w *whereBuilder) OrWhere(condition core.Expression) WhereBuilder {
	return &whereBuilder{
		condition: expr.Or(w.condition, condition),
		prev:      w.prev,
	}
}

// With adds a WITH clause
func (w *whereBuilder) With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
		prev:        w,
	}
}

// Returning adds a RETURN clause
func (w *whereBuilder) Returning(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
		prev:        w,
	}
}

// Build builds this WHERE into a complete statement
func (w *whereBuilder) Build() (core.Statement, error) {
	// Simple implementation for now
	var cypher string

	// If this builder has a previous clause, we need to build that first
	if w.prev != nil {
		prevStmt, err := w.prev.Build()
		if err != nil {
			return nil, err
		}
		cypher = prevStmt.Cypher() + " "
	}

	// Add WHERE keyword and condition
	cypher += "WHERE " + w.condition.String()

	return core.NewStatement(cypher, nil), nil
}
