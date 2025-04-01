// Package cypher provides a fluent API for building Cypher queries
package cypher

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/builder"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/renderer"
)

// Property creates a property expression for a node or relationship
func Property(entity, property string) core.Expression {
	return expr.Property(entity, property)
}

// Parameters creates a new parameter container
func Parameters() *core.Parameters {
	return core.NewParameters()
}

// Param creates a new parameter with the given value
func Param(value any) core.Expression {
	return core.NewParameter("", value)
}

// NamedParam creates a new named parameter
func NamedParam(name string, value any) core.Expression {
	return core.NewParameter(name, value)
}

// ParamWithValue creates a new named parameter (alias for NamedParam)
func ParamWithValue(name string, value any) core.Expression {
	return NamedParam(name, value)
}

// Literal creates a new literal expression
func Literal(value any) core.Expression {
	return core.NewLiteral(value)
}

// Node creates a new node pattern
func Node(labels ...string) core.NodeExpression {
	return ast.Node(labels...)
}

// NodeWithProperties creates a new node pattern with properties
func NodeWithProperties(properties map[string]any, labels ...string) core.NodeExpression {
	node := ast.Node(labels...)
	return node.WithProperties(convertProperties(properties))
}

// NamedNode creates a new named node pattern
func NamedNode(alias string, labels ...string) core.NodeExpression {
	return ast.Node(labels...).Named(alias)
}

// NamedNodeWithProperties creates a new named node pattern with properties
func NamedNodeWithProperties(alias string, properties map[string]any, labels ...string) core.NodeExpression {
	node := ast.Node(labels...).Named(alias)
	return node.WithProperties(convertProperties(properties))
}

// Pattern creates a pattern expression
func Pattern(elements ...core.PatternElement) core.Expression {
	return ast.Pattern(elements...)
}

// Path creates a path pattern
func Path(elements ...core.PatternElement) core.Expression {
	return ast.Path(elements...)
}

// Chain creates a relationship chain
func Chain(startNode core.NodeExpression, relationships ...core.RelationshipPattern) core.Expression {
	return ast.Chain(startNode, relationships...)
}

// RelateNodes creates a relationship between two nodes
func RelateNodes(fromNode core.NodeExpression, toNode core.NodeExpression, relType string) core.PatternElement {
	rel := fromNode.RelationshipTo(toNode, relType)
	return rel
}

// RelateBidirectionally creates a bidirectional relationship between two nodes
func RelateBidirectionally(fromNode core.NodeExpression, toNode core.NodeExpression, relType string) core.PatternElement {
	rel := fromNode.RelationshipBetween(toNode, relType)
	return rel
}

// Match creates a MATCH clause
func Match(pattern core.Expression) builder.MatchBuilder {
	return builder.Match(pattern)
}

// OptionalMatch creates an OPTIONAL MATCH clause
func OptionalMatch(pattern core.Expression) builder.MatchBuilder {
	return builder.OptionalMatch(pattern)
}

// Create creates a CREATE clause
func Create(pattern core.Expression) builder.CreateBuilder {
	return builder.Create(pattern)
}

// Merge creates a MERGE clause
func Merge(pattern core.Expression) builder.MergeBuilder {
	return builder.Merge(pattern)
}

// Return creates a RETURN clause
func Return(expressions ...core.Expression) builder.ReturnBuilder {
	return builder.Return(expressions...)
}

// With creates a WITH clause
func With(expressions ...core.Expression) builder.WithBuilder {
	return builder.With(expressions...)
}

// Where creates a WHERE clause
func Where(condition core.Expression) builder.WhereBuilder {
	return builder.Where(condition)
}

// OrderBy creates an ORDER BY clause
func OrderBy(expressions ...core.Expression) builder.OrderByBuilder {
	return builder.OrderBy(expressions...)
}

// Skip creates a SKIP clause
func Skip(count int) builder.SkipBuilder {
	return builder.Skip(count)
}

// Limit creates a LIMIT clause
func Limit(count int) builder.LimitBuilder {
	return builder.Limit(count)
}

// Delete creates a DELETE clause
func Delete(expressions ...core.Expression) builder.DeleteBuilder {
	return builder.Delete(expressions...)
}

// DetachDelete creates a DETACH DELETE clause
func DetachDelete(expressions ...core.Expression) builder.DeleteBuilder {
	return builder.DetachDelete(expressions...)
}

// Set creates a SET clause
func Set(expression core.Expression) builder.SetBuilder {
	return builder.Set(expression)
}

// Remove creates a REMOVE clause
func Remove(expression core.Expression) builder.RemoveBuilder {
	return builder.Remove(expression)
}

// Unwind creates an UNWIND clause
func Unwind(expression core.Expression, alias string) builder.UnwindBuilder {
	return builder.Unwind(expression, alias)
}

// Eq creates an equality expression
func Eq(left, right core.Expression) core.Expression {
	return expr.Equals(left, right)
}

// Ne creates a not-equals expression
func Ne(left, right core.Expression) core.Expression {
	return expr.NotEquals(left, right)
}

// Gt creates a greater-than expression
func Gt(left, right core.Expression) core.Expression {
	return expr.GreaterThan(left, right)
}

// Gte creates a greater-than-or-equal expression
func Gte(left, right core.Expression) core.Expression {
	return expr.GreaterThanEqual(left, right)
}

// Lt creates a less-than expression
func Lt(left, right core.Expression) core.Expression {
	return expr.LessThan(left, right)
}

