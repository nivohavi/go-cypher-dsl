package test

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/renderer"
)

func TestFormatter(t *testing.T) {
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
		t.Fatalf("Failed to build query: %v", err)
	}

	// Get the query string
	query := stmt.Cypher()

	// Create formatter with default options
	formatter := renderer.NewDefaultFormatter()
	formatted := formatter.Format(query)

	// Verify that formatting works and doesn't change the fundamental query
	if formatted == "" {
		t.Fatal("Formatted query is empty")
	}

	// Create formatter with custom options
	customFormatter := renderer.NewCypherFormatter(renderer.FormattingOptions{
		IndentString:     "    ", // 4 spaces
		KeywordCase:      renderer.KeywordCaseLower,
		ClauseNewline:    true,
		IndentSubClauses: true,
	})
	customFormatted := customFormatter.Format(query)

	if customFormatted == "" {
		t.Fatal("Custom formatted query is empty")
	}

	// Verify lowercase keywords in the custom formatted query
	if len(customFormatted) > 0 && customFormatted[0] == 'M' {
		t.Errorf("Expected lowercase 'match' but got uppercase: %s", customFormatted[:5])
	}
}
