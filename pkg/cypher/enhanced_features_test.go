package cypher

import (
	"testing"
)

func TestPropertyShorthand(t *testing.T) {
	// Define a node
	person := Node("Person").Named("p")

	// Use Prop() shorthand
	stmt, err := Match(person).
		Where(person.Prop("name").Eq("Tom Hanks")).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (p:`Person`) WHERE p.name = 'Tom Hanks' RETURN p"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestAutoTypeConversion(t *testing.T) {
	// Define a node
	movie := Node("Movie").Named("m").(*nodePattern)

	// Use automatic type conversion
	stmt, err := Match(movie).
		Where(movie.Property("released").Gt(2000)).
		Returning(movie).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (m:`Movie`) WHERE m.released > 2000 RETURN m"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestErrorHandling(t *testing.T) {
	// Define a node
	person := Node("Person").Named("p")

	// Check for errors in the chain
	builder := Match(person).
		Where(person.Property("name").Eq("Tom Hanks")).
		Returning(person)

	if builder.HasError() {
		t.Fatalf("Expected no error, got: %v", builder.Error())
	}

	// Create a builder with an error (empty expressions)
	errorBuilder := Match(person).
		Returning()

	if !errorBuilder.HasError() {
		t.Fatalf("Expected an error for empty return clause, got none")
	}
}

func TestSchemaFirstApproach(t *testing.T) {
	// Define a schema
	personSchema := NewTypedSchema("Person")
	nameProperty := personSchema.AddProperty("name")
	ageProperty := personSchema.AddProperty("age")

	// Create a node using the schema
	p := personSchema.Node("p")

	// Build a query using schema properties
	stmt, err := Match(p).
		Where(nameProperty.Of(p).Eq("Tom Hanks")).
		Returning(ageProperty.Of(p)).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (p:`Person`) WHERE p.name = 'Tom Hanks' RETURN p.age"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestHighLevelPatterns(t *testing.T) {
	// Test FindNodesByProperties
	findQuery := FindNodesByProperties("Person", map[string]interface{}{
		"name": "Tom Hanks",
		"born": 1956,
	})

	if findQuery.HasError() {
		t.Fatalf("Error in FindNodesByProperties: %v", findQuery.Error())
	}

	// Test CreateNodeWithProperties
	createQuery := CreateNodeWithProperties("Actor", map[string]interface{}{
		"name": "Keanu Reeves",
		"born": 1964,
	})

	if createQuery.HasError() {
		t.Fatalf("Error in CreateNodeWithProperties: %v", createQuery.Error())
	}

	// Test RelateNodes
	person := Node("Person").Named("p")
	movie := Node("Movie").Named("m")
	relQuery := RelateNodes(person, "ACTED_IN", movie, map[string]interface{}{
		"role": "Neo",
	})

	if relQuery.HasError() {
		t.Fatalf("Error in RelateNodes: %v", relQuery.Error())
	}
}
