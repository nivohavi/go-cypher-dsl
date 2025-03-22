package main

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/yourusername/go-cypher-dsl/pkg/cypher"
)

func main() {
	// Neo4j connection parameters
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "Password1" // Replace with your actual password

	// Create driver
	driver, err := neo4j.NewDriver(uri,
		neo4j.BasicAuth(username, password, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close()

	// Test the connection
	err = driver.VerifyConnectivity()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Neo4j: %v", err))
	}
	fmt.Println("Successfully connected to Neo4j")

	// Create a session
	session := driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite, // Use AccessModeRead for read-only queries
	})
	defer session.Close()

	// Build a simple example query with go-cypher-dsl
	// Let's find movies from the 1990s
	movie := cypher.Node("Movie").Named("m")
	released := movie.Property("released")

	// Create a condition for movies released between 1990 and 1999
	condition := released.Gte(1990).And(released.Lt(2000))

	// Create aliases for the properties
	titleAlias := cypher.As(movie.Property("title"), "title")
	releasedAlias := cypher.As(released, "year")

	query, _ := cypher.Match(movie).
		Where(condition).
		Returning(titleAlias, releasedAlias).
		OrderBy(released).
		Limit(cypher.Literal(5)).
		Build()

	// View the generated Cypher query
	fmt.Println("\nGenerated Cypher query:")
	fmt.Println(query.Cypher())
	fmt.Println("Parameters:", query.Params())

	// Execute the query
	result, err := session.Run(
		query.Cypher(),
		query.Params(),
	)

	if err != nil {
		panic(err)
	}

	// Process results
	fmt.Println("\nQuery results:")
	for result.Next() {
		record := result.Record()
		title, _ := record.Get("title")
		year, _ := record.Get("year")
		fmt.Printf("Movie: %s (%v)\n", title, year)
	}
}
