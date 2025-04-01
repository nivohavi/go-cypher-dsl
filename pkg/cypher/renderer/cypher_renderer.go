package renderer

import (
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// CypherRenderer renders a Cypher statement
type CypherRenderer struct {
	pretty       bool
	indentLevel  int
	indentString string
	parameters   *core.Parameters
}

// NewCypherRenderer creates a new CypherRenderer
func NewCypherRenderer() *CypherRenderer {
	return &CypherRenderer{
		pretty:       false,
		indentLevel:  0,
		indentString: "  ",
		parameters:   core.NewParameters(),
	}
}

// WithPrettyPrint enables or disables pretty printing
func (r *CypherRenderer) WithPrettyPrint(pretty bool) *CypherRenderer {
	r.pretty = pretty
	return r
}

// WithIndentString sets the indent string
func (r *CypherRenderer) WithIndentString(indent string) *CypherRenderer {
	r.indentString = indent
	return r
}

// WithParameters sets the parameters object
func (r *CypherRenderer) WithParameters(parameters *core.Parameters) *CypherRenderer {
	r.parameters = parameters
	return r
}

// Render renders a statement
func (r *CypherRenderer) Render(statement core.Statement) string {
	if statement == nil {
		return ""
	}

	// Simple implementation for now, just get the string representation
	cypher := statement.Cypher()

	if r.pretty {
		cypher = r.prettyPrint(cypher)
	}

	return cypher
}

// RenderWithParams renders a statement and returns the Cypher and parameters
func (r *CypherRenderer) RenderWithParams(statement core.Statement) (string, map[string]any) {
	if statement == nil {
		return "", nil
	}

	cypher := r.Render(statement)
	params := statement.Params()

	return cypher, params
}

// prettyPrint formats a Cypher query for better readability
func (r *CypherRenderer) prettyPrint(cypher string) string {
	// A simple implementation for now
	// A more sophisticated implementation would parse the query and format it properly

	// Split by keywords
	keywords := []string{"MATCH", "OPTIONAL MATCH", "WHERE", "WITH", "RETURN", "ORDER BY", "SKIP", "LIMIT", "CREATE", "MERGE", "DELETE", "DETACH DELETE", "SET", "REMOVE", "UNWIND"}
	for _, keyword := range keywords {
		// Replace keywords with newline + keyword
		cypher = strings.ReplaceAll(cypher, " "+keyword+" ", "\n"+keyword+" ")
	}

	// Add proper indentation
	lines := strings.Split(cypher, "\n")
	for i := 1; i < len(lines); i++ {
		lines[i] = r.indentString + lines[i]
	}

	return strings.Join(lines, "\n")
}

// RenderVisitor implements the visitor pattern for rendering expressions
type RenderVisitor struct {
	parameters *core.Parameters
}

// NewRenderVisitor creates a new RenderVisitor
func NewRenderVisitor() *RenderVisitor {
	return &RenderVisitor{
		parameters: core.NewParameters(),
	}
}

// WithParameters sets the parameters object
func (v *RenderVisitor) WithParameters(parameters *core.Parameters) *RenderVisitor {
	v.parameters = parameters
	return v
}

// Visit visits an expression and returns its string representation
func (v *RenderVisitor) Visit(element any) any {
	if element == nil {
		return "NULL"
	}

	// Use the String method of the expression if available
	if expr, ok := element.(core.Expression); ok {
		return expr.String()
	}

	// Use a literal expression for other values
	literal := core.NewLiteral(element)
	return literal.String()
}
