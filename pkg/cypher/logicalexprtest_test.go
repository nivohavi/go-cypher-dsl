package cypher

import (
	"testing"
)

func TestLogicalExpressionsWithNot(t *testing.T) {
	// Define a node
	movie := Node("Movie").Named("m").(*nodePattern)

	// Create the NOT expression
	notExpr := movie.Property("title").Contains("Matrix").Not()

	// Build a query with NOT
	stmt, err := Match(movie).
		Where(notExpr).
		Returning(movie).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (m:`Movie`) WHERE NOT (m.title CONTAINS 'Matrix') RETURN m"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}
