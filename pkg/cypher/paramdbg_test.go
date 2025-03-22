package cypher

import (
	"fmt"
	"testing"
)

func TestParamDebug(t *testing.T) {
	// Create a parameter
	param := ParamWithValue("name", "Tom Hanks")
	fmt.Printf("Parameter: %v, name: %s, value: %v\n", param, param.(*parameter).Name(), param.(*parameter).Value())

	// Create a node pattern with the parameter
	person := Node("Person").Named("p").(*nodePattern)
	expr := person.Property("name").Eq(param)

	// Create a simple statement with the parameter
	stmt, err := Match(person).
		Where(expr).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Error building statement: %v", err)
	}

	// Directly check the parameter collection logic
	cypherStmt := stmt.(*cypherStatement)

	// Render the query first
	query := cypherStmt.Cypher()
	fmt.Printf("Query: %s\n", query)

	// Test parameters
	collector := &parameterCollector{params: make(map[string]any)}
	collector.Visit(param)
	fmt.Printf("Direct collector with param: %v\n", collector.params)

	// Try to collect params from the builder
	builderCollector := &parameterCollector{params: make(map[string]any)}
	builderCollector.Visit(cypherStmt.builder)
	fmt.Printf("Builder collector: %v\n", builderCollector.params)

	// Test the actual statement params
	params := stmt.Params()
	fmt.Printf("Statement params: %v\n", params)
}
