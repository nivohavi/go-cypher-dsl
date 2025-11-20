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

### Function Expression Methods

Chain function expressions with aliases and use DISTINCT:

```go
person := cypher.Node("Person").Named("p")

// Count with alias - fluent chaining
stmt, _ := cypher.Match(person).
    Returning(cypher.Count(person).As("count")).
    Build()

// Count distinct values
stmt, _ := cypher.Match(person).
    Returning(
        cypher.Count(cypher.Distinct(person.Property("name"))).As("uniqueNames"),
    ).
    Build()

// Collect aggregation
stmt, _ := cypher.Match(person).
    Returning(
        cypher.Collect(person.Property("name")).As("names"),
    ).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (p:Person) RETURN collect(p.name) AS names
```

### Auto-quoted Aliases

Aliases are automatically quoted when they contain special characters:

```go
person := cypher.Node("Person").Named("p")

// Automatically quotes aliases with dots, spaces, or special characters
stmt, _ := cypher.Match(person).
    Returning(
        cypher.As(person.Property("name"), "Member.Name"),  // Renders as: p.name AS `Member.Name`
        cypher.As(person.Property("email"), "user_email"),  // Renders as: p.email AS user_email
    ).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (p:Person) RETURN p.name AS `Member.Name`, p.email AS user_email
```

### String Operators

Comprehensive string manipulation functions:

```go
person := cypher.Node("Person").Named("p")

// String concatenation
stmt, _ := cypher.Match(person).
    Returning(
        cypher.Concat(
            person.Property("firstName"),
            cypher.String(" "),
            person.Property("lastName"),
        ).As("fullName"),
    ).
    Build()

// String transformations
stmt, _ := cypher.Match(person).
    Where(
        cypher.Eq(
            cypher.ToLower(person.Property("email")),
            cypher.String("test@example.com"),
        ),
    ).
    Returning(
        cypher.Substring(person.Property("name"), cypher.Integer(0), cypher.Integer(10)).As("shortName"),
        cypher.Replace(person.Property("name"), cypher.String(" "), cypher.String("_")).As("normalizedName"),
        cypher.ToUpper(person.Property("status")),
        cypher.Trim(person.Property("notes")),
    ).
    Build()

// Split strings
stmt, _ := cypher.Match(person).
    Returning(
        cypher.Split(person.Property("tags"), cypher.String(",")).As("tagList"),
    ).
    Build()
```

### Raw Cypher Escape Hatch

For advanced use cases not yet covered by the DSL:

```go
// Insert literal Cypher when needed
stmt, _ := cypher.Match(person).
    Returning(
        cypher.RawCypher("p.name + ' ' + p.surname").As("fullName"),
    ).
    Build()

// Use in WHERE clauses for complex expressions
stmt, _ := cypher.Match(person).
    Where(cypher.RawCypher("p.age > 30 AND p.status = 'active'")).
    Returning(person).
    Build()
```

⚠️ **Security Warning:** Use `RawCypher()` with caution to avoid Cypher injection vulnerabilities. Only use when the DSL doesn't support a specific feature, and ensure any user input is properly sanitized.

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

## Complex Real-World Query Examples

The DSL supports building sophisticated queries for real-world scenarios. Here are some examples:

### E-commerce: Customer Purchase Analysis

```go
customer := cypher.Node("Customer").Named("customer")
order := cypher.Node("Order").Named("order")
product := cypher.Node("Product").Named("product")

rel1 := customer.RelationshipTo(order, "PURCHASED")
rel2 := order.RelationshipTo(product, "CONTAINS")

stmt, _ := cypher.Match(cypher.Path(customer, rel1, order, rel2, product)).
    Returning(
        customer.Property("name"),
        cypher.As(cypher.Count(cypher.Distinct(order)), "totalOrders"),
        cypher.As(cypher.Collect(cypher.Distinct(product.Property("name"))), "productsPurchased"),
    ).
    OrderBy(cypher.Desc(cypher.Var("totalOrders"))).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (customer:Customer)-[:PURCHASED]->(order:Order)-[:CONTAINS]->(product:Product)
// RETURN customer.name, count(DISTINCT order) AS totalOrders, collect(DISTINCT product.name) AS productsPurchased
// ORDER BY totalOrders DESC
```

