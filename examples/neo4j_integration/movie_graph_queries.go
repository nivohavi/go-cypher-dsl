package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/driver"
)

func main() {
	fmt.Println("Go Cypher DSL - Neo4j Integration Examples")
	fmt.Println("=========================================")

	// Neo4j connection parameters
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "password" // Replace with your actual password

	// Create driver
	neo4jDriver, err := neo4j.NewDriver(uri,
		neo4j.BasicAuth(username, password, ""))
	if err != nil {
		panic(err)
	}
	defer neo4jDriver.Close()

	// Test the connection
	err = neo4jDriver.VerifyConnectivity()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Neo4j: %v", err))
	}
	fmt.Println("Successfully connected to Neo4j")

	// Using the traditional Neo4j session approach
	traditionalSessionExample(neo4jDriver)

	// Using the improved SessionManager approach
	sessionManagerExample(neo4jDriver)
}

// Example of using the traditional Neo4j driver approach
func traditionalSessionExample(driver neo4j.Driver) {
	fmt.Println("\nTraditional Neo4j Session Examples")
	fmt.Println("=================================")

	// Create a session
	session := driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close()

	// Example 1: Find all Person nodes
	findAllPeople(session)

	// Example 2: Find specific person by name
	findPersonByName(session, "Tom Hanks")

	// Example 3: Find person-movie relationships
	findPersonMovieRelationships(session)

	// Example 4: Find movies released in a date range
	findMoviesInDateRange(session, 1990, 2000)

	// Example 5: Find actors and directors
	findActorsAndDirectors(session)
}

// Example of using the SessionManager helper
func sessionManagerExample(neo4jDriver neo4j.Driver) {
	fmt.Println("\nSessionManager Examples")
	fmt.Println("=====================")

	// Create the session manager and query helper
	sessionManager := driver.NewSessionManager(neo4jDriver)
	queryHelper := driver.NewQueryHelper()
	ctx := context.Background()

	// Example 1: Find all Person nodes using CollectList
	findAllPeopleWithManager(sessionManager, queryHelper, ctx)

	// Example 2: Find specific person by name using CollectSingle
	findPersonByNameWithManager(sessionManager, queryHelper, ctx, "Tom Hanks")

	// Example 3: Count movies
	countMoviesWithManager(sessionManager, queryHelper, ctx)

	// Example 4: Collect results as a map
	collectMovieMapWithManager(sessionManager, queryHelper, ctx)
}

// Traditional approach examples

