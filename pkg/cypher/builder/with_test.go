package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

func TestWith(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := With(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("With().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "WITH") {
		t.Errorf("Cypher() = %q, should contain 'WITH'", cypher)
	}
}

func TestWithWithWhere(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := With(node).
		Where(node.Property("age").Gt(30))
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("With().Where().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "WHERE") {
		t.Errorf("Cypher() = %q, should contain 'WHERE'", cypher)
	}
}

func TestWithWithOrderBy(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := With(node).
		OrderBy(node.Property("name")).
		Desc()
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("With().OrderBy().Desc().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "ORDER BY") {
		t.Errorf("Cypher() = %q, should contain 'ORDER BY'", cypher)
	}
}

func TestWithWithSkipLimit(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := With(node).
		Skip(10).
		Limit(5)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("With().Skip().Limit().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SKIP") {
		t.Errorf("Cypher() = %q, should contain 'SKIP'", cypher)
	}
	if !strings.Contains(cypher, "LIMIT") {
		t.Errorf("Cypher() = %q, should contain 'LIMIT'", cypher)
	}
}

func TestWithChaining(t *testing.T) {
	node1 := ast.Node("Person").Named("p")
	node2 := ast.Node("Movie").Named("m")
	builder := With(node1).
		With(node2)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("With().With().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	withs := strings.Count(cypher, "WITH")
	if withs < 2 {
		t.Errorf("Cypher() should contain at least 2 WITH clauses, got %d", withs)
	}
}

func TestWithWithMatch(t *testing.T) {
	node1 := ast.Node("Person").Named("p")
	node2 := ast.Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN")
	pattern := ast.Pattern(node1, rel, node2)

	builder := With(node1).
		Match(pattern)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("With().Match().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "MATCH") {
		t.Errorf("Cypher() = %q, should contain 'MATCH'", cypher)
	}
}

