package test

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
)

func TestParameterExtraction(t *testing.T) {
	// Test that parameters are correctly extracted from a complex query
	params := cypher.Parameters()

	// Create node patterns
	person := cypher.Node("Person").Named("p")
	movie := cypher.Node("Movie").Named("m")

	// Create a relationship pattern
	acted := person.RelationshipTo(movie, "ACTED_IN")

	// Create a query with multiple parameters
	stmt, err := cypher.Match(cypher.Pattern(acted)).
		Where(
			cypher.And(
				movie.Property("released").Gt(params.Add(2000)),
				person.Property("name").Eq(params.Add("Tom Hanks")),
			),
		).
		Returning(
			movie.Property("title"),
			movie.Property("released"),
		).
		OrderBy(movie.Property("released")).
		Desc().
		Build()

	if err != nil {
		t.Fatalf("Failed to build query: %v", err)
	}

	// Check that parameters were properly extracted
	extractedParams := stmt.Params()
	if len(extractedParams) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(extractedParams))
	}

	// Check values
	var foundYear, foundName bool
	for _, v := range extractedParams {
		switch v {
		case 2000:
			foundYear = true
		case "Tom Hanks":
			foundName = true
		}
	}

	if !foundYear {
		t.Error("Year parameter (2000) not found in extracted parameters")
	}

	if !foundName {
		t.Error("Name parameter ('Tom Hanks') not found in extracted parameters")
	}
}
