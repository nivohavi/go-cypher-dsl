package test

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/errors"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/validation"
)

func TestValidator(t *testing.T) {
	// Create a validator
	validator := validation.DefaultValidator()

	// Test valid query
	validQuery := "MATCH (p:Person) WHERE p.name = 'John' RETURN p"
	if !validator.IsValid(validQuery) {
		t.Errorf("Expected valid query to pass validation: %s", validQuery)
	}

	// Test invalid queries
	invalidQueries := map[string]string{
		"unmatched_parentheses": "MATCH (p:Person WHERE p.name = 'John' RETURN p",
		"unmatched_brackets":    "MATCH (p:Person)-[r:ACTED_IN->([m:Movie) RETURN p",
		"unmatched_braces":      "MATCH (p:Person {name: 'John') RETURN p",
	}

	for name, query := range invalidQueries {
		errs := validator.ValidateQuery(query)
		if len(errs) == 0 {
			t.Errorf("Expected invalid query '%s' to fail validation: %s", name, query)
		}
	}
}

func TestStatementValidation(t *testing.T) {
	// Create node patterns
	person := cypher.Node("Person").Named("p")

	// Build a valid query
	validStmt, err := cypher.Match(person).
		Where(person.Property("name").Eq("John")).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Failed to build query: %v", err)
	}

	// Set a validator on the statement
	validStmt.(*core.StatementImpl).SetValidator(validation.DefaultValidator())

	// Validate the statement
	valid, errs := validStmt.(*core.StatementImpl).Validate()
	if !valid || len(errs) > 0 {
		t.Errorf("Expected valid statement to pass validation: %s", validStmt.Cypher())
	}

	// Create an invalid statement directly
	invalidCypher := "MATCH (p:Person WHERE p.name = 'John' RETURN p"
	invalidStmt := core.NewStatement(invalidCypher, nil)

	// Set a validator on the statement
	invalidStmt.SetValidator(validation.DefaultValidator())

	// Validate the statement
	valid, errs = invalidStmt.Validate()
	if valid || len(errs) == 0 {
		t.Errorf("Expected invalid statement to fail validation: %s", invalidStmt.Cypher())
	}

	// Verify error type
	syntaxErr := false
	for _, err := range errs {
		if errors.IsSyntaxError(err) {
			syntaxErr = true
			break
		}
	}

	if !syntaxErr {
		t.Errorf("Expected syntax error in validation errors for: %s", invalidStmt.Cypher())
	}
}
