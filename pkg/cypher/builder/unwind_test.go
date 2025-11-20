package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestUnwind(t *testing.T) {
	list := core.NewParameter("", []string{"a", "b", "c"})
	builder := Unwind(list, "item")
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Unwind().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "UNWIND") {
		t.Errorf("Cypher() = %q, should contain 'UNWIND'", cypher)
	}
}

func TestUnwindWithMatch(t *testing.T) {
	list := core.NewParameter("", []string{"John", "Jane"})
	node := ast.Node("Person").Named("p")
	builder := Unwind(list, "name").
		Match(node).
		Where(node.Property("name").Eq(expr.NewVariableExpression("name")))
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Unwind().Match().Where().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "MATCH") {
		t.Errorf("Cypher() = %q, should contain 'MATCH'", cypher)
	}
}

func TestUnwindWithReturn(t *testing.T) {
	list := core.NewParameter("", []int{1, 2, 3})
	builder := Unwind(list, "num").
		Returning(expr.NewVariableExpression("num"))
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Unwind().Returning().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

