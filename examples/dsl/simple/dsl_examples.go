// Package main demonstrates simple usage patterns of the go-cypher-dsl
package main

import (
	"fmt"

	"github.com/yourusername/go-cypher-dsl/pkg/cypher"
)

func main() {
	fmt.Println("Go Cypher DSL - Simple Usage Examples")
	fmt.Println("====================================")

	// Run the examples
	basicMatchExample()
	propertyMatchExample()
	relationshipExample()
	parameterizedExample()
	aggregationExample()
}

// basicMatchExample demonstrates a very simple MATCH-RETURN query
func basicMatchExample() {
	fmt.Println("\n1. Basic MATCH-RETURN:")

	// Define a node
	person := cypher.Node("Person").Named("p")

	// Build a MATCH ... RETURN query
	stmt, err := cypher.Match(person).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	// Get the Cypher query
	query := stmt.Cypher()
	params := stmt.Params()

	fmt.Println(query)
	fmt.Printf("Parameters: %v\n", params)
}

// propertyMatchExample demonstrates matching with property expressions
func propertyMatchExample() {
	fmt.Println("\n2. Property Match Example:")

	// Define a node
	person := cypher.Node("Person").Named("p")

	// Build a query with property condition
	stmt, err := cypher.Match(person).
		Where(person.Property("name").Eq(cypher.Literal("John Doe"))).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	// Get the Cypher query
	query := stmt.Cypher()
	params := stmt.Params()

	fmt.Println(query)
	fmt.Printf("Parameters: %v\n", params)
}

// relationshipExample demonstrates working with relationships
func relationshipExample() {
	fmt.Println("\n3. Relationship Example:")

	// Define nodes and a relationship
	person := cypher.Node("Person").Named("p")
	movie := cypher.Node("Movie").Named("m")
	actedIn := person.RelationshipTo(movie, "ACTED_IN").Named("r")

	// Create a pattern with the relationship
	pattern := cypher.Pattern(person, actedIn, movie)

	// Build a query using the pattern
	stmt, err := cypher.Match(pattern).
		Where(person.Property("name").Eq(cypher.Literal("Tom Hanks"))).
		Returning(movie).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	// Get the Cypher query
	query := stmt.Cypher()
	params := stmt.Params()

	fmt.Println(query)
	fmt.Printf("Parameters: %v\n", params)
}

// parameterizedExample demonstrates using parameters
func parameterizedExample() {
	fmt.Println("\n4. Parameterized Query Example:")

	// Define nodes
	movie := cypher.Node("Movie").Named("m")

	// Build a query with parameters
	stmt, err := cypher.Match(movie).
		Where(movie.Property("released").Gt(cypher.ParamWithValue("year", 2000))).
		Returning(movie.Property("title"), movie.Property("released")).
		OrderBy(movie.Property("released")).
		Limit(cypher.Literal(5)).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	// Get the Cypher query
	query := stmt.Cypher()
	params := stmt.Params()

	fmt.Println(query)
	fmt.Printf("Parameters: %v\n", params)
}

// aggregationExample demonstrates using aggregation
func aggregationExample() {
	fmt.Println("\n5. Aggregation Example:")

	// Define nodes and a relationship
	person := cypher.Node("Person").Named("p")
	movie := cypher.Node("Movie").Named("m")
	actedIn := person.RelationshipTo(movie, "ACTED_IN").Named("r")

	// Create a pattern with the relationship
	pattern := cypher.Pattern(person, actedIn, movie)

	// Create a count expression
	movieCount := cypher.Function("count", movie)
	movieCountAlias := cypher.As(movieCount, "movieCount")

	// Build a query using aggregation
	stmt, err := cypher.Match(pattern).
		With(person, movieCountAlias).
		Returning(person.Property("name"), movieCountAlias).
		OrderBy(person.Property("name")).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	// Get the Cypher query
	query := stmt.Cypher()
	params := stmt.Params()

	fmt.Println(query)
	fmt.Printf("Parameters: %v\n", params)
}
