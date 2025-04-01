package main

import (
	"fmt"
	"os"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

func main() {
	fmt.Println("Go Cypher DSL Example")
	fmt.Println("---------------------")

	// Create parameters collection
	params := cypher.Parameters()

	// Create node patterns
	person := cypher.Node("Person").Named("p")
	movie := cypher.Node("Movie").Named("m")

	// Create a relationship pattern
	acted := cypher.RelateNodes(person, movie, "ACTED_IN")

	// Build a simple query using the fluent API
	stmt, err := cypher.Match(cypher.Pattern(acted)).
		Where(
			cypher.And(
				cypher.Gt(
					cypher.Property("m", "released"),
					params.Add(2000),
				),
				cypher.Eq(
					cypher.Property("p", "name"),
					params.Add("Tom Hanks"),
				),
			),
		).
		Returning(
			cypher.Property("m", "title"),
			cypher.Property("m", "released"),
		).
		OrderBy(cypher.Property("m", "released")).
		Desc().
		Build()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building query 1: %v\n", err)
		os.Exit(1)
	}

	// Print the pretty query
	fmt.Println("Example 1: Find Tom Hanks movies released after 2000")
	fmt.Println(cypher.PrettyPrint(stmt))
	fmt.Printf("\nParameters: %+v\n", stmt.Params())

	// Another example with a more complex pattern
	director := cypher.Node("Person").Named("d")
	directedRel := cypher.RelateNodes(movie, director, "DIRECTED_BY")

	// Build a query with multiple patterns
	stmt2, err := cypher.Match(cypher.Pattern(acted)).
		OptionalMatch(cypher.Pattern(directedRel)).
		Where(
			cypher.And(
				cypher.Eq(
					cypher.Property("p", "name"),
					params.Add("Tom Hanks"),
				),
				cypher.Ne(
					cypher.Property("d", "name"),
					params.Add("Steven Spielberg"),
				),
			),
		).
		Returning(
			cypher.Property("m", "title"),
			cypher.Property("d", "name"),
		).
		OrderBy(cypher.Property("m", "title")).
		Limit(10).
		Build()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building query 2: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nExample 2: Find movies where Tom Hanks acted but not directed by Spielberg")
	fmt.Println(cypher.PrettyPrint(stmt2))
	fmt.Printf("\nParameters: %+v\n", stmt2.Params())

	// Example with CREATE and SET
	// Use a single parameters container for all parameters
	params3 := cypher.Parameters()
	newPerson := cypher.Node("Person").Named("newPerson")
	newPersonWithProps := newPerson.WithProperties(map[string]core.Expression{
		"name": params3.Add("Leonardo DiCaprio"),
		"born": params3.Add(1974),
	})

	titanic := cypher.Node("Movie").Named("titanic")
	titanicWithProps := titanic.WithProperties(map[string]core.Expression{
		"title": params3.Add("Titanic"),
	})

	actedInRel := cypher.RelateNodes(newPerson, titanic, "ACTED_IN")
	oscarValue := params3.Add(1)

	stmt3, err := cypher.Create(cypher.Pattern(newPersonWithProps)).
		With(newPerson).
		Match(cypher.Pattern(titanicWithProps)).
		Create(cypher.Pattern(actedInRel)).
		Set(cypher.Property("newPerson", "oscars")).
		And(oscarValue).
		Returning(
			cypher.Property("newPerson", "name"),
			cypher.Property("titanic", "title"),
		).
		Build()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building query 3: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nExample 3: Create a person and relationship")
	fmt.Println(cypher.PrettyPrint(stmt3))
	fmt.Printf("\nParameters: %+v\n", stmt3.Params())
}
