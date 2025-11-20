package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

func TestDelete(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Delete(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Delete().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "DELETE") {
		t.Errorf("Cypher() = %q, should contain 'DELETE'", cypher)
	}
}

func TestDetachDelete(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := DetachDelete(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("DetachDelete().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "DETACH DELETE") {
		t.Errorf("Cypher() = %q, should contain 'DETACH DELETE'", cypher)
	}
}

func TestDeleteWithReturn(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Delete(node).
		Returning(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Delete().Returning().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

