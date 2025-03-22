package cypher

import (
	"fmt"
	"testing"
)

func TestWithClauseDebug(t *testing.T) {
	// Define nodes
	person := Node("Person").Named("p").(*nodePattern)
	movie := Node("Movie").Named("m").(*nodePattern)

	// Build a query with WITH
	matchBuilder := Match(person)
	matchWithWhere := matchBuilder.Where(person.Property("name").Eq("Tom Hanks"))
	withBuilder := matchWithWhere.With(person)
	matchBuilderAfterWith := withBuilder.Match(movie)
	matchBuilderWithWhere := matchBuilderAfterWith.Where(movie.Property("title").Contains(person.Property("name")))
	returnBuilder := matchBuilderWithWhere.Returning(movie)

	// Check all the builders
	fmt.Printf("matchBuilder.where: %v\n", matchBuilder.where != nil)
	fmt.Printf("matchWithWhere.where: %v\n", matchWithWhere.where != nil)
	fmt.Printf("matchBuilderAfterWith.where: %v\n", matchBuilderAfterWith.where != nil)
	fmt.Printf("matchBuilderWithWhere.where: %v\n", matchBuilderWithWhere.where != nil)

	// Build the statement
	stmt, err := returnBuilder.Build()
	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	// Debug the cypherStatement
	cypherStmt := stmt.(*cypherStatement)
	fmt.Printf("builder type: %T\n", cypherStmt.builder)

	// Render the Cypher
	actual := stmt.Cypher()
	expected := "MATCH (p:`Person`) WHERE p.name = 'Tom Hanks' WITH p MATCH (m:`Movie`) WHERE m.title CONTAINS p.name RETURN m"

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}
