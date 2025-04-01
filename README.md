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

### Basic Query Execution

```go
package main

import (
    "fmt"
    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
    "github.com/nivohavi/go-cypher-dsl/pkg/cypher"
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

### Working with Transactions

```go
func createMovieWithActors(driver neo4j.Driver, movieTitle string, releaseYear int, actors []Actor) error {
    session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close()
    
    // Execute within a transaction
    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        // Create movie node
        movieNode := cypher.Node("Movie").Named("m")
        movieWithProps := movieNode.WithProps(map[string]interface{}{
            "title": movieTitle,
            "released": releaseYear,
        })
        
        // Create the movie
        createMovieStmt, err := cypher.Create(movieWithProps).
            Returning(movieNode).
            Build()
            
        if err != nil {
            return nil, err
        }
        
        // Execute movie creation
        result, err := tx.Run(createMovieStmt.Cypher(), createMovieStmt.Params())
        if err != nil {
            return nil, err
        }
        
        // Create actors and relationships
        for _, actor := range actors {
            // Create/merge actor node
            actorNode := cypher.Node("Person").Named("p")
            mergeActorStmt, err := cypher.Merge(actorNode.WithProps(map[string]interface{}{
                "name": actor.Name,
            })).Returning(actorNode).Build()
            
            if err != nil {
                return nil, err
            }
            
            // Execute actor merge
            _, err = tx.Run(mergeActorStmt.Cypher(), mergeActorStmt.Params())
            if err != nil {
                return nil, err
            }
            
            // Create relationship
            actedInRel := actorNode.RelationshipTo(movieNode, "ACTED_IN").WithProps(map[string]interface{}{
                "role": actor.Role,
            })
            
            // Build relationship creation
            createRelStmt, err := cypher.Match(movieNode).
                Where(movieNode.Property("title").Eq(movieTitle)).
                Match(actorNode).
                Where(actorNode.Property("name").Eq(actor.Name)).
                Create(cypher.Pattern(actorNode, actedInRel, movieNode)).
                Build()
                
            if err != nil {
                return nil, err
            }
            
            // Execute relationship creation
            _, err = tx.Run(createRelStmt.Cypher(), createRelStmt.Params())
            if err != nil {
                return nil, err
            }
        }
        
        return result.Single(), nil
    })
    
    return err
}

type Actor struct {
    Name string
    Role string
}

// Usage:
// actors := []Actor{
//     {Name: "Tom Hanks", Role: "Forrest"},
//     {Name: "Robin Wright", Role: "Jenny"},
// }
// err := createMovieWithActors(driver, "Forrest Gump", 1994, actors)
```

### Parameterized Queries with Dynamic Conditions

```go
func findMoviesByFilters(session neo4j.Session, filters MovieFilters) ([]Movie, error) {
    // Create the base query
    movie := cypher.Node("Movie").Named("m")
    matchBuilder := cypher.Match(movie)
    
    // Build dynamic conditions based on filters
    var conditions []cypher.BooleanExpression
    
    if filters.Title != "" {
        conditions = append(conditions, 
            movie.Property("title").Contains(filters.Title))
    }
    
    if filters.MinYear > 0 {
        conditions = append(conditions, 
            movie.Property("released").Gte(filters.MinYear))
    }
    
    if filters.MaxYear > 0 {
        conditions = append(conditions, 
            movie.Property("released").Lte(filters.MaxYear))
    }
    
    if len(filters.Genres) > 0 {
        // Array containment check
        conditions = append(conditions,
            movie.Property("genres").In(filters.Genres))
    }
    
    if filters.DirectorName != "" {
        // Add relationship condition with director
        director := cypher.Node("Person").Named("d")
        directedRel := director.RelationshipTo(movie, "DIRECTED")
        
        // Build the pattern with an additional match
        matchBuilder = matchBuilder.Match(cypher.Pattern(director, directedRel, movie))
        conditions = append(conditions, 
            director.Property("name").Eq(filters.DirectorName))
    }
    
    // Combine all conditions with AND
    if len(conditions) > 0 {
        var condition cypher.BooleanExpression = conditions[0]
        for i := 1; i < len(conditions); i++ {
            condition = condition.And(conditions[i])
        }
        matchBuilder = matchBuilder.Where(condition)
    }
    
    // Add return and limit
    stmt, err := matchBuilder.
        Returning(movie).
        Limit(filters.Limit).
        Build()
        
    if err != nil {
        return nil, err
    }
    
    // Execute the query
    result, err := session.Run(stmt.Cypher(), stmt.Params())
    if err != nil {
        return nil, err
    }
    
    // Process results
    var movies []Movie
    for result.Next() {
        record := result.Record()
        movieNode, _ := record.Get("m")
        
        // Convert Neo4j Node to Movie struct
        movie := nodeToMovie(movieNode.(neo4j.Node))
        movies = append(movies, movie)
    }
    
    return movies, nil
}

