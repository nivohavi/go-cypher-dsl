package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

func TestCreate(t *testing.T) {
	node := ast.Node("Person").Named("p").WithProps(map[string]interface{}{
		"name": "John",
		"age":  30,
	})
	builder := Create(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Create().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "CREATE") {
		t.Errorf("Cypher() = %q, should contain 'CREATE'", cypher)
	}
}

func TestCreateWithReturn(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Create(node).
		Returning(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Create().Returning().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

func TestMultipleCreates(t *testing.T) {
	node1 := ast.Node("Person").Named("p1")
	node2 := ast.Node("Person").Named("p2")
	builder := Create(node1).
		Create(node2)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Create().Create().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	creates := strings.Count(cypher, "CREATE")
	if creates < 2 {
		t.Errorf("Cypher() should contain at least 2 CREATE clauses, got %d", creates)
	}
}

func TestCreateRelationship(t *testing.T) {
	node1 := ast.Node("Person").Named("p")
	node2 := ast.Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN")
	pattern := ast.Pattern(node1, rel, node2)

	builder := Create(pattern)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Create(pattern).Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "ACTED_IN") {
		t.Errorf("Cypher() = %q, should contain 'ACTED_IN'", cypher)
	}
}

