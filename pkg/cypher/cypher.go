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

// ComplexPath creates a path with multiple relationships in sequence
// This is useful for creating more complex patterns without nesting multiple calls
// Example: ComplexPath(user, "WORKS_FOR", company, "LOCATED_IN", city)
// Creates pattern: (user)-[:WORKS_FOR]->(company)-[:LOCATED_IN]->(city)
func ComplexPath(startNode core.NodeExpression, elementsInPath ...interface{}) core.Expression {
	if len(elementsInPath) < 2 || len(elementsInPath)%2 != 0 {
		panic("ComplexPath requires at least one relationship type and node pair, and must have an even number of elements")
	}

	var elements []core.PatternElement
	currentNode := startNode

	// Add the start node
	elements = append(elements, currentNode)

	// Process the path pairs: [relType, node, relType, node, ...]
	for i := 0; i < len(elementsInPath); i += 2 {
		relType, ok := elementsInPath[i].(string)
		if !ok {
			panic("Relationship type must be a string")
		}

		nextNode, ok := elementsInPath[i+1].(core.NodeExpression)
		if !ok {
			panic("Path elements must alternate between relationship types (string) and nodes")
		}

		// Create the relationship and add it to elements
		rel := currentNode.RelationshipTo(nextNode, relType)
		elements = append(elements, rel)

		// Update current node for next iteration
		currentNode = nextNode
	}

	return ast.Path(elements...)
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

// CompareProperty creates a fluent comparison condition between a property and a parameter value
// This simplifies common property comparisons by automatically creating both the property access
// and parameter in one call.
// Example: CompareProperty("n", "age", ">=", 30) -> n.age >= $param
func CompareProperty(entity, property, operator string, value any) core.Expression {
	propExpr := Property(entity, property)
	paramExpr := Param(value)

	switch operator {
	case "=", "==":
		return Eq(propExpr, paramExpr)
	case "!=", "<>":
		return Ne(propExpr, paramExpr)
	case ">":
		return Gt(propExpr, paramExpr)
	case ">=":
		return Gte(propExpr, paramExpr)
	case "<":
		return Lt(propExpr, paramExpr)
	case "<=":
		return Lte(propExpr, paramExpr)
	case "IN":
		// For IN, value should be a slice or array
		return expr.In(propExpr, value)
	case "CONTAINS":
		// For CONTAINS, value should be a string
		if strVal, ok := value.(string); ok {
			return expr.Contains(propExpr, strVal)
		}
		panic("CONTAINS operator requires a string value")
	case "STARTS WITH":
		// For STARTS WITH, value should be a string
		if strVal, ok := value.(string); ok {
			return expr.StartsWith(propExpr, strVal)
		}
		panic("STARTS WITH operator requires a string value")
	case "ENDS WITH":
		// For ENDS WITH, value should be a string
		if strVal, ok := value.(string); ok {
			return expr.EndsWith(propExpr, strVal)
		}
		panic("ENDS WITH operator requires a string value")
	default:
		panic("Unsupported operator: " + operator)
	}
}

// NamedCompareProperty creates a fluent comparison condition between a property and a named parameter
// Similar to CompareProperty but uses a named parameter
// Example: NamedCompareProperty("n", "age", ">=", "minAge", 30) -> n.age >= $minAge
func NamedCompareProperty(entity, property, operator, paramName string, value any) core.Expression {
	propExpr := Property(entity, property)
	paramExpr := NamedParam(paramName, value)

	switch operator {
	case "=", "==":
		return Eq(propExpr, paramExpr)
	case "!=", "<>":
		return Ne(propExpr, paramExpr)
	case ">":
		return Gt(propExpr, paramExpr)
	case ">=":
		return Gte(propExpr, paramExpr)
	case "<":
		return Lt(propExpr, paramExpr)
	case "<=":
		return Lte(propExpr, paramExpr)
	case "IN":
		// For IN, value should be a slice or array
		return expr.In(propExpr, value)
	case "CONTAINS":
		// For CONTAINS, value should be a string
		if strVal, ok := value.(string); ok {
			return expr.Contains(propExpr, strVal)
		}
		panic("CONTAINS operator requires a string value")
	case "STARTS WITH":
		// For STARTS WITH, value should be a string
		if strVal, ok := value.(string); ok {
			return expr.StartsWith(propExpr, strVal)
		}
		panic("STARTS WITH operator requires a string value")
	case "ENDS WITH":
		// For ENDS WITH, value should be a string
		if strVal, ok := value.(string); ok {
			return expr.EndsWith(propExpr, strVal)
		}
		panic("ENDS WITH operator requires a string value")
	default:
		panic("Unsupported operator: " + operator)
	}
}

// Render renders a statement to a Cypher query string
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

// Collect creates a COLLECT function expression
func Collect(expression core.Expression) core.Expression {
	return expr.Collect(expression)
}

// Distinct wraps an expression with DISTINCT keyword
func Distinct(expression core.Expression) core.Expression {
	return expr.Distinct(expression)
}

// String operators
// ================================================================

// Concat concatenates multiple string expressions using the + operator
func Concat(expressions ...core.Expression) core.Expression {
	return expr.Concat(expressions...)
}

// Substring creates a SUBSTRING function expression
// substring(expression, start [, length])
func Substring(expression core.Expression, start core.Expression, length ...core.Expression) core.Expression {
	return expr.Substring(expression, start, length...)
}

// Replace creates a REPLACE function expression
func Replace(expression, search, replace core.Expression) core.Expression {
	return expr.Replace(expression, search, replace)
}

// Split creates a SPLIT function expression
func Split(expression, delimiter core.Expression) core.Expression {
	return expr.Split(expression, delimiter)
}

// ToLower creates a toLower function expression
func ToLower(expression core.Expression) core.Expression {
	return expr.ToLower(expression)
}

// ToUpper creates a toUpper function expression
func ToUpper(expression core.Expression) core.Expression {
	return expr.ToUpper(expression)
}

// Trim creates a TRIM function expression
func Trim(expression core.Expression) core.Expression {
	return expr.Trim(expression)
}

// LTrim creates a lTrim function expression
func LTrim(expression core.Expression) core.Expression {
	return expr.LTrim(expression)
}

// RTrim creates a rTrim function expression
func RTrim(expression core.Expression) core.Expression {
	return expr.RTrim(expression)
}

// RawCypher creates a raw Cypher expression that will be inserted as-is into the query
// WARNING: Use with caution to avoid Cypher injection vulnerabilities.
// Only use this when the DSL doesn't support a specific Cypher feature.
func RawCypher(cypher string) core.Expression {
	return expr.RawCypher(cypher)
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
