package main

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/safety"
)

func TestSafety(t *testing.T) {
	// Set safety level to WARN
	safety.SetSafetyLevel(safety.SafetyLevelWarn)

	// Create a buffer to capture log output
	var buf bytes.Buffer
	logger := log.New(&buf, "[test] ", 0)
	safety.SetSafetyLogger(logger)

	// Test a safe string literal
	_ = core.NewLiteral("safe string")
	if buf.Len() > 0 {
		t.Errorf("Expected no warning for safe string, got: %s", buf.String())
	}

	// Reset buffer
	buf.Reset()

	// Test an unsafe string literal with special characters
	_ = core.NewLiteral("unsafe'string;with--special\"chars")
	if buf.Len() == 0 {
		t.Errorf("Expected warning for unsafe string, but got nothing")
	}

	// Reset buffer and test a query using unsafe string literal
	buf.Reset()

	person := cypher.Node("Person").Named("p")

	// Use an unsafe string directly
	stmt, err := cypher.Match(person).
		Where(person.Property("name").Eq("Robert'); DROP TABLE Users;--")).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Failed to build query: %v", err)
	}

	if buf.Len() == 0 {
		t.Errorf("Expected warning for SQL injection string in query, but got nothing")
	}

	// Ensure the cypher still works
	cypherStr := stmt.Cypher()
	if cypherStr == "" {
		t.Error("Expected Cypher query to be generated despite warning")
	}

	// Test with parameters instead (should be safe)
	buf.Reset()

	params := core.NewParameters()
	unsafeName := "Robert'); DROP TABLE Users;--"

	stmt2, err := cypher.Match(person).
		Where(person.Property("name").Eq(params.Add(unsafeName))).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Failed to build query: %v", err)
	}

	// Should not trigger warning with parameters
	if buf.Len() > 0 {
		t.Errorf("Expected no warning when using parameters, got: %s", buf.String())
	}

	// Check that parameters are used
	queryParams := stmt2.Params()
	if len(queryParams) == 0 {
		t.Error("Expected parameters to be present in query")
	}

	// Check that parameter has the unsafe value
	var foundUnsafeValue bool
	for _, v := range queryParams {
		if v == unsafeName {
			foundUnsafeValue = true
			break
		}
	}

	if !foundUnsafeValue {
		t.Error("Expected unsafe value to be stored as parameter")
	}
}

func TestMain(m *testing.M) {
	// Run the tests
	exitCode := m.Run()

	// Restore default logger before exit
	safety.SetSafetyLogger(log.New(os.Stderr, "[cypher-safety] ", log.LstdFlags))

	os.Exit(exitCode)
}
