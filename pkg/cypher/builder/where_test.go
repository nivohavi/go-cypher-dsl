package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestWhere(t *testing.T) {
	condition := expr.Equals(expr.Property("n", "age"), core.NewParameter("", 30))
	builder := Where(condition)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Where().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "WHERE") {
		t.Errorf("Cypher() = %q, should contain 'WHERE'", cypher)
	}
}

func TestWhereWithAndWhere(t *testing.T) {
	condition1 := expr.Equals(expr.Property("n", "age"), core.NewParameter("", 30))
	condition2 := expr.Equals(expr.Property("n", "active"), core.NewParameter("", true))
	builder := Where(condition1).
		AndWhere(condition2)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Where().AndWhere().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "WHERE") {
		t.Errorf("Cypher() = %q, should contain 'WHERE'", cypher)
	}
}

func TestWhereWithOrWhere(t *testing.T) {
	condition1 := expr.Equals(expr.Property("n", "name"), core.NewParameter("", "John"))
	condition2 := expr.Equals(expr.Property("n", "name"), core.NewParameter("", "Jane"))
	builder := Where(condition1).
		OrWhere(condition2)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Where().OrWhere().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "WHERE") {
		t.Errorf("Cypher() = %q, should contain 'WHERE'", cypher)
	}
}

func TestWhereWithReturn(t *testing.T) {
	node := ast.Node("Person").Named("n")
	condition := node.Property("age").Gt(30)
	builder := Where(condition).
		Returning(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Where().Returning().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

