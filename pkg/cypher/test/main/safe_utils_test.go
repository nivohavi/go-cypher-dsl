package main

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestSafeUtilities(t *testing.T) {
	// Test Safe function
	safeExpr := cypher.Safe("dangerous'value")

	// Verify it represents a parameter (by checking if it has a $ in its string representation)
	paramStr := safeExpr.String()
	if !strings.HasPrefix(paramStr, "$") {
		t.Errorf("Expected parameter expression (starting with $), got %s", paramStr)
	}

	// Create a query using Safe
	matchStmt, err := cypher.Match(
		cypher.SafeNodeProps("p", map[string]any{
			"name": "John",
			"age":  30,
		}, "Person"),
	).Where(
		expr.Equals(cypher.Property("p", "name"), cypher.Safe("John")),
	).Returning(
		cypher.Property("p", "name"),
	).Build()

	if err != nil {
		t.Fatalf("Failed to build match statement: %v", err)
	}

	// Verify query contains parameter placeholders and not raw values
	query := matchStmt.Cypher()
	if !strings.Contains(query, "MATCH (p:Person {") {
		t.Errorf("Expected MATCH clause with Person node, got %s", query)
	}
	if !strings.Contains(query, "name: $") {
		t.Errorf("Expected parameterized name property, got %s", query)
	}
	if !strings.Contains(query, "age: $") {
		t.Errorf("Expected parameterized age property, got %s", query)
	}
	if !strings.Contains(query, "WHERE p.name = $") {
		t.Errorf("Expected parameterized WHERE condition, got %s", query)
	}
	if strings.Contains(query, "John") {
		t.Errorf("Found raw string 'John' in query, expected parameterized value: %s", query)
	}
	if strings.Contains(query, "30") {
		t.Errorf("Found raw value '30' in query, expected parameterized value: %s", query)
	}

	// Test SafeMap function
	props := map[string]any{
		"name": "Alice",
		"age":  25,
	}

	safeProps := cypher.SafeMap(props)
	if len(safeProps) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(safeProps))
	}

	// Verify all values are represented as parameters (start with $ in string representation)
	for k, v := range safeProps {
		paramStr := v.String()
		if !strings.HasPrefix(paramStr, "$") {
			t.Errorf("Property %s not parameterized: %s", k, paramStr)
		}
	}

	// Test SafeRelProps function
	person1 := cypher.Node("Person").Named("p1")
	person2 := cypher.Node("Person").Named("p2")

	rel := cypher.SafeRelProps(person1, person2, "KNOWS", map[string]any{
		"since": 2010,
		"type":  "friend",
	})

	matchRelStmt, err := cypher.Match(rel).
		Returning(cypher.Property("p1", "name")).
		Build()

	if err != nil {
		t.Fatalf("Failed to build relationship statement: %v", err)
	}

	// Verify relationship properties are parameterized
	relQuery := matchRelStmt.Cypher()
	if !strings.Contains(relQuery, "MATCH (p1:Person)-[") {
		t.Errorf("Expected MATCH clause with relationship, got %s", relQuery)
	}
	if !strings.Contains(relQuery, "KNOWS {since: $") {
		t.Errorf("Expected parameterized since property, got %s", relQuery)
	}
	if !strings.Contains(relQuery, "type: $") {
		t.Errorf("Expected parameterized type property, got %s", relQuery)
	}
	if strings.Contains(relQuery, "2010") {
		t.Errorf("Found raw value '2010' in query, expected parameterized value: %s", relQuery)
	}
	if strings.Contains(relQuery, "friend") {
		t.Errorf("Found raw string 'friend' in query, expected parameterized value: %s", relQuery)
	}

	// Test ValidateStatement function
	validStmt, err := cypher.Match(cypher.Node("Person").Named("p")).
		Returning(cypher.Property("p", "name")).
		Build()

	if err != nil {
		t.Fatalf("Failed to build valid statement: %v", err)
	}

	valid, errors := cypher.ValidateStatement(validStmt)
	if !valid {
		t.Errorf("Expected valid statement, but validation failed with errors: %v", errors)
	}

	// For an invalid statement, we can test a nil statement
	valid, errors = cypher.ValidateStatement(nil)
	if valid {
		t.Error("Expected invalid statement for nil input, but validation passed")
	}
	if len(errors) == 0 {
		t.Error("Expected error for nil statement, but no errors returned")
	}
}
