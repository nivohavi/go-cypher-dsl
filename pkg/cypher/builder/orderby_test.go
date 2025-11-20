package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

func TestOrderBy(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := OrderBy(node.Property("name"))
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("OrderBy().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "ORDER BY") {
		t.Errorf("Cypher() = %q, should contain 'ORDER BY'", cypher)
	}
}

func TestOrderByAsc(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := OrderBy(node.Property("name")).
		Asc()
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("OrderBy().Asc().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "ASC") {
		t.Errorf("Cypher() = %q, should contain 'ASC'", cypher)
	}
}

func TestOrderByDesc(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := OrderBy(node.Property("name")).
		Desc()
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("OrderBy().Desc().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "DESC") {
		t.Errorf("Cypher() = %q, should contain 'DESC'", cypher)
	}
}

func TestOrderByWithSkipLimit(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := OrderBy(node.Property("name")).
		Skip(10).
		Limit(5)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("OrderBy().Skip().Limit().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SKIP") {
		t.Errorf("Cypher() = %q, should contain 'SKIP'", cypher)
	}
	if !strings.Contains(cypher, "LIMIT") {
		t.Errorf("Cypher() = %q, should contain 'LIMIT'", cypher)
	}
}

func TestOrderByMultipleExpressions(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := OrderBy(
		node.Property("lastName"),
		node.Property("firstName"),
	)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("OrderBy(...).Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "ORDER BY") {
		t.Errorf("Cypher() = %q, should contain 'ORDER BY'", cypher)
	}
}

