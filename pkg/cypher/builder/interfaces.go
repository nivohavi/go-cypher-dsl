// Package builder provides builder implementations for constructing Cypher queries
package builder

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// MatchBuilder builds MATCH clauses
type MatchBuilder interface {
	core.Buildable
	// Where adds a WHERE clause
	Where(condition core.Expression) MatchBuilder
	// OptionalMatch adds an OPTIONAL MATCH clause
	OptionalMatch(pattern core.Expression) MatchBuilder
	// Match adds a MATCH clause
	Match(pattern core.Expression) MatchBuilder
	// Create adds a CREATE clause
	Create(pattern core.Expression) CreateBuilder
	// Merge adds a MERGE clause
	Merge(pattern core.Expression) MergeBuilder
	// With adds a WITH clause
	With(expressions ...core.Expression) WithBuilder
	// Return adds a RETURN clause
	Returning(expressions ...core.Expression) ReturnBuilder
	// Delete adds a DELETE clause
	Delete(expressions ...core.Expression) DeleteBuilder
	// DetachDelete adds a DETACH DELETE clause
	DetachDelete(expressions ...core.Expression) DeleteBuilder
	// Set adds a SET clause
	Set(expression core.Expression) SetBuilder
	// Remove adds a REMOVE clause
	Remove(expression core.Expression) RemoveBuilder
	// Unwind adds an UNWIND clause
	Unwind(expression core.Expression, alias string) UnwindBuilder
}

// WhereBuilder builds WHERE clauses
type WhereBuilder interface {
	core.Buildable
	// Where adds another condition with AND
	AndWhere(condition core.Expression) WhereBuilder
	// OrWhere adds another condition with OR
	OrWhere(condition core.Expression) WhereBuilder
	// With adds a WITH clause
	With(expressions ...core.Expression) WithBuilder
	// Return adds a RETURN clause
	Returning(expressions ...core.Expression) ReturnBuilder
}

// WithBuilder builds WITH clauses
type WithBuilder interface {
	core.Buildable
	// Where adds a WHERE clause
	Where(condition core.Expression) WithBuilder
	// OrderBy adds an ORDER BY clause
	OrderBy(expressions ...core.Expression) WithOrderable
	// Skip adds a SKIP clause
	Skip(count int) WithBuilder
	// Limit adds a LIMIT clause
	Limit(count int) WithBuilder
	// Match adds a MATCH clause
	Match(pattern core.Expression) MatchBuilder
	// OptionalMatch adds an OPTIONAL MATCH clause
	OptionalMatch(pattern core.Expression) MatchBuilder
	// Create adds a CREATE clause
	Create(pattern core.Expression) CreateBuilder
	// Merge adds a MERGE clause
	Merge(pattern core.Expression) MergeBuilder
	// With adds another WITH clause
	With(expressions ...core.Expression) WithBuilder
	// Return adds a RETURN clause
	Returning(expressions ...core.Expression) ReturnBuilder
	// Delete adds a DELETE clause
	Delete(expressions ...core.Expression) DeleteBuilder
	// DetachDelete adds a DETACH DELETE clause
	DetachDelete(expressions ...core.Expression) DeleteBuilder
	// Set adds a SET clause
	Set(expression core.Expression) SetBuilder
	// Remove adds a REMOVE clause
	Remove(expression core.Expression) RemoveBuilder
	// Unwind adds an UNWIND clause
	Unwind(expression core.Expression, alias string) UnwindBuilder
}

// ReturnBuilder builds RETURN clauses
type ReturnBuilder interface {
	core.Buildable
	// OrderBy adds an ORDER BY clause
	OrderBy(expressions ...core.Expression) ReturnOrderable
	// Skip adds a SKIP clause
	Skip(count int) ReturnBuilder
	// Limit adds a LIMIT clause
	Limit(count int) ReturnBuilder
}

// ReturnOrderable is a ReturnBuilder that supports ORDER BY
type ReturnOrderable interface {
	ReturnBuilder
	// Asc specifies ascending order
	Asc() ReturnBuilder
	// Desc specifies descending order
	Desc() ReturnBuilder
}

// WithOrderable is a WithBuilder that supports ORDER BY
type WithOrderable interface {
	WithBuilder
	// Asc specifies ascending order
	Asc() WithBuilder
	// Desc specifies descending order
	Desc() WithBuilder
}

// CreateBuilder builds CREATE clauses
type CreateBuilder interface {
	core.Buildable
	// Create adds another CREATE clause
	Create(pattern core.Expression) CreateBuilder
	// With adds a WITH clause
	With(expressions ...core.Expression) WithBuilder
	// Return adds a RETURN clause
	Returning(expressions ...core.Expression) ReturnBuilder
	// Set adds a SET clause
	Set(expression core.Expression) SetBuilder
}

// MergeBuilder builds MERGE clauses
type MergeBuilder interface {
	core.Buildable
	// OnCreate adds an ON CREATE SET clause
	OnCreate(expression core.Expression) MergeBuilder
	// OnMatch adds an ON MATCH SET clause
	OnMatch(expression core.Expression) MergeBuilder
	// With adds a WITH clause
	With(expressions ...core.Expression) WithBuilder
	// Return adds a RETURN clause
	Returning(expressions ...core.Expression) ReturnBuilder
}

// DeleteBuilder builds DELETE clauses
type DeleteBuilder interface {
	core.Buildable
	// With adds a WITH clause
	With(expressions ...core.Expression) WithBuilder
	// Return adds a RETURN clause
	Returning(expressions ...core.Expression) ReturnBuilder
}

// SetBuilder builds SET clauses
type SetBuilder interface {
	core.Buildable
	// And adds another SET operation
	And(expression core.Expression) SetBuilder
	// With adds a WITH clause
	With(expressions ...core.Expression) WithBuilder
	// Return adds a RETURN clause
	Returning(expressions ...core.Expression) ReturnBuilder
}

// RemoveBuilder builds REMOVE clauses
type RemoveBuilder interface {
	core.Buildable
	// And adds another REMOVE operation
	And(expression core.Expression) RemoveBuilder
	// With adds a WITH clause
	With(expressions ...core.Expression) WithBuilder
	// Return adds a RETURN clause
	Returning(expressions ...core.Expression) ReturnBuilder
}

// UnwindBuilder builds UNWIND clauses
type UnwindBuilder interface {
	core.Buildable
	// Where adds a WHERE clause
	Where(condition core.Expression) WhereBuilder
	// Match adds a MATCH clause
	Match(pattern core.Expression) MatchBuilder
	// OptionalMatch adds an OPTIONAL MATCH clause
	OptionalMatch(pattern core.Expression) MatchBuilder
	// With adds a WITH clause
	With(expressions ...core.Expression) WithBuilder
	// Return adds a RETURN clause
	Returning(expressions ...core.Expression) ReturnBuilder
}

// OrderByBuilder builds ORDER BY clauses
type OrderByBuilder interface {
	core.Buildable
	// Asc specifies ascending order
	Asc() OrderByBuilder
	// Desc specifies descending order
	Desc() OrderByBuilder
	// Skip adds a SKIP clause
	Skip(count int) OrderByBuilder
	// Limit adds a LIMIT clause
	Limit(count int) OrderByBuilder
}

// LimitBuilder builds LIMIT clauses
type LimitBuilder interface {
	core.Buildable
}

// SkipBuilder builds SKIP clauses
type SkipBuilder interface {
	core.Buildable
	// Limit adds a LIMIT clause
	Limit(count int) LimitBuilder
}
