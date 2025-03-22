# go-cypher-dsl Examples

This directory contains various examples demonstrating how to use the Go Cypher DSL library. The examples are organized into the following categories:

## Directory Structure

```
examples/
├── dsl/               # Pure DSL usage examples without database connectivity
│   ├── simple/        # Basic DSL examples
│   └── complex/       # More complex DSL patterns (WIP)
│
├── neo4j/             # Examples that connect to Neo4j database
│   ├── basic/         # Basic Neo4j connectivity examples
│   └── movie_graph/   # Examples using the Neo4j movie graph dataset
│
└── patterns/          # Common Cypher patterns implemented with the DSL
```

## Running the Examples

Each subdirectory contains standalone Go programs that can be run directly. For example:

```bash
# Run a simple DSL example
cd examples/dsl/simple
go run dsl_examples.go

# Run a Neo4j connection example
cd examples/neo4j/basic
go run neo4j_connection.go
```

## Neo4j Examples

The Neo4j examples require a running Neo4j instance. The movie graph examples use the built-in Neo4j movie dataset, which can be loaded in Neo4j Browser with the `:play movies` command.

Default connection parameters:
- URL: `neo4j://localhost:7687` 
- Username: `neo4j`
- Password: `Password1` (you'll need to change this to match your setup)

## Example Categories

### DSL Examples

The DSL examples demonstrate how to build Cypher queries programmatically without executing them against a database.

### Neo4j Examples

The Neo4j examples show how to connect to a Neo4j database and execute queries built with the DSL.

### Pattern Examples

The pattern examples demonstrate common Cypher patterns implemented with the DSL, such as creating nodes, relationships, and more complex queries. 