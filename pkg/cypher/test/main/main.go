package main

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/util"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/renderer"
)

func main() {
	fmt.Println("Testing util.ExtractParameters...")
	TestExtractParametersDirectly()

	fmt.Println("\nTesting multiple patterns query...")
	TestQueryWithMultiplePatterns()

	fmt.Println("\nTesting formatter...")
	TestFormatter()
}

func TestExtractParametersDirectly() {
	// Create parameter
	params := make(map[string]any)
	param := cypher.ParamWithValue("test", "value")

	// Extract parameters
	util.ExtractParameters(param, params)

	// Verify the parameter was extracted
	fmt.Printf("Extracted parameters count: %d\n", len(params))

	// Check the value
	if val, ok := params["test"]; ok {
		fmt.Printf("Parameter 'test' value: %v\n", val)
	} else {
		fmt.Println("Parameter 'test' not found!")
	}
}

func TestQueryWithMultiplePatterns() {
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
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	// Verify the query was built correctly
	cypherQuery := stmt.Cypher()
	fmt.Printf("Generated Cypher query:\n%s\n", cypherQuery)

	// Check for specific patterns in the query
	expectedPatterns := []string{
		"MATCH",
		"OPTIONAL MATCH",
		"WHERE",
		"RETURN",
	}

	for _, pattern := range expectedPatterns {
		if !contains(cypherQuery, pattern) {
			fmt.Printf("ERROR: Expected pattern '%s' not found in query\n", pattern)
		} else {
			fmt.Printf("Pattern '%s' found in query\n", pattern)
		}
	}
}

func TestFormatter() {
	// Create a simple query
	person := cypher.Node("Person").Named("p")
	movie := cypher.Node("Movie").Named("m")
	actedIn := person.RelationshipTo(movie, "ACTED_IN")

	stmt, err := cypher.Match(cypher.Pattern(person, actedIn, movie)).
		Where(person.Property("name").Eq("Tom Hanks")).
		Returning(movie.Property("title"), movie.Property("released")).
		OrderBy(movie.Property("released")).
		Desc().
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	// Get the query string
	query := stmt.Cypher()
	fmt.Printf("Original query:\n%s\n\n", query)

	// Create formatter with default options
	formatter := renderer.NewDefaultFormatter()
	formatted := formatter.Format(query)
	fmt.Printf("Default formatted query:\n%s\n\n", formatted)

	// Create formatter with custom options
	customFormatter := renderer.NewCypherFormatter(renderer.FormattingOptions{
		IndentString:     "    ", // 4 spaces
		KeywordCase:      renderer.KeywordCaseLower,
		ClauseNewline:    true,
		IndentSubClauses: true,
	})
	customFormatted := customFormatter.Format(query)
	fmt.Printf("Custom formatted query:\n%s\n", customFormatted)
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
