// Package main demonstrates how to connect to Neo4j and execute Cypher queries
// generated with the go-cypher-dsl
package main

import (
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/yourusername/go-cypher-dsl/pkg/cypher"
)

func main() {
	fmt.Println("Neo4j Connection Example")
	fmt.Println("=======================")

	// Neo4j connection parameters
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "Password1" // Replace with your actual password

	// Create a Neo4j driver
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalf("Error creating Neo4j driver: %v", err)
	}
	defer driver.Close()

	// Verify connectivity
	err = driver.VerifyConnectivity()
	if err != nil {
		log.Fatalf("Error connecting to Neo4j: %v", err)
	}
	fmt.Println("Successfully connected to Neo4j")

	// Create a session for reading data
	session := driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close()

	// Example: Find movies released in the 1990s
	findMoviesExample(session)
}

// findMoviesExample demonstrates building and executing a query to find movies
func findMoviesExample(session neo4j.Session) {
	fmt.Println("\nFinding movies released in the 1990s:")

	// Define a movie node and properties
	movie := cypher.Node("Movie").Named("m")
	released := movie.Property("released")

	// Create a condition for movies from the 1990s
	condition := released.Gte(cypher.Literal(1990)).And(released.Lt(cypher.Literal(2000)))

	// Create aliases for better results
	titleAlias := cypher.As(movie.Property("title"), "title")
	yearAlias := cypher.As(released, "year")

	// Build the query
	stmt, err := cypher.Match(movie).
		Where(condition).
		Returning(titleAlias, yearAlias).
		OrderBy(released).
		Limit(cypher.Literal(5)).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	// Display the generated query
	fmt.Println("\nGenerated Cypher query:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query against Neo4j
	result, err := session.Run(
		stmt.Cypher(),
		stmt.Params(),
	)

	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Process and display the results
	fmt.Println("\nQuery results:")
	for result.Next() {
		record := result.Record()
		title, _ := record.Get("title")
		year, _ := record.Get("year")
		fmt.Printf("Movie: %s (%v)\n", title, year)
	}
}
