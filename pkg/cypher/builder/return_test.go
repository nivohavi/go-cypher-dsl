package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestReturn(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Return(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Return().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

func TestReturnWithOrderBy(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Return(node).
		OrderBy(node.Property("name")).
		Desc()
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Return().OrderBy().Desc().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "ORDER BY") {
		t.Errorf("Cypher() = %q, should contain 'ORDER BY'", cypher)
	}
	if !strings.Contains(cypher, "DESC") {
		t.Errorf("Cypher() = %q, should contain 'DESC'", cypher)
	}
}

func TestReturnWithSkipLimit(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Return(node).
		Skip(10).
		Limit(5)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Return().Skip().Limit().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SKIP") {
		t.Errorf("Cypher() = %q, should contain 'SKIP'", cypher)
	}
	if !strings.Contains(cypher, "LIMIT") {
		t.Errorf("Cypher() = %q, should contain 'LIMIT'", cypher)
	}
}

func TestReturnMultipleExpressions(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Return(
		node,
		node.Property("name"),
		expr.Count(node),
	)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Return(...).Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

