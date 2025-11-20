package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestMerge(t *testing.T) {
	node := ast.Node("Person").Named("p").WithProps(map[string]interface{}{
		"name": "John",
	})
	builder := Merge(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Merge().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "MERGE") {
		t.Errorf("Cypher() = %q, should contain 'MERGE'", cypher)
	}
}

func TestMergeWithOnCreate(t *testing.T) {
	node := ast.Node("Person").Named("p").WithProps(map[string]interface{}{
		"name": "John",
	})
	builder := Merge(node).
		OnCreate(expr.Equals(expr.Property("p", "created"), core.NewParameter("", 2023)))
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Merge().OnCreate().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "ON CREATE SET") {
		t.Errorf("Cypher() = %q, should contain 'ON CREATE SET'", cypher)
	}
}

func TestMergeWithOnMatch(t *testing.T) {
	node := ast.Node("Person").Named("p").WithProps(map[string]interface{}{
		"name": "John",
	})
	builder := Merge(node).
		OnMatch(expr.Equals(expr.Property("p", "updated"), core.NewParameter("", 2023)))
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Merge().OnMatch().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "ON MATCH SET") {
		t.Errorf("Cypher() = %q, should contain 'ON MATCH SET'", cypher)
	}
}

func TestMergeWithBothOnCreateAndOnMatch(t *testing.T) {
	node := ast.Node("Person").Named("p").WithProps(map[string]interface{}{
		"name": "John",
	})
	builder := Merge(node).
		OnCreate(expr.Equals(expr.Property("p", "created"), core.NewParameter("", 2023))).
		OnMatch(expr.Equals(expr.Property("p", "updated"), core.NewParameter("", 2023)))
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Merge().OnCreate().OnMatch().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "ON CREATE SET") {
		t.Errorf("Cypher() = %q, should contain 'ON CREATE SET'", cypher)
	}
	if !strings.Contains(cypher, "ON MATCH SET") {
		t.Errorf("Cypher() = %q, should contain 'ON MATCH SET'", cypher)
	}
}

func TestMergeWithReturn(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Merge(node).
		Returning(node)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Merge().Returning().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

