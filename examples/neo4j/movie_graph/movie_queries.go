// Package main demonstrates how to use go-cypher-dsl with the Neo4j movie graph dataset
package main

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
)

func main() {
	fmt.Println("Neo4j Movie Graph Examples")
	fmt.Println("==========================")

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
		AccessMode: neo4j.AccessModeRead, // Using read mode for these examples
	})
	defer session.Close()

	// Run each example
	runExample1(session)
	runExample2(session)
	runExample3(session)
	runExample4(session)
	runExample5(session)
}

// runExample1 demonstrates a simple MATCH-RETURN query
func runExample1(session neo4j.Session) {
	// Example 1: Simple MATCH-RETURN query
	// Notice no type assertion needed
	person := cypher.Node("Person").Named("p")

	stmt1, err := cypher.Match(person).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n1. Simple MATCH-RETURN:")
	fmt.Println(stmt1.Cypher())
	fmt.Println("Parameters:", stmt1.Params())

	// Execute the query
	result1, err := session.Run(
		stmt1.Cypher(),
		stmt1.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
	} else {
		fmt.Println("Results:")
		count := 0
		for result1.Next() && count < 3 {
			record := result1.Record()
			fmt.Printf("  Person: %v\n", record.Values[0])
			count++
		}
		if count == 3 {
			fmt.Println("  ... (more results)")
		}
	}
}

// runExample2 demonstrates property matching
func runExample2(session neo4j.Session) {
	// Example 2: MATCH-WHERE-RETURN with property expressions
	person := cypher.Node("Person").Named("p")

	stmt2, err := cypher.Match(person).
		Where(person.Property("name").Eq("Tom Hanks")).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n2. MATCH-WHERE-RETURN with property comparison:")
	fmt.Println(stmt2.Cypher())
	fmt.Println("Parameters:", stmt2.Params())

	// Execute the query
	result2, err := session.Run(
		stmt2.Cypher(),
		stmt2.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
	} else {
		fmt.Println("Results:")
		for result2.Next() {
			record := result2.Record()
			fmt.Printf("  Person: %v\n", record.Values[0])
		}
	}
}

// runExample3 demonstrates relationship patterns
func runExample3(session neo4j.Session) {
	// Example 3: Relationship patterns
	movie := cypher.Node("Movie").Named("m")
	person := cypher.Node("Person").Named("p")

	// Create a relationship between the nodes
	actedIn := person.RelationshipTo(movie, "ACTED_IN").Named("r")

	// Create a pattern with the relationship
	actedInPattern := cypher.Pattern(person, actedIn, movie)

	stmt3, err := cypher.Match(actedInPattern).
		Returning(person, movie).
		Limit(3).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n3. MATCH with relationship pattern:")
	fmt.Println(stmt3.Cypher())
	fmt.Println("Parameters:", stmt3.Params())

	// Execute the query
	result3, err := session.Run(
		stmt3.Cypher(),
		stmt3.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
	} else {
		fmt.Println("Results:")
		for result3.Next() {
			record := result3.Record()
			fmt.Printf("  Person: %v, Movie: %v\n", record.Values[0], record.Values[1])
		}
	}
}

// runExample4 demonstrates WITH clause and multiple conditions
func runExample4(session neo4j.Session) {
	// Example 4: Using WITH clause and multiple conditions
	// Create nodes and conditions
	tomHanks := cypher.Node("Person").Named("p")
	tomHandsProp := tomHanks.Property("name").Eq("Tom Hanks")

	movieNode := cypher.Node("Movie").Named("m")

	// Create relationship between nodes
	actedIn := tomHanks.RelationshipTo(movieNode, "ACTED_IN")

	// Create a pattern with the relationship
	actedInPattern := cypher.Pattern(tomHanks, actedIn, movieNode)

	stmt4, err := cypher.Match(tomHanks).
		Where(tomHandsProp).
		With(tomHanks).
		Match(actedInPattern).
		Where(movieNode.Property("released").Gt(cypher.ParamWithValue("year", 2000))).
		Returning(
			cypher.As(movieNode.Property("title"), "title"),
			cypher.As(movieNode.Property("released"), "year"),
		).
		OrderBy(cypher.Desc(movieNode.Property("released"))).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n4. Query with WITH clause and parameters:")
	fmt.Println(stmt4.Cypher())
	fmt.Println("Parameters:", stmt4.Params())

	// Execute the query
	result4, err := session.Run(
		stmt4.Cypher(),
		stmt4.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
	} else {
		fmt.Println("Results:")
		for result4.Next() {
			record := result4.Record()
			fmt.Printf("  Movie: %v (%v)\n", record.Values[0], record.Values[1])
		}
	}
}

// runExample5 demonstrates complex logical expressions
func runExample5(session neo4j.Session) {
	// Example 5: Complex logical expressions
	movie := cypher.Node("Movie").Named("m")

	condition := movie.Property("released").Gt(2000).
		And(movie.Property("released").Lt(2010)).
		Or(movie.Property("title").Contains("Matrix").Not())

	stmt5, err := cypher.Match(movie).
		Where(condition).
		Returning(movie.Property("title"), movie.Property("released")).
		OrderBy(movie.Property("released")).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n5. Query with complex logical expressions:")
	fmt.Println(stmt5.Cypher())
	fmt.Println("Parameters:", stmt5.Params())

	// Execute the query
	result5, err := session.Run(
		stmt5.Cypher(),
		stmt5.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
	} else {
		fmt.Println("Results:")
		for result5.Next() {
			record := result5.Record()
			fmt.Printf("  Movie: %v (%v)\n", record.Values[0], record.Values[1])
		}
	}
}
