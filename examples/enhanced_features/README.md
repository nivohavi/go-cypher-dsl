# Enhanced Features Examples

This directory contains examples demonstrating the newest and most advanced features of the Go Cypher DSL library. These examples showcase features that make building and executing Neo4j Cypher queries even more convenient and readable.

## Features Demonstrated

### 1. ComplexPath for Multi-Level Relationships

The `ComplexPath` helper function simplifies the creation of complex path patterns with multiple relationships. Instead of nesting relationship calls, you can define the entire path in a single linear function call:

```go
// Using ComplexPath
path := cypher.ComplexPath(
    user,
    "WORKS_AT", company,
    "LOCATED_IN", city,
    "PART_OF", country
)
```

This produces the same result as the more verbose traditional approach:

```go
// Traditional approach
path := cypher.Path(
    user, 
    user.RelationshipTo(company, "WORKS_AT"),
    company,
    company.RelationshipTo(city, "LOCATED_IN"),
    city,
    city.RelationshipTo(country, "PART_OF"),
    country
)
```

### 2. Property Comparison Helpers

The new property comparison helpers reduce the verbosity of common property comparisons:

```go
// Old way
cypher.Gt(
    cypher.Property("u", "age"),
    cypher.Param(30)
)

// New way
cypher.CompareProperty("u", "age", ">", 30)

// With named parameters
cypher.NamedCompareProperty("u", "age", ">", "minAge", 30)
```

These helpers support all common comparison operators: `=`, `!=`, `>`, `>=`, `<`, `<=`, `IN`, `CONTAINS`, `STARTS WITH`, and `ENDS WITH`.

### 3. Schema Management Helpers

The new schema package provides helpers for common database schema operations:

```go
// Create constraints
uniqueConstraint, _ := schema.CreateUniqueConstraint("user_email_unique", "User", "email")
nodeKeyConstraint, _ := schema.CreateNodeKeyConstraint("user_id_key", "User", "id")
existsConstraint, _ := schema.CreateExistsConstraint("user_name_exists", "User", "name")

// Create indexes
index, _ := schema.CreateIndex("user_name_idx", "User", "firstName", "lastName")
fullTextIndex, _ := schema.CreateFullTextIndex("content_search", 
    []string{"Post", "Comment"}, 
    []string{"title", "content"})

// Drop constraints and indexes
dropConstraint, _ := schema.DropConstraint("user_email_unique")
dropIndex, _ := schema.DropIndex("user_name_idx")

// Show all constraints and indexes
showConstraints, _ := schema.ShowConstraints()
showIndexes, _ := schema.ShowIndexes()
```

### 4. Neo4j Driver Integration

The new driver package provides simplified session management and result handling:

```go
// Create helpers
sessionManager := driver.NewSessionManager(neo4jDriver)
queryHelper := driver.NewQueryHelper()

// Execute read query with various result handlers
result, err := sessionManager.ExecuteRead(ctx, query, queryHelper.CollectSingle("user"))
listResult, err := sessionManager.ExecuteRead(ctx, query, queryHelper.CollectList("user"))
mapResult, err := sessionManager.ExecuteRead(ctx, query, queryHelper.CollectMap("userId", "userName"))
allResult, err := sessionManager.ExecuteRead(ctx, query, queryHelper.CollectAll())
count, err := sessionManager.ExecuteRead(ctx, query, queryHelper.CountResults())

// Execute write query
writeResult, err := sessionManager.ExecuteWrite(ctx, query, customHandler)

// Execute multiple write operations in a single transaction
batchResult, err := sessionManager.ExecuteBatchWrite(ctx, statements, batchHandler)
```

## Running the Example

To run this example:

```bash
go run example.go
```

Note: The driver integration part requires a running Neo4j instance. The example includes code that attempts to connect to Neo4j, but you may need to modify the connection details to match your setup. 