### Fraud Detection: Suspicious IP Patterns

```go
txn := cypher.Node("Transaction").Named("t")
account := cypher.Node("Account").Named("a")
ip := cypher.Node("IPAddress").Named("ip")

rel1 := txn.RelationshipTo(account, "FROM")
rel2 := account.RelationshipTo(ip, "HAS_IP")

stmt, _ := cypher.Match(cypher.Path(txn, rel1, account, rel2, ip)).
    With(
        ip,
        cypher.As(cypher.Count(cypher.Distinct(account)), "accountCount"),
    ).
    Where(
        cypher.Gt(cypher.Var("accountCount"), cypher.Integer(1)),
    ).
    Returning(
        cypher.As(ip.Property("address"), "SuspiciousIP"),
        cypher.Var("accountCount"),
    ).
    OrderBy(cypher.Desc(cypher.Var("accountCount"))).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (t:Transaction)-[:FROM]->(a:Account)-[:HAS_IP]->(ip:IPAddress)
// WITH ip, count(DISTINCT a) AS accountCount
// WHERE accountCount > 1
// RETURN ip.address AS SuspiciousIP, accountCount
// ORDER BY accountCount DESC
```

### Recommendation System: Collaborative Filtering

```go
u1 := cypher.Node("User").Named("u1")
u2 := cypher.Node("User").Named("u2")
p := cypher.Node("Product").Named("p")
rec := cypher.Node("Product").Named("rec")

rel1 := u1.RelationshipTo(p, "PURCHASED")
rel2 := u2.RelationshipTo(p, "PURCHASED")
rel3 := u2.RelationshipTo(rec, "PURCHASED")

stmt, _ := cypher.Match(cypher.Path(u1, rel1, p)).
    Match(cypher.Path(u2, rel2, p)).
    Match(cypher.Path(u2, rel3, rec)).
    Where(
        cypher.Not(cypher.Pattern(u1, u1.RelationshipTo(rec, "PURCHASED"), rec)),
    ).
    Returning(
        cypher.As(rec.Property("name"), "RecommendedProduct"),
        cypher.As(cypher.Count(cypher.Var("u2")), "RecommendationScore"),
    ).
    OrderBy(cypher.Desc(cypher.Var("RecommendationScore"))).
    Limit(5).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (u1:User)-[:PURCHASED]->(p:Product)
// MATCH (u2:User)-[:PURCHASED]->(p)
// MATCH (u2)-[:PURCHASED]->(rec:Product)
// WHERE NOT (u1)-[:PURCHASED]->(rec)
// RETURN rec.name AS RecommendedProduct, count(u2) AS RecommendationScore
// ORDER BY RecommendationScore DESC
// LIMIT 5
```

### E-commerce: Frequently Bought Together

```go
p1 := cypher.Node("Product").Named("p1")
p2 := cypher.Node("Product").Named("p2")
order := cypher.Node("Order").Named("order")

rel1 := order.RelationshipTo(p1, "CONTAINS")
rel2 := order.RelationshipTo(p2, "CONTAINS")

stmt, _ := cypher.Match(cypher.Path(order, rel1, p1)).
    Match(cypher.Path(order, rel2, p2)).
    Where(
        cypher.Ne(cypher.Var("p1"), cypher.Var("p2")),
    ).
    Returning(
        cypher.As(p1.Property("name"), "Product1"),
        cypher.As(p2.Property("name"), "Product2"),
        cypher.As(cypher.Count(cypher.Var("order")), "timesTogether"),
    ).
    OrderBy(cypher.Desc(cypher.Var("timesTogether"))).
    Limit(10).
    Build()

fmt.Println(stmt.Cypher())
// MATCH (order:Order)-[:CONTAINS]->(p1:Product)
// MATCH (order)-[:CONTAINS]->(p2:Product)
// WHERE p1 <> p2
// RETURN p1.name AS Product1, p2.name AS Product2, count(order) AS timesTogether
// ORDER BY timesTogether DESC
// LIMIT 10
```

## Testing

The package includes comprehensive test coverage for all features, including:
- Unit tests for expression building and rendering
- Integration tests for complex query patterns
- Real-world scenario tests (e-commerce, fraud detection, recommendations)
- Edge case handling and error scenarios

Run tests with:
```bash
go test ./...
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