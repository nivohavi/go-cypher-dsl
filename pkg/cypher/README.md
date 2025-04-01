# go-cypher-dsl Package

This package provides a fluent, type-safe way to build Cypher queries for Neo4j in Go. The DSL (Domain-Specific Language) approach helps reduce errors, improves maintainability, and provides better IDE support compared to string concatenation.

## Core Concepts

### Expressions

Everything in the DSL is based on `Expression` interfaces. An expression represents a part of a Cypher query that can be rendered as a string and may contain parameters.

Key expression types:
- Node patterns: `(n:Label)`
- Relationship patterns: `-[r:TYPE]->`
- Property expressions: `n.property`
- Comparison expressions: `n.property = value`
- Logical expressions: `expr1 AND expr2`
- Function calls: `count(n)`

### Builders

Builders provide a fluent API for constructing Cypher queries:
- `MatchBuilder`: For building MATCH clauses
- `CreateBuilder`: For building CREATE clauses
- `ReturnBuilder`: For building RETURN clauses
- `WithBuilder`: For building WITH clauses
- and more...

### Parameters

Parameters are handled automatically, providing protection against Cypher injection. Parameter values can be:
- Literal values: `"string"`, `42`, `true`
- Variable values (from Go variables)
- Named parameters that can be referenced in queries

## Basic Usage

### Creating Nodes and Relationships

```go
// Create a node pattern
person := cypher.Node("Person").Named("p")

// Create a relationship
movie := cypher.Node("Movie").Named("m")
actedIn := person.RelationshipTo(movie, "ACTED_IN")
```

### Building Queries

```go
// Create a MATCH-RETURN query
stmt, err := cypher.Match(person).
    Where(person.Property("name").Eq("Tom Hanks")).
    Returning(person).
    Build()

if err != nil {
    // Handle error
}

// Get the query string and parameters
query := stmt.Cypher()
params := stmt.Params()
```

### Using with Neo4j Driver

```go
// Create a Neo4j session
session := driver.NewSession(neo4j.SessionConfig{
    AccessMode: neo4j.AccessModeRead,
})
defer session.Close()

// Execute the query
result, err := session.Run(
    stmt.Cypher(),
    stmt.Params(),
)
```

## Best Practices

### 1. Use Named Nodes and Relationships

Always name your nodes and relationships to make the resulting Cypher query more readable and to allow for referencing them later:

```go
// Good
person := cypher.Node("Person").Named("p")

// Avoid
person := cypher.Node("Person") // No name
```

### 2. Use Parameter Values

For better security and performance, always use parameters for values rather than string concatenation:

```go
// Good
name := "Tom Hanks"
condition := person.Property("name").Eq(cypher.ParamWithValue("name", name))

// Avoid
condition := person.Property("name").Eq("Tom Hanks") // Hardcoded value
```

### 3. Reuse Expressions

Create expressions once and reuse them for better maintainability:

```go
// Create expressions once
person := cypher.Node("Person").Named("p")
nameProperty := person.Property("name")

// Reuse in multiple places
condition1 := nameProperty.Eq("Tom Hanks")
condition2 := nameProperty.Contains("Tom")
```

### 4. Handle Errors Properly

Always check for errors after building queries:

```go
stmt, err := builder.Build()
if err != nil {
    // Handle the error
    return nil, err
}
```

### 5. Use Schemas for Type Safety

For larger projects, define schemas to get better type safety:

```go
// Define a schema
personSchema := cypher.NewTypedSchema("Person")
nameProperty := personSchema.AddProperty("name")
ageProperty := personSchema.AddProperty("age")

// Create a node using the schema
p := personSchema.Node("p")

// Use schema properties
condition := nameProperty.Of(p).Eq("Tom Hanks")
```

## Common Patterns

### Find Nodes by ID

```go
person := cypher.Node("Person").Named("p")
idCondition := person.Property("id").Eq(123)

stmt, _ := cypher.Match(person).
    Where(idCondition).
    Returning(person).
    Build()
```

