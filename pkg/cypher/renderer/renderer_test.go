package renderer

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

func TestNewCypherRenderer(t *testing.T) {
	renderer := NewCypherRenderer()
	if renderer == nil {
		t.Error("NewCypherRenderer() returned nil")
	}
}

func TestRender(t *testing.T) {
	stmt := core.NewStatement("MATCH (n) RETURN n", nil)
	renderer := NewCypherRenderer()
	result := renderer.Render(stmt)

	if result != "MATCH (n) RETURN n" {
		t.Errorf("Render() = %q, want 'MATCH (n) RETURN n'", result)
	}
}

func TestRenderWithNilStatement(t *testing.T) {
	renderer := NewCypherRenderer()
	result := renderer.Render(nil)

	if result != "" {
		t.Errorf("Render(nil) = %q, want ''", result)
	}
}

func TestRenderWithPrettyPrint(t *testing.T) {
	stmt := core.NewStatement("MATCH (n) RETURN n", nil)
	renderer := NewCypherRenderer().WithPrettyPrint(true)
	result := renderer.Render(stmt)

	// Should contain newlines when pretty printed
	if !strings.Contains(result, "\n") {
		t.Errorf("Render() with pretty print = %q, should contain newlines", result)
	}
}

func TestRenderWithParams(t *testing.T) {
	params := map[string]any{"key": "value"}
	stmt := core.NewStatement("MATCH (n) RETURN n", params)
	renderer := NewCypherRenderer()
	cypher, resultParams := renderer.RenderWithParams(stmt)

	if cypher != "MATCH (n) RETURN n" {
		t.Errorf("RenderWithParams() cypher = %q, want 'MATCH (n) RETURN n'", cypher)
	}
	if resultParams["key"] != "value" {
		t.Errorf("RenderWithParams() params = %v, want map with key 'value'", resultParams)
	}
}

func TestRenderWithParamsNilStatement(t *testing.T) {
	renderer := NewCypherRenderer()
	cypher, params := renderer.RenderWithParams(nil)

	if cypher != "" {
		t.Errorf("RenderWithParams(nil) cypher = %q, want ''", cypher)
	}
	if params != nil {
		t.Errorf("RenderWithParams(nil) params = %v, want nil", params)
	}
}

func TestWithIndentString(t *testing.T) {
	stmt := core.NewStatement("MATCH (n) RETURN n", nil)
	renderer := NewCypherRenderer().WithPrettyPrint(true).WithIndentString("    ")
	result := renderer.Render(stmt)

	// Should contain the custom indent string
	if !strings.Contains(result, "    ") {
		t.Errorf("Render() with custom indent = %q, should contain '    '", result)
	}
}

