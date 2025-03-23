// Package main demonstrates common Cypher patterns using go-cypher-dsl
package main

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
)

func main() {
	fmt.Println("Common Cypher Pattern Examples")
	fmt.Println("=============================")

	// Demonstrate different patterns
	createNodePattern()
	matchNodeByIdPattern()
	deleteNodePattern()
	relationshipPattern()
	parameterizedQueriesPattern()
	aggregationPattern()
	coActorsPattern()
	movieRelationshipsPattern()
}

// createNodePattern demonstrates how to create a node
func createNodePattern() {
	fmt.Println("\n1. Create Node Pattern:")

	// Create a person node with properties
	person := cypher.Node("Person").Named("p")

	// Set properties using a map - using WithProps for automatic conversion
	personWithProps := person.WithProps(map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	})

	// Build the query
	stmt, err := cypher.Create(personWithProps).
		Returning(personWithProps).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())
}

// matchNodeByIdPattern demonstrates how to match a node by ID
func matchNodeByIdPattern() {
	fmt.Println("\n2. Match Node By ID Pattern:")

	// Match a node by ID similar to Java Cypher DSL
	person := cypher.Node("Person").Named("p")

	// Create an expression that compares the ID function to a literal value
	// The ID function returns an Expression that needs to be compared as a standalone predicate
	idEq := person.Property("id").Eq(123)

	// Create the statement
	stmt, err := cypher.Match(person).
		Where(idEq).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())
}

// deleteNodePattern demonstrates how to delete a node
func deleteNodePattern() {
	fmt.Println("\n3. Delete Node Pattern:")

	// Delete a node with properties
	person := cypher.Node("Person").Named("p")

	// Create a delete statement directly (without a preceding MATCH)
	// This is just for demonstration purposes
	deleteStmt, err := cypher.Delete(person).Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	// Just show the delete part for demonstration
	fmt.Println("DELETE part:")
	fmt.Println(deleteStmt.Cypher())
	fmt.Println("Parameters:", deleteStmt.Params())

	// Also demonstrate detach delete
	fmt.Println("\nDetach Delete Pattern:")

	// Create a separate detach delete statement
	detachDeleteStmt, err := cypher.DetachDelete(person).Build()

	if err != nil {
		fmt.Printf("Error building detach delete query: %v\n", err)
		return
	}

	// Just show the detach delete part for demonstration
	fmt.Println("DETACH DELETE part:")
	fmt.Println(detachDeleteStmt.Cypher())
	fmt.Println("Parameters:", detachDeleteStmt.Params())
}

// relationshipPattern demonstrates how to create and query relationships
func relationshipPattern() {
	fmt.Println("\n4. Relationship Pattern:")

	// Find all movies of Tom Hanks using improved API style
	tom := cypher.Node("Person").Named("tom")
	tomHanksMovies := cypher.Node("Movie").Named("m")

	// Create relationship between tom and movies
	actedInRel := tom.RelationshipTo(tomHanksMovies, "ACTED_IN")

	// Create a pattern with the relationship
	actedInPattern := cypher.Pattern(tom, actedInRel, tomHanksMovies)

	// Build the query with property condition
	stmt, err := cypher.Match(actedInPattern).
		Where(tom.Property("name").Eq("Tom Hanks")).
		Returning(tom, tomHanksMovies).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())
}

// parameterizedQueriesPattern demonstrates how to use parameters
func parameterizedQueriesPattern() {
	fmt.Println("\n5. Parameterized Queries Pattern:")

	// Find movies released in the 1990s (Java Cypher DSL Example 7)
	nineties := cypher.Node("Movie").Named("nineties")
	released := nineties.Property("released")

	// Create condition like the Java example
	condition := released.Gte(1990).And(released.Lt(2000))

	// Build the query
	stmt, err := cypher.Match(nineties).
		Where(condition).
		Returning(nineties.Property("title")).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())
}

// aggregationPattern demonstrates how to use aggregation functions
func aggregationPattern() {
	fmt.Println("\n6. Aggregation Pattern:")

	// Find people and limit results like Java Cypher DSL Example 6
	people := cypher.Node("Person").Named("people")

	// Build the query
	stmt, err := cypher.Match(people).
		Returning(people.Property("name")).
		Limit(10).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())
}

// coActorsPattern demonstrates finding co-actors (Java Cypher DSL Example 10)
func coActorsPattern() {
	fmt.Println("\n7. Co-actors Pattern:")

	// Tom Hanks' co-actors pattern
	tom := cypher.Node("Person").Named("tom")
	movie := cypher.Node("Movie").Named("m")
	coActors := cypher.Node("Person").Named("coActors")

	// Create the complex relationship pattern
	tomToMovie := tom.RelationshipTo(movie, "ACTED_IN")
	movieToCoActors := movie.RelationshipFrom(coActors, "ACTED_IN")

	// Create the path pattern with both relationships
	path := cypher.Pattern(tom, tomToMovie, movie, movieToCoActors, coActors)

	// Build the query
	stmt, err := cypher.Match(path).
		Where(tom.Property("name").Eq("Tom Hanks")).
		Returning(coActors.Property("name")).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())
}

// movieRelationshipsPattern demonstrates finding relationships to a movie (Java Cypher DSL Example 11)
func movieRelationshipsPattern() {
	fmt.Println("\n8. Movie Relationships Pattern:")

	// How people are related to Cloud Atlas
	cloudAtlas := cypher.Node("Movie").Named("cloudAtlas")
	people := cypher.Node("Person").Named("people")

	// Create bidirectional relationship
	relatedTo := people.RelationshipBetween(cloudAtlas).Named("relatedTo")

	// Create the pattern with the relationship
	relatedPattern := cypher.Pattern(people, relatedTo, cloudAtlas)

	// Use type function to get relationship type
	relType := cypher.Function("type", relatedTo)

	// Build the query
	stmt, err := cypher.Match(relatedPattern).
		Where(cloudAtlas.Property("title").Eq("Cloud Atlas")).
		Returning(people.Property("name"), relType, relatedTo).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())
}
