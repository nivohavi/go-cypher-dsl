package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestRemove(t *testing.T) {
	removeExpr := expr.Property("p", "age")
	builder := Remove(removeExpr)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Remove().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "REMOVE") {
		t.Errorf("Cypher() = %q, should contain 'REMOVE'", cypher)
	}
}

func TestRemoveWithAnd(t *testing.T) {
	removeExpr1 := expr.Property("p", "age")
	removeExpr2 := expr.Property("p", "temp")
	builder := Remove(removeExpr1).
		And(removeExpr2)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Remove().And().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "REMOVE") {
		t.Errorf("Cypher() = %q, should contain 'REMOVE'", cypher)
	}
}

func TestRemoveWithReturn(t *testing.T) {
	node := ast.Node("Person").Named("p")
	removeExpr := expr.Property("p", "age")
	builder := Remove(removeExpr).
		Returning(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Remove().Returning().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

