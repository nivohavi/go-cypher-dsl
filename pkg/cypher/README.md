# Go Cypher DSL

A fluent, type-safe domain-specific language (DSL) for building Cypher queries in Go.

## Features

- Fluent, chainable API that closely mirrors Cypher syntax
- Type-safe query construction to catch errors at compile time
- Parameter support to prevent Cypher injection
- Pretty printing for debugging
- Full support for all major Cypher clauses and expressions
- Immutable builders that can be reused and composed

## Installation

```bash
go get github.com/nivohavi/go-cypher-dsl
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/nivohavi/go-cypher-dsl/pkg2/cypher"
)

func main() {
    // Build a simple query
    query := cypher.Match(
        cypher.Node("Person").Named("p").WithProps(map[string]interface{}{
            "name": "John",
        }),
    ).
    Return(
        cypher.Property("p", "name"),
        cypher.Property("p", "age"),
    ).
    Build()

    // Render the query for debugging
    fmt.Println(cypher.PrettyPrint(query))

    // Get the Cypher string and parameters for execution
    cypherStr, params := cypher.RenderWithParams(query)
    fmt.Println("Cypher:", cypherStr)
    fmt.Println("Params:", params)
}
```

### More Complex Example

```go
package main

import (
    "fmt"
    "github.com/nivohavi/go-cypher-dsl/pkg2/cypher"
)

func main() {
    // Create parameter container
    params := cypher.Parameters()

    // Build a more complex query with parameters
    person := cypher.NamedNode("p", "Person")
    movie := cypher.NamedNode("m", "Movie")
    
    // Create a relationship pattern
    acted := cypher.RelateNodes(person, movie, "ACTED_IN")
    
    query := cypher.Match(cypher.Pattern(acted)).
        Where(
            cypher.And(
                cypher.Gt(
                    cypher.Property("m", "released"),
                    params.Add(2000),
                ),
                cypher.Eq(
                    cypher.Property("p", "name"),
                    params.Add("Tom Hanks"),
                ),
            ),
        ).
        Return(
            cypher.Property("m", "title"),
            cypher.Property("m", "released"),
        ).
        OrderBy(cypher.Property("m", "released")).
        Desc().
        Build()

    // Print the pretty query
    fmt.Println(cypher.PrettyPrint(query))
}
```

## API Reference

### Node Pattern Functions

- `Node(labels ...string)` - Creates a node pattern with optional labels
- `NamedNode(alias string, labels ...string)` - Creates a named node with optional labels
- `NodeWithProperties(properties map[string]any, labels ...string)` - Creates a node with properties
- `NamedNodeWithProperties(alias string, properties map[string]any, labels ...string)` - Creates a named node with properties

### Relationship Functions

- `RelateNodes(left, right NodeExpression, types ...string)` - Creates a relationship between nodes
- `RelateBidirectionally(left, right NodeExpression, types ...string)` - Creates a bidirectional relationship

### Pattern Functions

- `Pattern(element PatternElement)` - Creates a pattern expression
- `Path(alias string, element PatternElement)` - Creates a named path

### Clause Functions

- `Match(pattern Expression)` - Creates a MATCH clause
- `OptionalMatch(pattern Expression)` - Creates an OPTIONAL MATCH clause
- `Where(condition Expression)` - Creates a WHERE clause
- `With(expressions ...Expression)` - Creates a WITH clause
- `Return(expressions ...Expression)` - Creates a RETURN clause
- `OrderBy(expressions ...Expression)` - Creates an ORDER BY clause
- `Limit(count int)` - Creates a LIMIT clause
- `Skip(count int)` - Creates a SKIP clause
- `Create(pattern Expression)` - Creates a CREATE clause
- `Merge(pattern Expression)` - Creates a MERGE clause
- `Delete(expressions ...Expression)` - Creates a DELETE clause
- `DetachDelete(expressions ...Expression)` - Creates a DETACH DELETE clause
- `Set(expression Expression)` - Creates a SET clause
- `Remove(expression Expression)` - Creates a REMOVE clause
- `Unwind(expression Expression, alias string)` - Creates an UNWIND clause

### Expression Functions

- `Literal(value any)` - Creates a literal expression
- `Param(value any)` - Creates a parameter expression
- `NamedParam(name string, value any)` - Creates a named parameter
- `Eq(left, right Expression)` - Creates an equality expression
- `Ne(left, right Expression)` - Creates a not-equal expression
- `Gt(left, right Expression)` - Creates a greater-than expression
- `Gte(left, right Expression)` - Creates a greater-than-or-equal expression
- `Lt(left, right Expression)` - Creates a less-than expression
- `Lte(left, right Expression)` - Creates a less-than-or-equal expression
- `And(left, right Expression)` - Creates a logical AND expression
- `Or(left, right Expression)` - Creates a logical OR expression
- `Xor(left, right Expression)` - Creates a logical XOR expression
- `Not(expression Expression)` - Creates a logical NOT expression

### Rendering Functions

- `Render(statement Statement)` - Renders a statement to a Cypher string
- `RenderWithParams(statement Statement)` - Renders a statement and returns Cypher and parameters
- `PrettyPrint(statement Statement)` - Renders a statement with pretty printing

## License

MIT 