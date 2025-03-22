package main

import (
	"fmt"
	"log"

	"github.com/yourusername/go-cypher-dsl/pkg/cypher"
)

func main() {
	// Example 1: Basic MATCH-RETURN query
	example1()

	// Example 2: Query with parameters
	example2()

	// Example 3: More complex query with relationships
	example3()

	// Example 4: Using the static schema
	example4()
}

func example1() {
	fmt.Println("Example 1: Basic MATCH-RETURN query")

	// Define a node
	person := cypher.Node("Person").Named("p")

	// Build a MATCH ... RETURN query
	stmt, err := cypher.Match(person).
		Returning(person).
		Build()

	if err != nil {
		log.Fatalf("Error building query: %v", err)
	}

	// Get the Cypher query
	query := stmt.Cypher()
	fmt.Println(query)
	fmt.Println()

	// This will output:
	// MATCH (p:`Person`) RETURN p
}

func example2() {
	fmt.Println("Example 2: Query with parameters")

	// Create a node with the ability to access properties
	personBase := cypher.Node("Person")
	person := personBase.Named("p")

	// Type assertion only for Property method
	nodeWithProps, ok := personBase.(interface {
		Property(name string) cypher.PropertyExpression
	})
	if !ok {
		log.Fatal("Could not assert node type")
	}

	// Build a MATCH ... WHERE ... RETURN query with parameter
	stmt, err := cypher.Match(person).
		Where(nodeWithProps.Property("name").Eq(cypher.ParamWithValue("name", "Tom Hanks"))).
		Returning(person).
		Build()

	if err != nil {
		log.Fatalf("Error building query: %v", err)
	}

	// Get the Cypher query and parameters
	query := stmt.Cypher()
	params := stmt.Params()

	fmt.Println(query)
	fmt.Printf("Parameters: %v\n", params)
	fmt.Println()

	// This will output:
	// MATCH (p:`Person`) WHERE p.name = $name RETURN p
	// Parameters: map[name:Tom Hanks]
}

func example3() {
	fmt.Println("Example 3: More complex query with relationships")

	// Create nodes and relationships
	personBase := cypher.Node("Person")
	person := personBase.Named("p")

	movieBase := cypher.Node("Movie")
	movie := movieBase.Named("m")

	acted := cypher.Relationship("ACTED_IN").Named("r").Build()

	// Type assertions for property access
	personWithProps, ok := personBase.(interface {
		Property(name string) cypher.PropertyExpression
	})
	if !ok {
		log.Fatal("Could not assert person type")
	}

	movieWithProps, ok := movieBase.(interface {
		Property(name string) cypher.PropertyExpression
	})
	if !ok {
		log.Fatal("Could not assert movie type")
	}

	// Build a pattern path (person)-[acted]->(movie)
	path := cypher.Pattern(person, acted, movie)

	// Create a logical AND condition using And method from ComparisonExpression
	nameCondition := personWithProps.Property("name").Eq(cypher.Literal("Keanu Reeves"))
	yearCondition := movieWithProps.Property("released").Gt(cypher.Literal(2000))

	// Build a MATCH query with the path and conditions
	stmt, err := cypher.Match(path).
		Where(nameCondition.And(yearCondition)).
		Returning(
			movieWithProps.Property("title"),
			cypher.As(movieWithProps.Property("released"), "year"),
			acted,
		).
		OrderBy(movieWithProps.Property("released")).
		Build()

	if err != nil {
		log.Fatalf("Error building query: %v", err)
	}

	// Get the Cypher query
	query := stmt.Cypher()
	fmt.Println(query)
	fmt.Println()

	// This will output:
	// MATCH (p:`Person`)-[r:`ACTED_IN`]->(m:`Movie`)
	// WHERE (p.name = 'Keanu Reeves' AND m.released > 2000)
	// RETURN m.title, m.released AS year, r ORDER BY m.released
}

func example4() {
	fmt.Println("Example 4: Using the static schema")

	// Use the static schema
	person := cypher.Person()
	personNode := person.Named("p")
	movie := cypher.Movie()
	movieNode := movie.Named("m")
	actedIn := cypher.ActedIn()
	actedRel := actedIn.Direction(cypher.OUTGOING).Named("r").Build()

	// Build a query with the static schema
	stmt, err := cypher.Match(cypher.Pattern(personNode, actedRel, movieNode)).
		Where(person.Name().Eq(cypher.Literal("Keanu Reeves"))).
		Returning(movie.Title(), actedIn.Role()).
		Build()

	if err != nil {
		log.Fatalf("Error building query: %v", err)
	}

	// Get the Cypher query
	query := stmt.Cypher()
	fmt.Println(query)

	// This will output:
	// MATCH (p:`Person`)-[r:`ACTED_IN`]->(m:`Movie`)
	// WHERE p.name = 'Keanu Reeves'
	// RETURN m.title, r.role
}
