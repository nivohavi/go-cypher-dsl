# go-cypher-dsl

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/nivohavi/go-cypher-dsl.svg)](https://pkg.go.dev/github.com/nivohavi/go-cypher-dsl)
[![Go Report Card](https://goreportcard.com/badge/github.com/nivohavi/go-cypher-dsl)](https://goreportcard.com/report/github.com/nivohavi/go-cypher-dsl)

A fluent, type-safe Cypher query builder for Go. Build Neo4j queries programmatically with syntax completion, reducing errors and improving development speed.

## Installation

```bash
go get github.com/nivohavi/go-cypher-dsl
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/nivohavi/go-cypher-dsl/pkg/cypher"
)

func main() {
    // Create a node pattern
    person := cypher.Node("Person").Named("p")
    
    // Build a query
    stmt, _ := cypher.Match(person).
        Where(person.Property("name").Eq("John")).
        Returning(person).
        Build()
    
    // Get the query string and parameters
    query := stmt.Cypher()
    params := stmt.Params()
    
    fmt.Println(query)
    // MATCH (p:`Person`) WHERE p.name = 'John' RETURN p
    
    fmt.Println(params)
    // map[]
}
```

## Examples

Check out the [examples](examples/) directory for comprehensive examples covering various usage patterns:

```
examples/
├── basics/              # Core query building examples using improved API
├── enhanced_features/   # Examples showcasing new advanced features
├── neo4j_integration/   # Examples of Neo4j database integration
└── schema_management/   # Examples of schema management (constraints and indexes)
```

To run an example:

```bash
# Run the basic query building example
cd examples/basics
go run query_building.go
```

## Core Concepts

### Nodes and Relationships

```go
// Create nodes
person := cypher.Node("Person").Named("p")
movie := cypher.Node("Movie").Named("m")

// Create a relationship
acted := person.RelationshipTo(movie, "ACTED_IN").Named("r")

// Create a pattern with the relationship
pattern := cypher.Pattern(person, acted, movie)

// Use in a query
query, _ := cypher.Match(pattern).
    Returning(person, movie).
    Build()

fmt.Println(query.Cypher())
// MATCH (p:`Person`)-[r:`ACTED_IN`]->(m:`Movie`) RETURN p, m
```

## Enhanced Features

### Complex Path Construction

Create paths with multiple relationships more easily using `ComplexPath`:

```go
// Define nodes
user := cypher.Node("User").Named("u")
company := cypher.Node("Company").Named("c")
city := cypher.Node("City").Named("city")

// Create a complex path with a single function call
path := cypher.ComplexPath(
    user,
    "WORKS_AT", company,
    "LOCATED_IN", city
)

// Use in a query
stmt, _ := cypher.Match(path).
    Where(cypher.Property("u", "name").Eq("John")).
    Returning(
        cypher.Property("u", "name"),
        cypher.Property("city", "name")
    ).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (u:User)-[:WORKS_AT]->(c:Company)-[:LOCATED_IN]->(city:City) WHERE u.name = $p0 RETURN u.name, city.name
```

### Simplified Property Comparison Helpers

Compare properties more concisely using our comparison helpers:

```go
// Traditional approach
query1, _ := cypher.Match(userNode).
    Where(
        cypher.And(
            cypher.Gt(
                cypher.Property("u", "age"),
                cypher.Param(30)
            ),
            cypher.Eq(
                cypher.Property("u", "active"),
                cypher.Param(true)
            )
        )
    ).
    Returning(cypher.Property("u", "name")).
    Build()

// Simplified comparison using CompareProperty
query2, _ := cypher.Match(userNode).
    Where(
        cypher.And(
            cypher.CompareProperty("u", "age", ">", 30),
            cypher.CompareProperty("u", "active", "=", true)
        )
    ).
    Returning(cypher.Property("u", "name")).
    Build()

// Named parameters with NamedCompareProperty
query3, _ := cypher.Match(userNode).
    Where(
        cypher.NamedCompareProperty("u", "age", ">", "minAge", 30)
    ).
    Returning(cypher.Property("u", "name")).
    Build()
```

### Schema Management Helpers

Create and manage database constraints and indexes:

```go
// Create a unique constraint
uniqueConstraint, _ := schema.CreateUniqueConstraint("user_email_unique", "User", "email")
fmt.Println(uniqueConstraint.Cypher())
// CREATE CONSTRAINT user_email_unique IF NOT EXISTS FOR (n:User) REQUIRE n.email IS UNIQUE

// Create a node key constraint
nodeKeyConstraint, _ := schema.CreateNodeKeyConstraint("user_id_key", "User", "id")
fmt.Println(nodeKeyConstraint.Cypher())
// CREATE CONSTRAINT user_id_key IF NOT EXISTS FOR (n:User) REQUIRE (n.id) IS NODE KEY

// Create an index on multiple properties
index, _ := schema.CreateIndex("user_name_idx", "User", "firstName", "lastName")
fmt.Println(index.Cypher())
// CREATE INDEX user_name_idx IF NOT EXISTS FOR (n:User) ON (n.firstName, n.lastName)

// Create a full-text index
fullTextIndex, _ := schema.CreateFullTextIndex("content_search", 
    []string{"Post", "Comment"}, 
    []string{"title", "content"})
fmt.Println(fullTextIndex.Cypher())
// CALL db.index.fulltext.createNodeIndex('content_search', ['Post', 'Comment'], ['title', 'content'])
```

### Neo4j Driver Integration

Simplified query execution and result handling:

```go
// Create the session manager
sessionManager := driver.NewSessionManager(neo4jDriver)
queryHelper := driver.NewQueryHelper()
ctx := context.Background()

// Create a query to execute
query, _ := cypher.Match(userNode).
    Returning(userNode).
    Limit(5).
    Build()

// Execute a read query and collect a single result
result, err := sessionManager.ExecuteRead(ctx, query, 
    queryHelper.CollectSingle("u"))

// Execute a read query and collect multiple results as a list
listResult, err := sessionManager.ExecuteRead(ctx, query, 
    queryHelper.CollectList("u"))

// Execute a read query and just count the results
countResult, err := sessionManager.ExecuteRead(ctx, query, 
    queryHelper.CountResults())

// Execute a batch of write operations in a single transaction
statements := []core.Statement{statement1, statement2}
batchResult, err := sessionManager.ExecuteBatchWrite(ctx, statements, 
    func(results []neo4j.Result) (any, error) {
        // Process results
        return nil, nil
    })
```

### Properties and Conditions

```go
person := cypher.Node("Person").Named("p")

// Property access
nameProperty := person.Property("name")

// Create conditions
condition := nameProperty.Eq("Tom Hanks").
    And(person.Property("born").Gt(1950))

// Build query with condition
query, _ := cypher.Match(person).
    Where(condition).
    Returning(person).
    Build()

fmt.Println(query.Cypher())
// MATCH (p:`Person`) WHERE p.name = 'Tom Hanks' AND p.born > 1950 RETURN p
```

## Common Query Patterns

### Basic MATCH-RETURN

```go
person := cypher.Node("Person").Named("p")

stmt, _ := cypher.Match(person).
    Returning(person).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (p:`Person`) RETURN p
```

### Filtering with WHERE

```go
person := cypher.Node("Person").Named("p")

stmt, _ := cypher.Match(person).
    Where(person.Property("name").Eq("Tom Hanks").
        Or(person.Property("name").Eq("Tom Cruise"))).
    Returning(person).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (p:`Person`) WHERE p.name = 'Tom Hanks' OR p.name = 'Tom Cruise' RETURN p
```

### Relationship Patterns

```go
person := cypher.Node("Person").Named("p")
movie := cypher.Node("Movie").Named("m")
acted := person.RelationshipTo(movie, "ACTED_IN").Named("r")

stmt, _ := cypher.Match(cypher.Pattern(person, acted, movie)).
    Where(movie.Property("title").Eq("The Matrix")).
    Returning(person.Property("name")).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (p:`Person`)-[r:`ACTED_IN`]->(m:`Movie`) WHERE m.title = 'The Matrix' RETURN p.name
```

### Creating Data

```go
person := cypher.Node("Person").Named("p")
personProps := person.WithProps(map[string]interface{}{
    "name": "Keanu Reeves",
    "born": 1964,
})

stmt, _ := cypher.Create(personProps).
    Returning(person).
    Build()

fmt.Println(stmt.Cypher())
// CREATE (p:`Person` {name: 'Keanu Reeves', born: 1964}) RETURN p
```

### Merging Data

```go
person := cypher.Node("Person").Named("p")
personProps := person.WithProps(map[string]interface{}{
    "name": "Keanu Reeves",
})

stmt, _ := cypher.Merge(personProps).
    OnCreate(cypher.Set(cypher.Property("p", "created").Eq(2023))).
    OnMatch(cypher.Set(cypher.Property("p", "updated").Eq(2023))).
    Returning(person).
    Build()

fmt.Println(stmt.Cypher())
// MERGE (p:`Person` {name: 'Keanu Reeves'}) ON CREATE SET p.created = 2023 ON MATCH SET p.updated = 2023 RETURN p
```

## Features

### Error Handling

Errors are accumulated during the query building process:

```go
builder := cypher.Match(person).
    Where(person.Property("name").Eq("Tom Hanks")).
    Returning() // Error: empty return clause

if builder.HasError() {
    fmt.Println("Error:", builder.Error())
    return
}
```

### Pretty Printing

Format your Cypher queries for better readability:

```go
// Build a complex query
stmt, _ := cypher.Match(pattern1).
    Where(condition1).
    With(expressions...).
    Match(pattern2).
    Where(condition2).
    Returning(returnExpressions...).
    OrderBy(orderByExpressions...).
    Limit(10).
    Build()

// Format the query with the pretty printer
formattedQuery := cypher.PrettyPrint(stmt.Cypher())
fmt.Println(formattedQuery)
/*
MATCH (p:`Person`)-[r:`ACTED_IN`]->(m:`Movie`)
WHERE p.name = 'Tom Hanks'
WITH p, m
MATCH (m)-[d:`DIRECTED_BY`]->(director:`Person`)
WHERE director.name = 'Steven Spielberg'
RETURN p.name, m.title, director.name
ORDER BY m.year DESC
LIMIT 10
*/
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.