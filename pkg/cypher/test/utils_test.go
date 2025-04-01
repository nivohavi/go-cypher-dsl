package test

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/util"
)

func TestExtractParametersDirectly(t *testing.T) {
	// Create parameter
	params := make(map[string]any)
	param := cypher.ParamWithValue("test", "value")

	// Extract parameters
	util.ExtractParameters(param, params)

	// Verify the parameter was extracted
	if len(params) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(params))
	}

	// Check the value
	if val, ok := params["test"]; !ok || val != "value" {
		t.Errorf("Parameter 'test' not found or has wrong value. Got: %v", val)
	}
}

func TestQueryWithMultiplePatterns(t *testing.T) {
	// Create nodes and relationships for a more complex query
	person := cypher.Node("Person").Named("p")
	movie := cypher.Node("Movie").Named("m")
	actedIn := person.RelationshipTo(movie, "ACTED_IN")

	director := cypher.Node("Person").Named("d")
	directedRel := movie.RelationshipTo(director, "DIRECTED_BY")

	// Build a query with multiple patterns
	stmt, err := cypher.Match(cypher.Pattern(actedIn)).
		OptionalMatch(cypher.Pattern(directedRel)).
		Where(
			cypher.And(
				person.Property("name").Eq("Tom Hanks"),
				cypher.Ne(director.Property("name"), cypher.Literal("Steven Spielberg")),
			),
		).
		Returning(
			movie.Property("title"),
			director.Property("name"),
		).
		Build()

	if err != nil {
		t.Fatalf("Failed to build query: %v", err)
	}

	// Verify the query was built correctly
	cypher := stmt.Cypher()
	if cypher == "" {
		t.Fatal("Generated Cypher is empty")
	}

	// Check for specific patterns in the query
	expectedPatterns := []string{
		"MATCH",
		"OPTIONAL MATCH",
		"WHERE",
		"RETURN",
	}

	for _, pattern := range expectedPatterns {
		if !contains(cypher, pattern) {
			t.Errorf("Expected pattern '%s' not found in: %s", pattern, cypher)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != "" && substr != "" && s != substr && s[:len(s)-(len(s)-len(substr))] == substr
}
