# go-cypher-dsl Examples

This directory contains various examples demonstrating how to use the Go Cypher DSL library. The examples are organized into the following categories:

## Directory Structure

```
examples/
├── basics/              # Core query building examples using improved API
├── enhanced_features/   # Examples showcasing new advanced features
├── neo4j_integration/   # Examples of Neo4j database integration
└── schema_management/   # Examples of schema management (constraints and indexes)
```

## Running the Examples

Each subdirectory contains standalone Go programs that can be run directly. For example:

```bash
# Run the basic query building example
cd examples/basics
go run query_building.go

# Run the enhanced features example
cd examples/enhanced_features
go run example.go

# Run the Neo4j integration example
cd examples/neo4j_integration
go run movie_graph_queries.go
```

## Neo4j Examples

The Neo4j integration examples require a running Neo4j instance. The movie graph examples use the built-in Neo4j movie dataset, which can be loaded in Neo4j Browser with the `:play movies` command.

Default connection parameters:
- URL: `neo4j://localhost:7687` 
- Username: `neo4j`
- Password: `password` (you'll need to change this to match your setup)

## Example Categories

### Basics

The basics examples demonstrate fundamental query building with the improved API:
- Creating nodes and relationships
- Building MATCH, CREATE, MERGE and DELETE queries
- Adding WHERE conditions
- Property handling
- Pattern matching
- Complex logical expressions

### Enhanced Features

The enhanced features examples demonstrate the newest capabilities of the library:
- **ComplexPath**: Simplified creation of multi-level relationship paths
- **Property Comparison Helpers**: More concise property comparisons with `CompareProperty` and `NamedCompareProperty`
- **Schema Management**: Easier creation of constraints and indexes
- **Neo4j Driver Integration**: Simplified execution with `SessionManager` and result handling

### Neo4j Integration

The Neo4j integration examples show how to connect to a Neo4j database and execute queries built with the DSL, including:
- Traditional Neo4j session approach
- Improved SessionManager approach with various query helpers
- Result handling patterns

### Schema Management

The schema management examples demonstrate how to create and manage database constraints and indexes using the DSL. 