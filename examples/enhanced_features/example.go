package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/driver"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/schema"
)

func main() {
	fmt.Println("Go Cypher DSL - Enhanced Features Demo")
	fmt.Println("======================================")

	// Connection details (replace with your actual Neo4j credentials)
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "password"

	// Demo the new features - these don't require a database connection
	demonstrateComplexPaths()
	demonstrateComparisonHelpers()
	demonstrateSchemaHelpers()

	// Connect to Neo4j to demonstrate driver helpers
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Printf("Error connecting to Neo4j: %v\n", err)
		return
	}
	defer driver.Close()

	// Test connection
	err = driver.VerifyConnectivity()
	if err != nil {
		fmt.Printf("Failed to connect to Neo4j: %v\n", err)
		return
	}

	demonstrateDriverHelpers(driver)
}

// Feature 1: ComplexPath for multi-level relationships
func demonstrateComplexPaths() {
	fmt.Println("\n1. ComplexPath for Multi-Level Relationships")
	fmt.Println("-------------------------------------------")

	// Create nodes for our path
	user := cypher.Node("User").Named("u")
	company := cypher.Node("Company").Named("c")
	city := cypher.Node("City").Named("city")
	country := cypher.Node("Country").Named("country")

	// Create a complex path using the traditional approach (nesting relationships)
	path1, err := cypher.Match(
		cypher.Path(
			user,
			user.RelationshipTo(company, "WORKS_AT"),
			company,
			company.RelationshipTo(city, "LOCATED_IN"),
			city,
			city.RelationshipTo(country, "PART_OF"),
			country,
		),
	).
		Where(cypher.Eq(
			cypher.Property("u", "name"),
			cypher.Param("John"),
		)).
		Returning(
			cypher.Property("u", "name"),
			cypher.Property("c", "name"),
			cypher.Property("country", "name"),
		).
		Build()

	if err != nil {
		fmt.Printf("Error building query 1: %v\n", err)
		return
	}

	fmt.Println("Traditional approach:")
	fmt.Println(cypher.PrettyPrint(path1))

	// Create the same path using the new ComplexPath helper
	path2, err := cypher.Match(
		cypher.ComplexPath(
			user,
			"WORKS_AT", company,
			"LOCATED_IN", city,
			"PART_OF", country,
		),
	).
		Where(cypher.Eq(
			cypher.Property("u", "name"),
			cypher.Param("John"),
		)).
		Returning(
			cypher.Property("u", "name"),
			cypher.Property("c", "name"),
			cypher.Property("country", "name"),
		).
		Build()

	if err != nil {
		fmt.Printf("Error building query 2: %v\n", err)
		return
	}

	fmt.Println("\nUsing new ComplexPath helper:")
	fmt.Println(cypher.PrettyPrint(path2))
}

// Feature 2: Property comparison helpers
func demonstrateComparisonHelpers() {
	fmt.Println("\n2. Property Comparison Helpers")
	fmt.Println("----------------------------")

	// Traditional approach for comparisons
	userNode := cypher.Node("User").Named("u")
	query1, err := cypher.Match(userNode).
		Where(
			cypher.And(
				cypher.Gt(
					cypher.Property("u", "age"),
					cypher.Param(30),
				),
				cypher.Eq(
					cypher.Property("u", "active"),
					cypher.Param(true),
				),
			),
		).
		Returning(cypher.Property("u", "name")).
		Build()

	if err != nil {
		fmt.Printf("Error building query 1: %v\n", err)
		return
	}

	fmt.Println("Traditional approach for comparisons:")
	fmt.Println(cypher.PrettyPrint(query1))
	fmt.Println("Parameters:", query1.Params())

	// Using the new CompareProperty helper
	query2, err := cypher.Match(userNode).
		Where(
			cypher.And(
				cypher.CompareProperty("u", "age", ">", 30),
				cypher.CompareProperty("u", "active", "=", true),
			),
		).
		Returning(cypher.Property("u", "name")).
		Build()

	if err != nil {
		fmt.Printf("Error building query 2: %v\n", err)
		return
	}

	fmt.Println("\nUsing new CompareProperty helper:")
	fmt.Println(cypher.PrettyPrint(query2))
	fmt.Println("Parameters:", query2.Params())

	// Using the NamedCompareProperty helper
	query3, err := cypher.Match(userNode).
		Where(
			cypher.And(
				cypher.NamedCompareProperty("u", "age", ">", "minAge", 30),
				cypher.NamedCompareProperty("u", "active", "=", "isActive", true),
			),
		).
		Returning(cypher.Property("u", "name")).
		Build()

	if err != nil {
		fmt.Printf("Error building query 3: %v\n", err)
		return
	}

	fmt.Println("\nUsing new NamedCompareProperty helper:")
	fmt.Println(cypher.PrettyPrint(query3))
	fmt.Println("Parameters:", query3.Params())

	// String-specific operations
	query4, err := cypher.Match(userNode).
		Where(
			cypher.Or(
				cypher.CompareProperty("u", "name", "STARTS WITH", "Jo"),
				cypher.CompareProperty("u", "email", "CONTAINS", "@example.com"),
			),
		).
		Returning(cypher.Property("u", "name")).
		Build()

	if err != nil {
		fmt.Printf("Error building query 4: %v\n", err)
		return
	}

	fmt.Println("\nString operations with CompareProperty:")
	fmt.Println(cypher.PrettyPrint(query4))
	fmt.Println("Parameters:", query4.Params())
}

