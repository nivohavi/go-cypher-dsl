// Package main demonstrates common Cypher patterns using go-cypher-dsl
package main

import (
	"fmt"

	"github.com/yourusername/go-cypher-dsl/pkg/cypher"
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
	personName := person.Property("name")
	personAge := person.Property("age")

	// Build the query
	stmt, err := cypher.Create(person).
		Set(personName.To(cypher.Literal("John Doe"))).
		Set(personAge.To(cypher.Literal(30))).
		Returning(person).
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

	// Use id function with person as argument
	idFunc := cypher.Function("id", person)

	// Create the statement
	stmt, err := cypher.Match(person).
		Where(idFunc.Eq(cypher.Literal(123))).
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

	// Match and delete in a single statement
	stmt, err := cypher.Match(person).
		Where(person.Property("name").Eq(cypher.Literal("John Doe"))).
		Delete(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Also demonstrate detach delete
	fmt.Println("\nDetach Delete Pattern:")

	stmt2, err := cypher.Match(person).
		Where(person.Property("name").Eq(cypher.Literal("John Doe"))).
		DetachDelete(person).
		Build()

	if err != nil {
		fmt.Printf("Error building detach delete query: %v\n", err)
		return
	}

	fmt.Println(stmt2.Cypher())
	fmt.Println("Parameters:", stmt2.Params())
}

// relationshipPattern demonstrates how to create and query relationships
func relationshipPattern() {
	fmt.Println("\n4. Relationship Pattern:")

	// Find all movies of Tom Hanks (similar to Java Cypher DSL example)
	tom := cypher.Node("Person").Named("tom").WithProperties(map[string]cypher.Expression{
		"name": cypher.Literal("Tom Hanks"),
	})
	tomHanksMovies := cypher.Node("Movie").Named("m")

	// Create the relationship pattern
	actedIn := tom.RelationshipTo(tomHanksMovies, "ACTED_IN")

	// Build the query
	stmt, err := cypher.Match(actedIn).
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
	condition := released.Gte(cypher.Literal(1990)).And(released.Lt(cypher.Literal(2000)))

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
		Limit(cypher.Literal(10)).
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
	tom := cypher.Node("Person").Named("tom").WithProperties(map[string]cypher.Expression{
		"name": cypher.Literal("Tom Hanks"),
	})
	movie := cypher.Node("Movie").Named("m")
	coActors := cypher.Node("Person").Named("coActors")

	// Create the complex relationship pattern
	path := tom.RelationshipTo(movie, "ACTED_IN").RelationshipFrom(coActors, "ACTED_IN")

	// Build the query
	stmt, err := cypher.Match(path).
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
	cloudAtlas := cypher.Node("Movie").Named("cloudAtlas").WithProperties(map[string]cypher.Expression{
		"title": cypher.Literal("Cloud Atlas"),
	})
	people := cypher.Node("Person").Named("people")

	// Create bidirectional relationship
	relatedTo := people.RelationshipBetween(cloudAtlas).Named("relatedTo")

	// Use type function similar to Java's Cypher.type()
	relType := cypher.Function("type", relatedTo)

	// Build the query
	stmt, err := cypher.Match(relatedTo).
		Returning(people.Property("name"), relType, relatedTo).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())
}