func findAllPeople(session neo4j.Session) {
	// Create a node pattern
	person := cypher.Node("Person").Named("p")

	// Build a query with MATCH and RETURN
	stmt, err := cypher.Match(person).
		Returning(person).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n1. Find All People (limited to 5):")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query
	result, err := session.Run(
		stmt.Cypher(),
		stmt.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Process results
	fmt.Println("Results:")
	for result.Next() {
		record := result.Record()
		fmt.Printf("  Person: %v\n", record.Values[0])
	}
}

func findPersonByName(session neo4j.Session, name string) {
	// Create a node pattern
	person := cypher.Node("Person").Named("p")

	// Build a query with WHERE condition
	stmt, err := cypher.Match(person).
		Where(person.Property("name").Eq(name)).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Printf("\n2. Find Person by Name (%s):\n", name)
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query
	result, err := session.Run(
		stmt.Cypher(),
		stmt.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Process results
	fmt.Println("Results:")
	for result.Next() {
		record := result.Record()
		fmt.Printf("  Person: %v\n", record.Values[0])
	}
}

func findPersonMovieRelationships(session neo4j.Session) {
	// Create node patterns
	person := cypher.Node("Person").Named("p")
	movie := cypher.Node("Movie").Named("m")

	// Create a relationship from person to movie
	acted := person.RelationshipTo(movie, "ACTED_IN").Named("r")

	// Create a pattern using the relationship
	pattern := cypher.Path(person, acted, movie)

	// Build a query with the pattern
	stmt, err := cypher.Match(pattern).
		Returning(
			person.Property("name"),
			movie.Property("title"),
		).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n3. Find Person-Movie Relationships:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query
	result, err := session.Run(
		stmt.Cypher(),
		stmt.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Process results
	fmt.Println("Results:")
	for result.Next() {
		record := result.Record()
		fmt.Printf("  Actor: %v, Movie: %v\n", record.Values[0], record.Values[1])
	}
}

func findMoviesInDateRange(session neo4j.Session, startYear, endYear int) {
	// Create a node pattern
	movie := cypher.Node("Movie").Named("m")

	// Create condition for date range
	dateRangeCondition := movie.Property("released").Gte(startYear).
		And(movie.Property("released").Lte(endYear))

	// Build a query with the condition
	stmt, err := cypher.Match(movie).
		Where(dateRangeCondition).
		Returning(
			movie.Property("title"),
			movie.Property("released"),
		).
		OrderBy(movie.Property("released")).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Printf("\n4. Find Movies Released Between %d and %d:\n", startYear, endYear)
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query
	result, err := session.Run(
		stmt.Cypher(),
		stmt.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Process results
	fmt.Println("Results:")
	for result.Next() {
		record := result.Record()
		fmt.Printf("  Movie: %v (%v)\n", record.Values[0], record.Values[1])
	}
}

func findActorsAndDirectors(session neo4j.Session) {
	// Create node patterns
	actor := cypher.Node("Person").Named("actor")
	movie := cypher.Node("Movie").Named("movie")
	director := cypher.Node("Person").Named("director")

	// Create relationships
	actedIn := actor.RelationshipTo(movie, "ACTED_IN")
	directedBy := movie.RelationshipTo(director, "DIRECTED_BY")

	// Build a query with multiple patterns
	stmt, err := cypher.Match(cypher.Path(actor, actedIn, movie, directedBy, director)).
		Returning(
			cypher.As(actor.Property("name"), "actor"),
			cypher.As(movie.Property("title"), "movie"),
			cypher.As(director.Property("name"), "director"),
		).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n5. Find Actors and Their Directors:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query
	result, err := session.Run(
		stmt.Cypher(),
		stmt.Params(),
	)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Process results
	fmt.Println("Results:")
	for result.Next() {
		record := result.Record()
		actorName, _ := record.Get("actor")
		movieTitle, _ := record.Get("movie")
		directorName, _ := record.Get("director")
		fmt.Printf("  Actor: %v, Movie: %v, Director: %v\n", actorName, movieTitle, directorName)
	}
}

// SessionManager approach examples

func findAllPeopleWithManager(sessionManager *driver.SessionManager, queryHelper *driver.QueryHelper, ctx context.Context) {
	// Create a node pattern
	person := cypher.Node("Person").Named("p")

	// Build a query with MATCH and RETURN
	stmt, err := cypher.Match(person).
		Returning(person).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n1. Find All People Using SessionManager:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query with CollectList handler
	result, err := sessionManager.ExecuteRead(ctx, stmt, queryHelper.CollectList("p"))
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Process results
	peopleList, ok := result.([]any)
	if ok && len(peopleList) > 0 {
		fmt.Printf("Results: Found %d people\n", len(peopleList))
		for i, person := range peopleList {
			fmt.Printf("  Person %d: %v\n", i+1, person)
		}
	} else {
		fmt.Println("No people found")
	}
}

func findPersonByNameWithManager(sessionManager *driver.SessionManager, queryHelper *driver.QueryHelper, ctx context.Context, name string) {
	// Create a node pattern
	person := cypher.Node("Person").Named("p")

	// Build a query with WHERE condition
	stmt, err := cypher.Match(person).
		Where(person.Property("name").Eq(name)).
		Returning(person).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Printf("\n2. Find Person by Name (%s) Using SessionManager:\n", name)
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query with CollectSingle handler
	result, err := sessionManager.ExecuteRead(ctx, stmt, queryHelper.CollectSingle("p"))
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Process result
	if result != nil {
		fmt.Printf("Result: Found person: %v\n", result)
	} else {
		fmt.Printf("No person found with name '%s'\n", name)
	}
}

func countMoviesWithManager(sessionManager *driver.SessionManager, queryHelper *driver.QueryHelper, ctx context.Context) {
	// Create a node pattern
	movie := cypher.Node("Movie").Named("m")

	// Build a query
	stmt, err := cypher.Match(movie).
		Returning(movie).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n3. Count Movies Using SessionManager:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query with CountResults handler
	result, err := sessionManager.ExecuteRead(ctx, stmt, queryHelper.CountResults())
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Display count
	fmt.Printf("Result: Found %v movies\n", result)
}

func collectMovieMapWithManager(sessionManager *driver.SessionManager, queryHelper *driver.QueryHelper, ctx context.Context) {
	// Create a node pattern
	movie := cypher.Node("Movie").Named("m")

	// Build a query
	stmt, err := cypher.Match(movie).
		Returning(
			cypher.As(movie.Property("title"), "title"),
			cypher.As(movie.Property("released"), "year"),
		).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("\n4. Collect Movies as Map Using SessionManager:")
	fmt.Println(stmt.Cypher())
	fmt.Println("Parameters:", stmt.Params())

	// Execute the query with CollectMap handler
	result, err := sessionManager.ExecuteRead(ctx, stmt, queryHelper.CollectMap("title", "year"))
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Display map
	movieMap, ok := result.(map[string]any)
	if ok && len(movieMap) > 0 {
		fmt.Printf("Results: Found %d movies\n", len(movieMap))
		for title, year := range movieMap {
			fmt.Printf("  %s: %v\n", title, year)
		}
	} else {
		fmt.Println("No movies found or could not convert result to map")
	}
}