// Feature 3: Schema helpers
func demonstrateSchemaHelpers() {
	fmt.Println("\n3. Schema Operation Helpers")
	fmt.Println("--------------------------")

	// Create a unique constraint
	uniqueConstraint, _ := schema.CreateUniqueConstraint("user_email_unique", "User", "email")
	fmt.Println("Create unique constraint:")
	fmt.Println(uniqueConstraint.Cypher())

	// Create a node key constraint
	nodeKeyConstraint, _ := schema.CreateNodeKeyConstraint("user_id_key", "User", "id")
	fmt.Println("\nCreate node key constraint:")
	fmt.Println(nodeKeyConstraint.Cypher())

	// Create an index
	index, _ := schema.CreateIndex("user_name_idx", "User", "firstName", "lastName")
	fmt.Println("\nCreate index:")
	fmt.Println(index.Cypher())

	// Create a full-text index
	fullTextIndex, _ := schema.CreateFullTextIndex("content_search",
		[]string{"Post", "Comment"},
		[]string{"title", "content"})
	fmt.Println("\nCreate full-text index:")
	fmt.Println(fullTextIndex.Cypher())

	// Show all constraints
	showConstraints, _ := schema.ShowConstraints()
	fmt.Println("\nShow constraints:")
	fmt.Println(showConstraints.Cypher())
}

// Feature 4: Neo4j driver helpers
func demonstrateDriverHelpers(neo4jDriver neo4j.Driver) {
	fmt.Println("\n4. Neo4j Driver Integration Helpers")
	fmt.Println("---------------------------------")

	// Create the session manager
	sessionManager := driver.NewSessionManager(neo4jDriver)
	queryHelper := driver.NewQueryHelper()
	ctx := context.Background()

	// Create a sample query
	userNode := cypher.Node("User").Named("u")
	query, err := cypher.Match(userNode).
		Returning(userNode).
		Limit(5).
		Build()

	if err != nil {
		fmt.Printf("Error building query: %v\n", err)
		return
	}

	fmt.Println("Query to execute:")
	fmt.Println(cypher.PrettyPrint(query))

	// Execute read query with handler for single result
	result, err := sessionManager.ExecuteRead(ctx, query, queryHelper.CollectSingle("u"))
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	fmt.Println("\nResult of CollectSingle:")
	if result != nil {
		fmt.Printf("Found user: %v\n", result)
	} else {
		fmt.Println("No users found")
	}

	// Execute read query with handler for list of results
	listResult, err := sessionManager.ExecuteRead(ctx, query, queryHelper.CollectList("u"))
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	fmt.Println("\nResult of CollectList:")
	userList, ok := listResult.([]any)
	if ok && len(userList) > 0 {
		fmt.Printf("Found %d users\n", len(userList))
		for i, user := range userList {
			if i < 3 { // Show only first 3 to avoid too much output
				fmt.Printf("- User %d: %v\n", i+1, user)
			}
		}
		if len(userList) > 3 {
			fmt.Printf("...and %d more\n", len(userList)-3)
		}
	} else {
		fmt.Println("No users found")
	}

	// Execute read query with handler for counting results
	countResult, err := sessionManager.ExecuteRead(ctx, query, queryHelper.CountResults())
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	fmt.Println("\nResult of CountResults:")
	fmt.Printf("Total users: %v\n", countResult)

	// Example of batch operations
	fmt.Println("\nBatch Operation Example:")
	fmt.Println("(Simulated - not actually executed to avoid data modifications)")

	// Create statements for batch operation
	createUser1, _ := cypher.Create(
		cypher.NodeWithProperties(
			map[string]any{"name": "Alice", "age": 30},
			"User",
		),
	).Build()

	createUser2, _ := cypher.Create(
		cypher.NodeWithProperties(
			map[string]any{"name": "Bob", "age": 28},
			"User",
		),
	).Build()

	fmt.Println("Statement 1: " + createUser1.Cypher())
	fmt.Println("Statement 2: " + createUser2.Cypher())
	fmt.Println("These would be executed in a single transaction using ExecuteBatchWrite")
}
