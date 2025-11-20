package ast

import (
	"testing"
)

func TestPattern(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN")

	pattern := Pattern(node1, rel, node2)
	result := pattern.String()

	// Should contain both nodes and relationship
	if !contains(result, "p") || !contains(result, "m") || !contains(result, "ACTED_IN") {
		t.Errorf("Pattern() = %q, should contain nodes and relationship", result)
	}
}

func TestPath(t *testing.T) {
	node1 := Node("Person").Named("p")
	node2 := Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN")

	path := Path(node1, rel, node2)
	result := path.String()

	// Should contain both nodes and relationship
	if !contains(result, "p") || !contains(result, "m") || !contains(result, "ACTED_IN") {
		t.Errorf("Path() = %q, should contain nodes and relationship", result)
	}
}

func TestComplexPattern(t *testing.T) {
	person := Node("Person").Named("p")
	movie := Node("Movie").Named("m")
	director := Node("Person").Named("d")

	actedIn := person.RelationshipTo(movie, "ACTED_IN")
	directed := movie.RelationshipTo(director, "DIRECTED_BY")

	// Create a pattern with multiple relationships
	pattern1 := Pattern(person, actedIn, movie)
	pattern2 := Pattern(movie, directed, director)

	result1 := pattern1.String()
	result2 := pattern2.String()

	if !contains(result1, "ACTED_IN") {
		t.Errorf("Complex pattern 1 = %q, should contain ACTED_IN", result1)
	}
	if !contains(result2, "DIRECTED_BY") {
		t.Errorf("Complex pattern 2 = %q, should contain DIRECTED_BY", result2)
	}
}

func TestPatternLogicalOps(t *testing.T) {
	node1 := Node("Person").Named("p1")
	node2 := Node("Person").Named("p2")
	rel := node1.RelationshipTo(node2, "KNOWS")
	pattern1 := Pattern(node1, rel, node2)

	node3 := Node("Person").Named("p3")
	rel2 := node1.RelationshipTo(node3, "KNOWS")
	pattern2 := Pattern(node1, rel2, node3)

	// Test And
	andExpr := pattern1.And(pattern2)
	if andExpr == nil {
		t.Error("pattern1.And(pattern2) returned nil")
	}

	// Test Or
	orExpr := pattern1.Or(pattern2)
	if orExpr == nil {
		t.Error("pattern1.Or(pattern2) returned nil")
	}

	// Test Not
	notExpr := pattern1.Not()
	if notExpr == nil {
		t.Error("pattern1.Not() returned nil")
	}
}

func TestChain(t *testing.T) {
	user := Node("User").Named("u")
	company := Node("Company").Named("c")
	city := Node("City").Named("city")

	rel1 := user.RelationshipTo(company, "WORKS_AT")
	rel2 := company.RelationshipTo(city, "LOCATED_IN")

	chain := Chain(user, rel1, rel2)
	result := chain.String()

	// Should contain all nodes and relationships
	if !contains(result, "u") || !contains(result, "c") || !contains(result, "city") {
		t.Errorf("Chain() = %q, should contain all nodes", result)
	}
	if !contains(result, "WORKS_AT") || !contains(result, "LOCATED_IN") {
		t.Errorf("Chain() = %q, should contain all relationships", result)
	}
}

func TestChainEmpty(t *testing.T) {
	node := Node("Person").Named("p")
	chain := Chain(node)
	result := chain.String()

	// Should just be the node
	if result != node.String() {
		t.Errorf("Chain(node) = %q, want %q", result, node.String())
	}
}


