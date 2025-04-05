package basics

import (
	"fmt"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
)

func main() {
	fmt.Println("Go Cypher DSL - Basic Query Building Examples")
	fmt.Println("===========================================")

	// Demonstrate different query patterns
	basicMatchReturn()
	propertyConditions()
	relationshipPatterns()
	complexLogicalConditions()
	withClauseAndPagination()
	createNodes()
	mergePattern()
	deletePattern()
}

// Basic MATCH-RETURN query
func basicMatchReturn() {
	fmt.Println("\n1. Basic MATCH-RETURN Query")
	fmt.Println("-------------------------")

	// Create a node pattern
	person := cypher.Node("Person").Named("p")

	// Build a query with MATCH and RETURN
	stmt, err := cypher.Match(person).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// You can also return specific properties
	stmt2, err := cypher.Match(person).
		Returning(
			// Using an alias for the returned property
			cypher.As(person.Property("name"), "personName"),
			person.Property("born"),
		).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nReturning specific properties:")
	fmt.Println(stmt2.Cypher())
}

// Property condition examples
func propertyConditions() {
	fmt.Println("\n2. Property Conditions")
	fmt.Println("--------------------")

	// Create a node pattern
	person := cypher.Node("Person").Named("p")

	// Build a query with WHERE condition
	stmt, err := cypher.Match(person).
		// Direct property condition on the node
		Where(person.Property("name").Eq("Tom Hanks")).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("Simple property condition:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Numeric comparison
	stmt2, err := cypher.Match(person).
		Where(person.Property("born").Gt(1960)).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nNumeric comparison:")
	fmt.Println(stmt2.Cypher())
	fmt.Println("Parameters:", stmt2.Params())

	// String operations
	stmt3, err := cypher.Match(person).
		Where(person.Property("name").Contains("Tom")).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nString operation:")
	fmt.Println(stmt3.Cypher())
	fmt.Println("Parameters:", stmt3.Params())
}

// Relationship pattern examples
func relationshipPatterns() {
	fmt.Println("\n3. Relationship Patterns")
	fmt.Println("----------------------")

	// Create node patterns
	person := cypher.Node("Person").Named("p")
	movie := cypher.Node("Movie").Named("m")

	// Create a relationship from person to movie
	acted := person.RelationshipTo(movie, "ACTED_IN").Named("r")

	// Create a pattern using the relationship
	pattern := cypher.Path(person, acted, movie)

	// Build a query with the pattern
	stmt, err := cypher.Match(pattern).
		Where(person.Property("name").Eq("Tom Hanks")).
		Returning(movie.Property("title"), movie.Property("released")).
		OrderBy(movie.Property("released")).
		Desc().
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("Actor-Movie relationship:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Complex path with multiple relationships
	// Create another node and relationship
	director := cypher.Node("Person").Named("d")
	directed := movie.RelationshipTo(director, "DIRECTED_BY").Named("dr")

	// Build a query with multiple patterns
	stmt2, err := cypher.Match(pattern).
		OptionalMatch(cypher.Path(movie, directed, director)).
		Where(person.Property("name").Eq("Tom Hanks")).
		Returning(
			movie.Property("title"),
			cypher.As(director.Property("name"), "director"),
		).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nMultiple relationships:")
	fmt.Println(stmt2.Cypher())
	fmt.Println("Parameters:", stmt2.Params())

	// Using ComplexPath helper for multi-level relationships
	user := cypher.Node("User").Named("u")
	company := cypher.Node("Company").Named("c")
	city := cypher.Node("City").Named("city")

	// Build a query with ComplexPath
	stmt3, err := cypher.Match(
		cypher.ComplexPath(
			user,
			"WORKS_AT", company,
			"LOCATED_IN", city,
		),
	).
		Where(user.Property("name").Eq("John")).
		Returning(
			user.Property("name"),
			company.Property("name"),
			city.Property("name"),
		).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nUsing ComplexPath helper:")
	fmt.Println(stmt3.Cypher())
	fmt.Println("Parameters:", stmt3.Params())
}

// Complex logical condition examples
func complexLogicalConditions() {
	fmt.Println("\n4. Complex Logical Conditions")
	fmt.Println("---------------------------")

	// Create node pattern
	movie := cypher.Node("Movie").Named("m")

	// AND condition
	andCondition := movie.Property("released").Gt(2000).
		And(movie.Property("released").Lt(2010))

	stmt, err := cypher.Match(movie).
		Where(andCondition).
		Returning(movie.Property("title"), movie.Property("released")).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("AND condition:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// OR condition
	orCondition := movie.Property("title").Contains("Matrix").
		Or(movie.Property("title").Contains("Inception"))

	stmt2, err := cypher.Match(movie).
		Where(orCondition).
		Returning(movie.Property("title")).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nOR condition:")
	fmt.Println(stmt2.Cypher())
	fmt.Println("Parameters:", stmt2.Params())

	// NOT condition
	notCondition := movie.Property("title").Contains("Matrix").Not()

	stmt3, err := cypher.Match(movie).
		Where(notCondition).
		Returning(movie.Property("title")).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nNOT condition:")
	fmt.Println(stmt3.Cypher())
	fmt.Println("Parameters:", stmt3.Params())

	// Complex condition combination
	complexCondition := movie.Property("released").Gt(2000).
		And(movie.Property("released").Lt(2010)).
		Or(movie.Property("title").Contains("Matrix").Not())

	stmt4, err := cypher.Match(movie).
		Where(complexCondition).
		Returning(movie.Property("title"), movie.Property("released")).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nComplex condition combination:")
	fmt.Println(stmt4.Cypher())
	fmt.Println("Parameters:", stmt4.Params())

	// Using comparison helpers
	stmt5, err := cypher.Match(movie).
		Where(
			cypher.And(
				cypher.CompareProperty("m", "released", ">", 2000),
				cypher.CompareProperty("m", "released", "<", 2010),
			),
		).
		Returning(movie.Property("title")).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nUsing comparison helpers:")
	fmt.Println(stmt5.Cypher())
	fmt.Println("Parameters:", stmt5.Params())
}

// WITH clause and pagination examples
func withClauseAndPagination() {
	fmt.Println("\n5. WITH Clause and Pagination")
	fmt.Println("--------------------------")

	// Create node patterns
	person := cypher.Node("Person").Named("p")
	movie := cypher.Node("Movie").Named("m")
	acted := person.RelationshipTo(movie, "ACTED_IN")

	// Using WITH to chain queries
	stmt, err := cypher.Match(person).
		Where(person.Property("name").Eq("Tom Hanks")).
		With(person).
		Match(cypher.Path(person, acted, movie)).
		Returning(movie.Property("title")).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("Using WITH clause:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Pagination example with SKIP and LIMIT
	stmt2, err := cypher.Match(movie).
		Returning(movie.Property("title")).
		OrderBy(movie.Property("released")).
		Desc().
		Skip(10).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nPagination with SKIP and LIMIT:")
	fmt.Println(stmt2.Cypher())
	fmt.Println("Parameters:", stmt2.Params())
}

// Examples of creating nodes
func createNodes() {
	fmt.Println("\n6. Creating Nodes")
	fmt.Println("---------------")

	// Create a person node with properties
	person := cypher.Node("Person").Named("p")

	// Set properties using WithProps
	personWithProps := person.WithProps(map[string]interface{}{
		"name": "John Doe",
		"born": 1990,
	})

	// Build CREATE query
	stmt, err := cypher.Create(personWithProps).
		Returning(personWithProps).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("Creating a node with properties:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Create a relationship between nodes
	movieNode := cypher.Node("Movie").Named("m")
	movieWithProps := movieNode.WithProps(map[string]interface{}{
		"title":    "Example Movie",
		"released": 2023,
	})

	actedIn := person.RelationshipTo(movieNode, "ACTED_IN").WithProps(map[string]interface{}{
		"roles": []string{"Main Character"},
	})

	// Build CREATE query for relationship
	stmt2, err := cypher.Match(person).
		Where(person.Property("name").Eq("John Doe")).
		Create(cypher.Path(person, actedIn, movieWithProps)).
		Returning(person, movieNode).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nCreating a relationship:")
	fmt.Println(stmt2.Cypher())
	fmt.Println("Parameters:", stmt2.Params())
}

// Example of MERGE pattern
func mergePattern() {
	fmt.Println("\n7. MERGE Pattern")
	fmt.Println("--------------")

	// Define node to be merged
	person := cypher.Node("Person").Named("p").WithProps(map[string]interface{}{
		"name": "Keanu Reeves",
	})

	// MERGE query with ON CREATE SET and ON MATCH SET
	stmt, err := cypher.Merge(person).
		OnCreate(cypher.Eq(
			cypher.Property("p", "created"),
			cypher.ParamWithValue("timestamp", 1617234567),
		)).
		OnMatch(cypher.Eq(
			cypher.Property("p", "updated"),
			cypher.ParamWithValue("timestamp", 1617234567),
		)).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("MERGE with ON CREATE and ON MATCH:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())
}

// Example of DELETE pattern
func deletePattern() {
	fmt.Println("\n8. DELETE Pattern")
	fmt.Println("--------------")

	// Define node to be deleted
	person := cypher.Node("Person").Named("p")

	// DELETE query
	stmt, err := cypher.Match(person).
		Where(person.Property("name").Eq("John Doe")).
		Delete(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("Simple DELETE:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// DETACH DELETE query
	stmt2, err := cypher.Match(person).
		Where(person.Property("name").Eq("John Doe")).
		DetachDelete(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\nDETACH DELETE (deletes relationships too):")
	fmt.Println(stmt2.Cypher())
	fmt.Println("Parameters:", stmt2.Params())
}
