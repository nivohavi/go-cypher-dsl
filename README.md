# go-cypher-dsl

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A fluent, type-safe Cypher query builder for Go. Build Neo4j queries programmatically with syntax completion, reducing errors and improving development speed.

## Installation

```bash
go get github.com/yourusername/go-cypher-dsl
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/yourusername/go-cypher-dsl/pkg/cypher"
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

## Core Concepts

### Nodes and Relationships

```go
// Create nodes
person := cypher.Node("Person").Named("p")
movie := cypher.Node("Movie").Named("m")

// Create a relationship
acted := person.RelationshipTo(movie, "ACTED_IN").Named("r")

// Create a pattern
pattern := cypher.Pattern(person, acted, movie)

// Use in a query
query, _ := cypher.Match(pattern).
    Returning(person, movie).
    Build()

fmt.Println(query.Cypher())
// MATCH (p:`Person`)-[r:`ACTED_IN`]->(m:`Movie`) RETURN p, m
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

### Type Conversion

The library automatically converts Go values to Cypher expressions:

```go
// Automatically converts 1950 to a literal
person.Property("born").Gt(1950)

// Automatically converts string array to literal array
genres := []string{"Action", "Drama"}
movie.Property("genres").In(genres)
```

### Property Shorthand

Use `Prop()` as a shorter alternative to `Property()`:

```go
// These are equivalent
person.Property("name").Eq("Tom Hanks")
person.Prop("name").Eq("Tom Hanks")
```

### Schema-First Approach

Define schemas for better type safety:

```go
// Define a schema
personSchema := cypher.NewTypedSchema("Person")
nameProperty := personSchema.AddProperty("name")
ageProperty := personSchema.AddProperty("age")

// Create a node using the schema
p := personSchema.Node("p")

// Build a query using schema properties
stmt, _ := cypher.Match(p).
    Where(nameProperty.Of(p).Eq("Tom Hanks")).
    Returning(ageProperty.Of(p)).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (p:`Person`) WHERE p.name = 'Tom Hanks' RETURN p.age
```

## Common Query Patterns

### Finding Nodes

```go
// Find nodes by properties
findQuery := cypher.FindNodesByProperties("Person", map[string]interface{}{
    "name": "Tom Hanks",
    "born": 1956,
})

fmt.Println(findQuery.Cypher())
// MATCH (n:`Person`) WHERE n.name = 'Tom Hanks' AND n.born = 1956 RETURN n
```

### Creating Nodes

```go
// Create a node with properties
createQuery := cypher.CreateNodeWithProperties("Actor", map[string]interface{}{
    "name": "Keanu Reeves",
    "born": 1964,
})

fmt.Println(createQuery.Cypher())
// CREATE (n:`Actor` {name: 'Keanu Reeves', born: 1964}) RETURN n
```

### Creating Relationships

```go
// Get existing nodes
person := cypher.Node("Person").Named("p")
movie := cypher.Node("Movie").Named("m")

// Create a relationship between nodes
relQuery := cypher.RelateNodes(person, "ACTED_IN", movie, map[string]interface{}{
    "role": "Neo",
})

fmt.Println(relQuery.Cypher())
// CREATE (p:`Person`)-[:`ACTED_IN` {role: 'Neo'}]->(m:`Movie`) RETURN p, r, m
```

### Deleting Nodes

```go
person := cypher.Node("Person").Named("p")

// Delete a node (with optional detach)
deleteQuery := cypher.DeleteNode(person, true) // true for DETACH DELETE
```

## Advanced Examples

### Complex Query with Multiple Clauses

```go
person := cypher.Node("Person").Named("p")
movie := cypher.Node("Movie").Named("m")
acted := person.RelationshipTo(movie, "ACTED_IN").Named("r")

stmt, _ := cypher.Match(cypher.Pattern(person, acted, movie)).
    Where(person.Prop("name").Eq("Tom Hanks")).
    With(person, movie, cypher.Count(acted).As("roles")).
    Where(cypher.Name("roles").Gt(1)).
    Returning(movie.Prop("title"), cypher.Name("roles")).
    OrderByDesc(cypher.Name("roles")).
    Limit(10).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (p:`Person`)-[r:`ACTED_IN`]->(m:`Movie`) 
// WHERE p.name = 'Tom Hanks' 
// WITH p, m, count(r) AS roles 
// WHERE roles > 1 
// RETURN m.title, roles 
// ORDER BY roles DESC 
// LIMIT 10
```

### Combining Multiple Statements with Transactions

```go
// Create a transaction
tx := cypher.NewTransaction()

// Add multiple statements
matchStmt, _ := cypher.Match(person).
    Where(person.Prop("name").Eq("Tom Hanks")).
    Returning(person).
    Build()

createStmt, _ := cypher.Create(movie).
    Returning(movie).
    Build()

tx.Add(matchStmt)
tx.Add(createStmt)

// Execute the transaction with Neo4j
```

## Integration with Neo4j Driver

```go
package main

import (
    "context"
    "fmt"
    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
    "github.com/yourusername/go-cypher-dsl/pkg/cypher"
)

func main() {
    // Create driver
    driver, err := neo4j.NewDriver("neo4j://localhost:7687", 
        neo4j.BasicAuth("neo4j", "password", ""))
    if err != nil {
        panic(err)
    }
    defer driver.Close()

    session := driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    // Build query with go-cypher-dsl
    person := cypher.Node("Person").Named("p")
    movie := cypher.Node("Movie").Named("m")
    acted := person.RelationshipTo(movie, "ACTED_IN").Named("r")
    
    stmt, _ := cypher.Match(cypher.Pattern(person, acted, movie)).
        Where(person.Prop("name").Eq("Tom Hanks")).
        Returning(movie.Prop("title").As("title"), 
                  acted.Prop("role").As("role")).
        Build()
    
    // Execute the query
    result, err := session.Run(
        context.Background(),
        stmt.Cypher(),
        stmt.Params(),
    )
    
    if err != nil {
        panic(err)
    }
    
    // Process results
    for result.Next(context.Background()) {
        record := result.Record()
        title, _ := record.Get("title")
        role, _ := record.Get("role")
        fmt.Printf("Movie: %s, Role: %s\n", title, role)
    }
}
```

## Examples

For more examples, check the [examples](./examples) directory.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.