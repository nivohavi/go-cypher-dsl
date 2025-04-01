package main

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/schema"
)

func main() {
	fmt.Println("Go Cypher DSL - Schema Management Examples")
	fmt.Println("=========================================")

	// Demonstrate schema management statements generation
	demonstrateConstraints()
	demonstrateIndexes()
	demonstrateFullTextIndexes()
	demonstrateDropOperations()

	// Demonstrate executing schema statements against Neo4j (commented out by default)
	// To run these examples, uncomment them and provide your Neo4j credentials
	// demonstrateSchemaOperationsOnNeo4j()
}

func demonstrateConstraints() {
	fmt.Println("\n1. Constraint Creation")
	fmt.Println("--------------------")

	// Example 1: Create a unique constraint for User's email
	userEmailUnique, err := schema.CreateUniqueConstraint("user_email_unique", "User", "email")
	if err != nil {
		fmt.Printf("Error creating uniqueness constraint: %v\n", err)
		return
	}

	fmt.Println("Uniqueness Constraint:")
	fmt.Println(userEmailUnique.Cypher())
	fmt.Println("Parameters:", userEmailUnique.Params())

	// Example 2: Create a node existence constraint (property must exist)
	personNameExists, err := schema.CreateExistsConstraint("person_name_exists", "Person", "name")
	if err != nil {
		fmt.Printf("Error creating node property existence constraint: %v\n", err)
		return
	}

	fmt.Println("\nNode Property Existence Constraint:")
	fmt.Println(personNameExists.Cypher())
	fmt.Println("Parameters:", personNameExists.Params())

	// Example 3: Create a relationship property existence constraint
	reviewRatingExists, err := schema.CreateRelationshipConstraint("review_rating_exists", "REVIEWED", "rating")
	if err != nil {
		fmt.Printf("Error creating relationship property existence constraint: %v\n", err)
		return
	}

	fmt.Println("\nRelationship Property Existence Constraint:")
	fmt.Println(reviewRatingExists.Cypher())
	fmt.Println("Parameters:", reviewRatingExists.Params())

	// Example 4: Create a node key constraint (composite unique constraint)
	personCompositeKey, err := schema.CreateNodeKeyConstraint("person_composite_key", "Person", "firstName", "lastName", "dob")
	if err != nil {
		fmt.Printf("Error creating node key constraint: %v\n", err)
		return
	}

	fmt.Println("\nNode Key Constraint (Composite):")
	fmt.Println(personCompositeKey.Cypher())
	fmt.Println("Parameters:", personCompositeKey.Params())
}

func demonstrateIndexes() {
	fmt.Println("\n2. Index Creation")
	fmt.Println("---------------")

	// Example 1: Create a single-property index
	userNameIndex, err := schema.CreateIndex("user_name_idx", "User", "name")
	if err != nil {
		fmt.Printf("Error creating index: %v\n", err)
		return
	}

	fmt.Println("Single Property Index:")
	fmt.Println(userNameIndex.Cypher())
	fmt.Println("Parameters:", userNameIndex.Params())

	// Example 2: Create a composite index (multiple properties)
	movieCompositeIndex, err := schema.CreateIndex("movie_composite_idx", "Movie", "title", "released")
	if err != nil {
		fmt.Printf("Error creating composite index: %v\n", err)
		return
	}

	fmt.Println("\nComposite Index:")
	fmt.Println(movieCompositeIndex.Cypher())
	fmt.Println("Parameters:", movieCompositeIndex.Params())
}

func demonstrateFullTextIndexes() {
	fmt.Println("\n3. Full-Text Index Creation")
	fmt.Println("-------------------------")

	// Example 1: Create a full-text index on post content
	contentSearch, err := schema.CreateFullTextIndex(
		"content_search", []string{"Post", "Comment"}, []string{"title", "content"})
	if err != nil {
		fmt.Printf("Error creating full-text index: %v\n", err)
		return
	}

	fmt.Println("Full-Text Index:")
	fmt.Println(contentSearch.Cypher())
	fmt.Println("Parameters:", contentSearch.Params())
}

