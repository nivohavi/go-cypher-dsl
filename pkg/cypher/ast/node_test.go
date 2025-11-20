package ast

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestNode(t *testing.T) {
	tests := []struct {
		name     string
		labels   []string
		expected string
	}{
		{"single label", []string{"Person"}, "(:Person)"},
		{"multiple labels", []string{"Person", "Actor"}, "(:Person:Actor)"},
		{"no labels", []string{}, "()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Node(tt.labels...)
			result := node.String()
			if result != tt.expected {
				t.Errorf("Node(%v).String() = %q, want %q", tt.labels, result, tt.expected)
			}
		})
	}
}

func TestNodeNamed(t *testing.T) {
	node := Node("Person").Named("p")
	result := node.String()
	expected := "(p:Person)"
	if result != expected {
		t.Errorf("Node(\"Person\").Named(\"p\").String() = %q, want %q", result, expected)
	}
}

func TestNodeWithLabels(t *testing.T) {
	node := Node("Person").WithLabels("Actor")
	result := node.String()
	// Should contain both labels
	if !contains(result, "Person") || !contains(result, "Actor") {
		t.Errorf("Node().WithLabels() = %q, should contain both labels", result)
	}
}

func TestNodeWithProperties(t *testing.T) {
	node := Node("Person").Named("p")
	props := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	nodeWithProps := node.WithProps(props)
	result := nodeWithProps.String()
	// Should contain properties
	if !contains(result, "name") || !contains(result, "age") {
		t.Errorf("Node().WithProps() = %q, should contain properties", result)
	}
}

func TestNodeWithPropertiesExpressions(t *testing.T) {
	node := Node("Person").Named("p")
	props := map[string]interface{}{
		"name": expr.String("John"),
		"age":  expr.Integer(30),
	}
	nodeWithProps := node.WithProps(props)
	result := nodeWithProps.String()
	// Should contain properties
	if !contains(result, "name") || !contains(result, "age") {
		t.Errorf("Node().WithProps() = %q, should contain properties", result)
	}
}

func TestNodeProperty(t *testing.T) {
	node := Node("Person").Named("p")
	prop := node.Property("name")
	result := prop.String()
	expected := "p.name"
	if result != expected {
		t.Errorf("Node().Property() = %q, want %q", result, expected)
	}
}

func TestNodeRelationshipTo(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN")
	result := rel.String()
	// Should contain relationship type
	if !contains(result, "ACTED_IN") {
		t.Errorf("RelationshipTo() = %q, should contain relationship type", result)
	}
}

func TestNodeRelationshipFrom(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node2.RelationshipFrom(node1, "DIRECTED")
	result := rel.String()
	// Should contain relationship type
	if !contains(result, "DIRECTED") {
		t.Errorf("RelationshipFrom() = %q, should contain relationship type", result)
	}
}

func TestNodeRelationshipBetween(t *testing.T) {
	node1 := Node("Person").Named("p1")
	node2 := Node("Person").Named("p2")
	rel := node1.RelationshipBetween(node2, "KNOWS")
	result := rel.String()
	// Should contain relationship type
	if !contains(result, "KNOWS") {
		t.Errorf("RelationshipBetween() = %q, should contain relationship type", result)
	}
}

func TestNodeSymbolicName(t *testing.T) {
	node := Node("Person").Named("p")
	name := node.SymbolicName()
	if name != "p" {
		t.Errorf("SymbolicName() = %q, want %q", name, "p")
	}
}

func TestNodeLogicalOps(t *testing.T) {
	node1 := Node("Person").Named("p1")
	node2 := Node("Person").Named("p2")

	// Test And
	andExpr := node1.And(node2)
	if andExpr == nil {
		t.Error("node1.And(node2) returned nil")
	}

	// Test Or
	orExpr := node1.Or(node2)
	if orExpr == nil {
		t.Error("node1.Or(node2) returned nil")
	}

	// Test Not
	notExpr := node1.Not()
	if notExpr == nil {
		t.Error("node1.Not() returned nil")
	}
}