// Lte creates a less-than-or-equal expression
func Lte(left, right core.Expression) core.Expression {
	return expr.LessThanEqual(left, right)
}

// And creates a logical AND expression
func And(left, right core.Expression) core.Expression {
	return expr.And(left, right)
}

// Or creates a logical OR expression
func Or(left, right core.Expression) core.Expression {
	return expr.Or(left, right)
}

// Xor creates a logical XOR expression
func Xor(left, right core.Expression) core.Expression {
	return expr.Xor(left, right)
}

// Not creates a logical NOT expression
func Not(expression core.Expression) core.Expression {
	return expr.Not(expression)
}

// Render renders a statement using the default renderer
func Render(statement core.Statement) string {
	r := renderer.NewCypherRenderer()
	return r.Render(statement)
}

// RenderWithParams renders a statement and returns the Cypher and parameters
func RenderWithParams(statement core.Statement) (string, map[string]any) {
	r := renderer.NewCypherRenderer()
	return r.RenderWithParams(statement)
}

// PrettyPrint renders a statement with pretty printing
func PrettyPrint(statement core.Statement) string {
	r := renderer.NewCypherRenderer().WithPrettyPrint(true)
	return r.Render(statement)
}

// Literal utility functions

// String creates a string literal
func String(value string) core.Expression {
	return expr.String(value)
}

// Integer creates an integer literal
func Integer(value int64) core.Expression {
	return expr.Integer(value)
}

// Float creates a float literal
func Float(value float64) core.Expression {
	return expr.Float(value)
}

// Boolean creates a boolean literal
func Boolean(value bool) core.Expression {
	return expr.Boolean(value)
}

// Null creates a null literal
func Null() core.Expression {
	return expr.Null()
}

// List creates a list expression
func List(elements ...core.Expression) core.Expression {
	return expr.List(elements...)
}

// Map creates a map expression
func Map(entries map[string]core.Expression) core.Expression {
	return expr.Map(entries)
}

// Equals creates an equality comparison expression
func Equals(left, right core.Expression) core.Expression {
	return expr.Equals(left, right)
}

// NotEquals creates a not-equals comparison expression
func NotEquals(left, right core.Expression) core.Expression {
	return expr.NotEquals(left, right)
}

// GreaterThan creates a greater-than comparison expression
func GreaterThan(left, right core.Expression) core.Expression {
	return expr.GreaterThan(left, right)
}

// LessThan creates a less-than comparison expression
func LessThan(left, right core.Expression) core.Expression {
	return expr.LessThan(left, right)
}

// GreaterThanEqual creates a greater-than-or-equal comparison expression
func GreaterThanEqual(left, right core.Expression) core.Expression {
	return expr.GreaterThanEqual(left, right)
}

// LessThanEqual creates a less-than-or-equal comparison expression
func LessThanEqual(left, right core.Expression) core.Expression {
	return expr.LessThanEqual(left, right)
}

// In creates an IN comparison with a list of values
func In(left core.Expression, values ...any) core.Expression {
	return expr.In(left, values...)
}

// StartsWith creates a STARTS WITH comparison
func StartsWith(left core.Expression, value string) core.Expression {
	return expr.StartsWith(left, value)
}

// EndsWith creates an ENDS WITH comparison
func EndsWith(left core.Expression, value string) core.Expression {
	return expr.EndsWith(left, value)
}

// Contains creates a CONTAINS comparison
func Contains(left core.Expression, value string) core.Expression {
	return expr.Contains(left, value)
}

// RegularExpression creates a regular expression comparison
func RegularExpression(left core.Expression, pattern string) core.Expression {
	return expr.RegularExpression(left, pattern)
}

// As creates an alias for an expression
func As(expression core.Expression, alias string) core.Expression {
	return expr.As(expression, alias)
}

// Desc adds DESC ordering for an ORDER BY expression
func Desc(expression core.Expression) core.Expression {
	return expr.Desc(expression)
}

// Asc adds ASC ordering for an ORDER BY expression (default)
func Asc(expression core.Expression) core.Expression {
	return expr.Asc(expression)
}

// Function creates a function call expression
func Function(name string, args ...core.Expression) core.Expression {
	return expr.Function(name, args...)
}

// Count creates a COUNT function expression
func Count(expression core.Expression) core.Expression {
	return expr.Count(expression)
}

// CountStar creates a COUNT(*) function expression
func CountStar() core.Expression {
	return expr.CountStar()
}

// Sum creates a SUM function expression
func Sum(expression core.Expression) core.Expression {
	return expr.Sum(expression)
}

// Avg creates an AVG function expression
func Avg(expression core.Expression) core.Expression {
	return expr.Avg(expression)
}

// Min creates a MIN function expression
func Min(expression core.Expression) core.Expression {
	return expr.Min(expression)
}

// Max creates a MAX function expression
func Max(expression core.Expression) core.Expression {
	return expr.Max(expression)
}

// Utility functions

// convertProperties converts a map of Go values to a map of Expression values
func convertProperties(properties map[string]any) map[string]core.Expression {
	result := make(map[string]core.Expression, len(properties))
	for k, v := range properties {
		result[k] = expr.LiteralFromValue(v)
	}
	return result
}

// Variables
// ================================================================

// Var creates a variable reference
func Var(name string) core.Expression {
	return expr.NewVariableExpression(name)
}

// Alias expressions
// ================================================================