func demonstrateDropOperations() {
	fmt.Println("\n4. Drop Operations")
	fmt.Println("----------------")

	// Example 1: Drop a constraint
	dropConstraint, err := schema.DropConstraint("user_email_unique")
	if err != nil {
		fmt.Printf("Error creating drop constraint statement: %v\n", err)
		return
	}

	fmt.Println("Drop Constraint:")
	fmt.Println(dropConstraint.Cypher())
	fmt.Println("Parameters:", dropConstraint.Params())

	// Example 2: Drop an index
	dropIndex, err := schema.DropIndex("user_name_idx")
	if err != nil {
		fmt.Printf("Error creating drop index statement: %v\n", err)
		return
	}

	fmt.Println("\nDrop Index:")
	fmt.Println(dropIndex.Cypher())
	fmt.Println("Parameters:", dropIndex.Params())

	// Example 3: Show all constraints
	showConstraints, err := schema.ShowConstraints()
	if err != nil {
		fmt.Printf("Error creating show constraints statement: %v\n", err)
		return
	}

	fmt.Println("\nShow All Constraints:")
	fmt.Println(showConstraints.Cypher())
	fmt.Println("Parameters:", showConstraints.Params())

	// Example 4: Show all indexes
	showIndexes, err := schema.ShowIndexes()
	if err != nil {
		fmt.Printf("Error creating show indexes statement: %v\n", err)
		return
	}

	fmt.Println("\nShow All Indexes:")
	fmt.Println(showIndexes.Cypher())
	fmt.Println("Parameters:", showIndexes.Params())
}

// This function demonstrates executing schema operations on an actual Neo4j database
// It's commented out by default to avoid making changes to user's database unintentionally
func demonstrateSchemaOperationsOnNeo4j() {
	fmt.Println("\n5. Executing Schema Operations on Neo4j")
	fmt.Println("--------------------------------------")

	// Neo4j connection parameters
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "password" // Replace with your actual password

	// Create driver
	driver, err := neo4j.NewDriver(uri,
		neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Printf("Error creating Neo4j driver: %v\n", err)
		return
	}
	defer driver.Close()

	// Test the connection
	err = driver.VerifyConnectivity()
	if err != nil {
		fmt.Printf("Failed to connect to Neo4j: %v\n", err)
		return
	}
	fmt.Println("Successfully connected to Neo4j")

	// Create a session for schema operations
	session := driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close()

	// Create a unique constraint
	userEmailUnique, err := schema.CreateUniqueConstraint("user_email_unique", "User", "email")
	if err != nil {
		fmt.Printf("Error creating uniqueness constraint: %v\n", err)
		return
	}

	fmt.Println("\nExecuting Uniqueness Constraint Creation:")
	fmt.Println(userEmailUnique.Cypher())

	// Execute the constraint creation
	_, err = session.Run(userEmailUnique.Cypher(), userEmailUnique.Params())
	if err != nil {
		fmt.Printf("Error executing constraint creation: %v\n", err)
		return
	}
	fmt.Println("Successfully created constraint")

	// Create an index
	userNameIndex, err := schema.CreateIndex("user_name_idx", "User", "name")
	if err != nil {
		fmt.Printf("Error creating index: %v\n", err)
		return
	}

	fmt.Println("\nExecuting Index Creation:")
	fmt.Println(userNameIndex.Cypher())

	// Execute the index creation
	_, err = session.Run(userNameIndex.Cypher(), userNameIndex.Params())
	if err != nil {
		fmt.Printf("Error executing index creation: %v\n", err)
		return
	}
	fmt.Println("Successfully created index")

	// Show all constraints
	showConstraints, err := schema.ShowConstraints()
	if err != nil {
		fmt.Printf("Error creating show constraints statement: %v\n", err)
		return
	}

	fmt.Println("\nListing All Constraints:")
	result, err := session.Run(showConstraints.Cypher(), showConstraints.Params())
	if err != nil {
		fmt.Printf("Error executing show constraints: %v\n", err)
		return
	}

	// Process results
	fmt.Println("Existing Constraints:")
	for result.Next() {
		record := result.Record()
		fmt.Printf("  Constraint: %v\n", record.Values)
	}

	// Show all indexes
	showIndexes, err := schema.ShowIndexes()
	if err != nil {
		fmt.Printf("Error creating show indexes statement: %v\n", err)
		return
	}

	fmt.Println("\nListing All Indexes:")
	result, err = session.Run(showIndexes.Cypher(), showIndexes.Params())
	if err != nil {
		fmt.Printf("Error executing show indexes: %v\n", err)
		return
	}

	// Process results
	fmt.Println("Existing Indexes:")
	for result.Next() {
		record := result.Record()
		fmt.Printf("  Index: %v\n", record.Values)
	}
}
