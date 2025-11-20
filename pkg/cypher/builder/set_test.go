package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestSet(t *testing.T) {
	setExpr := expr.Equals(expr.Property("p", "age"), core.NewParameter("", 30))
	builder := Set(setExpr)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Set().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SET") {
		t.Errorf("Cypher() = %q, should contain 'SET'", cypher)
	}
}

func TestSetWithAnd(t *testing.T) {
	setExpr1 := expr.Equals(expr.Property("p", "age"), core.NewParameter("", 30))
	setExpr2 := expr.Equals(expr.Property("p", "active"), core.NewParameter("", true))
	builder := Set(setExpr1).
		And(setExpr2)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Set().And().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SET") {
		t.Errorf("Cypher() = %q, should contain 'SET'", cypher)
	}
}

func TestSetWithReturn(t *testing.T) {
	node := ast.Node("Person").Named("p")
	setExpr := expr.Equals(expr.Property("p", "age"), core.NewParameter("", 30))
	builder := Set(setExpr).
		Returning(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Set().Returning().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

