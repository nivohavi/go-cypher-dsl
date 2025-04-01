package builder

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// Match creates a new MATCH clause
func Match(pattern core.Expression) MatchBuilder {
	return &matchBuilder{
		pattern:  pattern,
		optional: false,
	}
}

// OptionalMatch creates a new OPTIONAL MATCH clause
func OptionalMatch(pattern core.Expression) MatchBuilder {
	return &matchBuilder{
		pattern:  pattern,
		optional: true,
	}
}

// Create creates a new CREATE clause
func Create(pattern core.Expression) CreateBuilder {
	return &createBuilder{
		pattern: pattern,
	}
}

// Merge creates a new MERGE clause
func Merge(pattern core.Expression) MergeBuilder {
	return &mergeBuilder{
		pattern: pattern,
	}
}

// Return creates a new RETURN clause
func Return(expressions ...core.Expression) ReturnBuilder {
	return &returnBuilder{
		expressions: expressions,
	}
}

// With creates a new WITH clause
func With(expressions ...core.Expression) WithBuilder {
	return &withBuilder{
		expressions: expressions,
	}
}

// Where creates a new WHERE clause
func Where(condition core.Expression) WhereBuilder {
	return &whereBuilder{
		condition: condition,
	}
}

// OrderBy creates a new ORDER BY clause
func OrderBy(expressions ...core.Expression) OrderByBuilder {
	return &orderByBuilder{
		expressions: expressions,
	}
}

// Limit creates a new LIMIT clause
func Limit(limit int) LimitBuilder {
	return &limitBuilder{
		limit: limit,
	}
}

// Skip creates a new SKIP clause
func Skip(skip int) SkipBuilder {
	return &skipBuilder{
		skip: skip,
	}
}

// Delete creates a new DELETE clause
func Delete(expressions ...core.Expression) DeleteBuilder {
	return &deleteBuilder{
		expressions: expressions,
		detach:      false,
	}
}

// DetachDelete creates a new DETACH DELETE clause
func DetachDelete(expressions ...core.Expression) DeleteBuilder {
	return &deleteBuilder{
		expressions: expressions,
		detach:      true,
	}
}

// Set creates a new SET clause
func Set(expression core.Expression) SetBuilder {
	return &setBuilder{
		expressions: []core.Expression{expression},
	}
}

// Remove creates a new REMOVE clause
func Remove(expression core.Expression) RemoveBuilder {
	return &removeBuilder{
		expressions: []core.Expression{expression},
	}
}

// Unwind creates a new UNWIND clause
func Unwind(expression core.Expression, alias string) UnwindBuilder {
	return &unwindBuilder{
		expression: expression,
		alias:      alias,
	}
}