### Create a Node with Properties

```go
person := cypher.Node("Person").Named("p")
personProps := person.WithProps(map[string]interface{}{
    "name": "John Doe",
    "age":  30,
})

stmt, _ := cypher.Create(personProps).
    Returning(person).
    Build()
```

### Create a Relationship

```go
person := cypher.Node("Person").Named("p")
movie := cypher.Node("Movie").Named("m")
actedIn := person.RelationshipTo(movie, "ACTED_IN")

stmt, _ := cypher.Match(person).
    Where(person.Property("name").Eq("Tom Hanks")).
    Match(movie).
    Where(movie.Property("title").Eq("Forrest Gump")).
    Create(actedIn).
    Returning(person, movie).
    Build()
```

### Pagination

```go
person := cypher.Node("Person").Named("p")

stmt, _ := cypher.Match(person).
    Returning(person).
    OrderBy(person.Property("name")).
    Skip(20).
    Limit(10).
    Build()
```

### Aggregation

```go
movie := cypher.Node("Movie").Named("m")
director := cypher.Node("Person").Named("d")
directedRel := movie.RelationshipTo(director, "DIRECTED_BY")

stmt, _ := cypher.Match(cypher.Pattern(movie, directedRel, director)).
    Returning(
        director.Property("name"),
        cypher.Count(movie).As("movieCount"),
    ).
    OrderBy(cypher.Desc(cypher.Var("movieCount"))).
    Limit(10).
    Build()
```

## Advanced Features

### Path Variables

```go
// Define path pattern
person1 := cypher.Node("Person").Named("p1")
person2 := cypher.Node("Person").Named("p2")
movie := cypher.Node("Movie").Named("m")

path := cypher.Path(
    person1,
    person1.RelationshipTo(movie, "ACTED_IN"),
    movie,
    movie.RelationshipTo(person2, "ACTED_IN"),
    person2,
).Named("path")

stmt, _ := cypher.Match(path).
    Where(person1.Property("name").Eq("Tom Hanks")).
    And(person2.Property("name").Eq("Kevin Bacon")).
    Returning(path).
    Build()
```

### Unwind Operations

```go
names := []string{"Tom Hanks", "Tom Cruise", "Tom Holland"}
namesParam := cypher.ParamWithValue("names", names)

person := cypher.Node("Person").Named("p")

stmt, _ := cypher.Unwind(namesParam, "name").
    Match(person).
    Where(person.Property("name").Eq(cypher.Var("name"))).
    Returning(person).
    Build()
```

### Subqueries

```go
// Main query
person := cypher.Node("Person").Named("p")
movie := cypher.Node("Movie").Named("m")
actedIn := person.RelationshipTo(movie, "ACTED_IN")

// Build with a subquery using WITH
stmt, _ := cypher.Match(person).
    Where(person.Property("name").Eq("Tom Hanks")).
    With(person).
    Match(cypher.Pattern(person, actedIn, movie)).
    Returning(movie).
    Build()
```

## Formatting

The library includes a formatter for pretty-printing Cypher queries:

```go
// Get the query string
query := stmt.Cypher()

// Format with default options
formatted := cypher.PrettyPrint(query)

// Format with custom options
formatter := cypher.NewCypherFormatter(cypher.FormattingOptions{
    IndentString:     "    ",  // 4 spaces
    KeywordCase:      cypher.KeywordCaseUpper,
    ClauseNewline:    true,
    IndentSubClauses: true,
})
formatted = formatter.Format(query)
```

## Error Handling

Errors are accumulated during query building:

```go
builder := cypher.Match(person).
    Where(person.Property("name").Eq("Tom Hanks")).
    Returning() // Error: empty return clause

if builder.HasError() {
    fmt.Println("Error:", builder.Error())
}
```

## Contributing

Contributions to improve the DSL are welcome! Please see the main project README for guidelines. 