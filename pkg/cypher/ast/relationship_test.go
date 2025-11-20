package ast

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

func TestRelationshipPattern(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN")

	result := rel.String()
	// Should contain relationship type
	if !contains(result, "ACTED_IN") {
		t.Errorf("RelationshipTo() = %q, should contain relationship type", result)
	}
}

func TestRelationshipNamed(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN").Named("r")

	result := rel.String()
	// Should contain alias
	if !contains(result, "r") {
		t.Errorf("Relationship().Named() = %q, should contain alias", result)
	}
}

func TestRelationshipWithProperties(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN")
	props := map[string]interface{}{
		"roles": []string{"Lead"},
		"year":  2020,
	}
	relWithProps := rel.WithProps(props)

	result := relWithProps.String()
	// Should contain relationship type
	// Note: Properties are stored but may not be in String() representation
	if !contains(result, "ACTED_IN") {
		t.Errorf("Relationship().WithProps() = %q, should contain ACTED_IN", result)
	}
}

func TestRelationshipProperty(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN").Named("r")
	prop := rel.Property("roles")

	result := prop.String()
	expected := "r.roles"
	if result != expected {
		t.Errorf("Relationship().Property() = %q, want %q", result, expected)
	}
}

func TestRelationshipDirection(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")

	// Outgoing
	relOut := node1.RelationshipTo(node2, "ACTED_IN")
	if relOut.Direction() != core.OUTGOING {
		t.Errorf("RelationshipTo() direction = %v, want OUTGOING", relOut.Direction())
	}

	// Incoming
	relIn := node2.RelationshipFrom(node1, "DIRECTED")
	if relIn.Direction() != core.INCOMING {
		t.Errorf("RelationshipFrom() direction = %v, want INCOMING", relIn.Direction())
	}

	// Bidirectional
	relBoth := node1.RelationshipBetween(node2, "KNOWS")
	if relBoth.Direction() != core.BIDIRECTIONAL {
		t.Errorf("RelationshipBetween() direction = %v, want BIDIRECTIONAL", relBoth.Direction())
	}
}

func TestRelationshipTypes(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN", "PERFORMED_IN")

	types := rel.Types()
	if len(types) != 2 {
		t.Errorf("Types() length = %d, want 2", len(types))
	}
	if types[0] != "ACTED_IN" || types[1] != "PERFORMED_IN" {
		t.Errorf("Types() = %v, want [ACTED_IN, PERFORMED_IN]", types)
	}
}

func TestRelationshipStartEndNodes(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN")

	start := rel.StartNode()
	end := rel.EndNode()

	if start.SymbolicName() != "p" {
		t.Errorf("StartNode() = %q, want 'p'", start.SymbolicName())
	}
	if end.SymbolicName() != "m" {
		t.Errorf("EndNode() = %q, want 'm'", end.SymbolicName())
	}
}

func TestRelationshipLogicalOps(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel1 := node1.RelationshipTo(node2, "ACTED_IN")
	rel2 := node1.RelationshipTo(node2, "DIRECTED")

	// Test And
	andExpr := rel1.And(rel2)
	if andExpr == nil {
		t.Error("rel1.And(rel2) returned nil")
	}

	// Test Or
	orExpr := rel1.Or(rel2)
	if orExpr == nil {
		t.Error("rel1.Or(rel2) returned nil")
	}

	// Test Not
	notExpr := rel1.Not()
	if notExpr == nil {
		t.Error("rel1.Not() returned nil")
	}
}

func TestRelationshipStringRepresentation(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")

	// Outgoing
	relOut := node1.RelationshipTo(node2, "ACTED_IN")
	resultOut := relOut.String()
	if !contains(resultOut, "->") {
		t.Errorf("Outgoing relationship = %q, should contain '->'", resultOut)
	}

	// Incoming
	relIn := node2.RelationshipFrom(node1, "DIRECTED")
	resultIn := relIn.String()
	if !contains(resultIn, "<-") {
		t.Errorf("Incoming relationship = %q, should contain '<-'", resultIn)
	}
}