type MovieFilters struct {
    Title        string
    MinYear      int
    MaxYear      int
    Genres       []string
    DirectorName string
    Limit        int
}

type Movie struct {
    Title    string
    Released int
    Genres   []string
    // Other properties...
}

func nodeToMovie(node neo4j.Node) Movie {
    props := node.Props()
    
    // Extract genres as string slice
    var genres []string
    if genresValue, ok := props["genres"].([]interface{}); ok {
        for _, g := range genresValue {
            if genre, ok := g.(string); ok {
                genres = append(genres, genre)
            }
        }
    }
    
    // Create movie from node properties
    return Movie{
        Title:    props["title"].(string),
        Released: int(props["released"].(int64)),
        Genres:   genres,
    }
}

// Usage:
// filters := MovieFilters{
//     MinYear: 1990,
//     MaxYear: 2000,
//     DirectorName: "Steven Spielberg",
//     Limit: 5,
// }
// movies, err := findMoviesByFilters(session, filters)
```

### Path Aggregation and Shortest Path

```go
func findShortestPath(session neo4j.Session, startActorName, endActorName string) ([]string, error) {
    // Define nodes
    startActor := cypher.Node("Person").Named("start")
    endActor := cypher.Node("Person").Named("end")
    
    // Create a variable length path between actors
    // This uses a Cypher shortestPath function
    pathExpr := cypher.Function("shortestPath", cypher.PatternPath(
        startActor,
        startActor.RelationshipTo(endActor, "*").Named("r").WithProps(map[string]interface{}{
            "length": "*",  // Variable length path
        }),
        endActor,
    ))
    
    // Build the query
    stmt, err := cypher.Match(startActor).
        Where(startActor.Property("name").Eq(startActorName)).
        Match(endActor).
        Where(endActor.Property("name").Eq(endActorName)).
        With(startActor, endActor).
        Match(cypher.As(pathExpr, "p")).
        Returning(cypher.Name("p")).
        Build()
        
    if err != nil {
        return nil, err
    }
    
    // Execute the query
    result, err := session.Run(stmt.Cypher(), stmt.Params())
    if err != nil {
        return nil, err
    }
    
    // Process result
    if result.Next() {
        record := result.Record()
        path, _ := record.Get("p")
        
        // Extract node names from path
        return extractPathNodeNames(path.(neo4j.Path)), nil
    }
    
    return nil, fmt.Errorf("no path found between %s and %s", startActorName, endActorName)
}

func extractPathNodeNames(path neo4j.Path) []string {
    var names []string
    
    // Add all node names from the path
    for _, node := range path.Nodes() {
        if name, ok := node.Props()["name"].(string); ok {
            names = append(names, name)
        }
    }
    
    return names
}

// Usage:
// path, err := findShortestPath(session, "Kevin Bacon", "Tom Hanks")
// if err == nil {
//     fmt.Println("Bacon path:", strings.Join(path, " -> "))
// }
```

## Examples

For more examples, check the [examples](./examples) directory.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Recent Enhancements

The library has recently been enhanced with several powerful features to make working with Neo4j and Cypher even more efficient:

1. **ComplexPath Helper** - Create multi-level relationship paths more easily without nesting relationship calls.
2. **Property Comparison Helpers** - Reduce verbosity when writing common property comparisons with `CompareProperty` and `NamedCompareProperty`.
3. **Schema Management Helpers** - Simplify the creation and management of constraints and indexes with the `schema` package.
4. **Neo4j Driver Integration** - Streamline query execution and result handling with the `driver` package's `SessionManager` and `QueryHelper`.

See the [enhanced_features examples](examples/enhanced_features) for detailed demonstrations of these new capabilities.