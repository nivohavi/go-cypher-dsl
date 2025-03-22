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

// Create a pattern with the relationship
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
// Values are automatically converted to appropriate Cypher expressions
condition := person.Property("born").Gt(1950)

// Use parameters for query parameterization
paramCondition := movie.Property("released").Gt(cypher.ParamWithValue("year", 2000))

// Arrays are handled nicely too
genres := []string{"Action", "Drama"}
inGenres := movie.Property("genres").In(genres)
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
person := cypher.Node("Person").Named("p")

// Set properties using a map
personWithProps := person.WithProps(map[string]interface{}{
    "name": "John Doe",
    "age":  30,
})

// Build the query
stmt, _ := cypher.Create(personWithProps).
    Returning(personWithProps).
    Build()

fmt.Println(stmt.Cypher())
// CREATE (p:`Person` {name: 'John Doe', age: 30}) RETURN p
```

### Creating Relationships

```go
// Create relationship between nodes
person := cypher.Node("Person").Named("p")
movie := cypher.Node("Movie").Named("m")

// Create a relationship pattern
actedIn := person.RelationshipTo(movie, "ACTED_IN")

// Add properties to the relationship
actedInWithProps := actedIn.WithProps(map[string]interface{}{
    "role": "Neo",
})

// Create a pattern with the relationship
pattern := cypher.Pattern(person, actedInWithProps, movie)

// Build the query
stmt, _ := cypher.Create(pattern).
    Returning(person, actedInWithProps, movie).
    Build()

fmt.Println(stmt.Cypher())
// CREATE (p:`Person`)-[:`ACTED_IN` {role: 'Neo'}]->(m:`Movie`) RETURN p, r, m
```

### Deleting Nodes

```go
// Delete a node
person := cypher.Node("Person").Named("p")

// Create a delete statement
deleteStmt, _ := cypher.Delete(person).Build()

// For detach delete
detachDeleteStmt, _ := cypher.DetachDelete(person).Build()
```

## Advanced Examples

### Complex Query with Multiple Clauses

```go
person := cypher.Node("Person").Named("p")
movie := cypher.Node("Movie").Named("m")
acted := person.RelationshipTo(movie, "ACTED_IN").Named("r")

// Create the pattern
pattern := cypher.Pattern(person, acted, movie)

// Create a count expression
rolesCount := cypher.As(cypher.Function("count", acted), "roles")

stmt, _ := cypher.Match(pattern).
    Where(person.Property("name").Eq("Tom Hanks")).
    With(person, movie, rolesCount).
    Where(cypher.Name("roles").Gt(1)).
    Returning(movie.Property("title"), cypher.Name("roles")).
    OrderBy(cypher.Desc(cypher.Name("roles"))).
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

### Finding Co-actors Pattern

```go
// Tom Hanks' co-actors pattern
tom := cypher.Node("Person").Named("tom")
movie := cypher.Node("Movie").Named("m")
coActors := cypher.Node("Person").Named("coActors")

// Create the complex relationship pattern
tomToMovie := tom.RelationshipTo(movie, "ACTED_IN")
movieToCoActors := movie.RelationshipFrom(coActors, "ACTED_IN")

// Create the path pattern with both relationships
path := cypher.Pattern(tom, tomToMovie, movie, movieToCoActors, coActors)

stmt, _ := cypher.Match(path).
    Where(tom.Property("name").Eq("Tom Hanks")).
    Returning(coActors.Property("name")).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (tom:`Person`)-[:`ACTED_IN`]->(m:`Movie`)<-[:`ACTED_IN`]-(coActors:`Person`) 
// WHERE tom.name = 'Tom Hanks' 
// RETURN coActors.name
```

## Integration with Neo4j Driver

```go
package main

import (
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
    
    // Create the pattern
    pattern := cypher.Pattern(person, acted, movie)
    
    // Build the statement
    stmt, _ := cypher.Match(pattern).
        Where(person.Property("name").Eq("Tom Hanks")).
        Returning(
            cypher.As(movie.Property("title"), "title"),
            cypher.As(acted.Property("role"), "role"),
        ).
        Build()
    
    // Execute the query
    result, err := session.Run(
        stmt.Cypher(),
        stmt.Params(),
    )
    
    if err != nil {
        panic(err)
    }
    
    // Process results
    for result.Next() {
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