package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestMatch(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Match(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "MATCH") {
		t.Errorf("Cypher() = %q, should contain 'MATCH'", cypher)
	}
	if !strings.Contains(cypher, "Person") {
		t.Errorf("Cypher() = %q, should contain 'Person'", cypher)
	}
}

func TestMatchWithWhere(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Match(node).
		Where(node.Property("name").Eq("John"))
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().Where().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "WHERE") {
		t.Errorf("Cypher() = %q, should contain 'WHERE'", cypher)
	}
}

func TestMatchWithReturn(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Match(node).
		Returning(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().Returning().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

func TestOptionalMatch(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := OptionalMatch(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("OptionalMatch().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "OPTIONAL MATCH") {
		t.Errorf("Cypher() = %q, should contain 'OPTIONAL MATCH'", cypher)
	}
}

func TestMultipleMatches(t *testing.T) {
	node1 := ast.Node("Person").Named("p")
	node2 := ast.Node("Movie").Named("m")
	builder := Match(node1).
		Match(node2)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().Match().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	// Should contain multiple MATCH clauses
	matches := strings.Count(cypher, "MATCH")
	if matches < 2 {
		t.Errorf("Cypher() should contain at least 2 MATCH clauses, got %d", matches)
	}
}

func TestMatchWithCreate(t *testing.T) {
	node1 := ast.Node("Person").Named("p")
	node2 := ast.Node("Movie").Named("m")
	rel := node1.RelationshipTo(node2, "ACTED_IN")
	pattern := ast.Pattern(node1, rel, node2)

	builder := Match(node1).
		Create(pattern)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().Create().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "CREATE") {
		t.Errorf("Cypher() = %q, should contain 'CREATE'", cypher)
	}
}

func TestMatchWithMerge(t *testing.T) {
	node := ast.Node("Person").Named("p").WithProps(map[string]interface{}{
		"name": "John",
	})

	builder := Match(node).
		Merge(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().Merge().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "MERGE") {
		t.Errorf("Cypher() = %q, should contain 'MERGE'", cypher)
	}
}

func TestMatchWithDelete(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Match(node).
		Delete(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().Delete().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "DELETE") {
		t.Errorf("Cypher() = %q, should contain 'DELETE'", cypher)
	}
}

func TestMatchWithDetachDelete(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Match(node).
		DetachDelete(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().DetachDelete().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "DETACH DELETE") {
		t.Errorf("Cypher() = %q, should contain 'DETACH DELETE'", cypher)
	}
}

func TestMatchWithSet(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Match(node).
		Set(expr.Equals(expr.Property("p", "age"), core.NewParameter("", 30)))
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().Set().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SET") {
		t.Errorf("Cypher() = %q, should contain 'SET'", cypher)
	}
}

func TestMatchWithWith(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Match(node).
		With(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Match().With().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "WITH") {
		t.Errorf("Cypher() = %q, should contain 'WITH'", cypher)
	}
}

