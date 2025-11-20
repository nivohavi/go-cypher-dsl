package cypher

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

func TestEmptyReturn(t *testing.T) {
	node := ast.Node("Person").Named("p")
	builder := Match(node).
		Returning()
	stmt, err := builder.Build()

	// This should either error or handle gracefully
	if err == nil {
		cypher := stmt.Cypher()
		// If it doesn't error, the query should still be valid
		if !strings.Contains(cypher, "MATCH") {
			t.Errorf("EmptyReturn query should contain MATCH")
		}
	}
}

func TestNullHandling(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Where(node.Property("name").IsNull()).
		Returning(node).
		Build()

	if err != nil {
		t.Fatalf("NullHandling query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "IS NULL") {
		t.Errorf("NullHandling query should contain 'IS NULL'")
	}
}

func TestNotNullHandling(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Where(node.Property("name").IsNotNull()).
		Returning(node).
		Build()

	if err != nil {
		t.Fatalf("NotNullHandling query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	// Should contain either "IS NOT NULL" or "NOT ... IS NULL"
	if !strings.Contains(cypher, "IS NOT NULL") && !strings.Contains(cypher, "NOT") {
		t.Errorf("NotNullHandling query should contain 'IS NOT NULL' or 'NOT'")
	}
}

func TestComplexNestedConditions(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Where(
			And(
				Or(
					node.Property("age").Gt(18),
					node.Property("age").Lt(65),
				),
				Not(node.Property("deleted").Eq(true)),
			),
		).
		Returning(node).
		Build()

	if err != nil {
		t.Fatalf("ComplexNestedConditions query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "WHERE") {
		t.Errorf("ComplexNestedConditions query should contain WHERE")
	}
}

func TestMultipleLabels(t *testing.T) {
	node := ast.Node("Person", "Actor").Named("p")
	stmt, err := Match(node).
		Returning(node).
		Build()

	if err != nil {
		t.Fatalf("MultipleLabels query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "Person") || !strings.Contains(cypher, "Actor") {
		t.Errorf("MultipleLabels query should contain both labels")
	}
}

func TestParameterSubstitution(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Where(node.Property("name").Eq(Param("John"))).
		Returning(node).
		Build()

	if err != nil {
		t.Fatalf("ParameterSubstitution query Build() error = %v", err)
	}

	params := stmt.Params()
	if len(params) == 0 {
		t.Error("ParameterSubstitution query should have parameters")
	}
}

func TestNamedParameterSubstitution(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Where(node.Property("name").Eq(NamedParam("personName", "John"))).
		Returning(node).
		Build()

	if err != nil {
		t.Fatalf("NamedParameterSubstitution query Build() error = %v", err)
	}

	params := stmt.Params()
	if params["personName"] != "John" {
		t.Errorf("NamedParameterSubstitution params = %v, want personName='John'", params)
	}
}

func TestInOperator(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Where(node.Property("status").In("active", "pending", "inactive")).
		Returning(node).
		Build()

	if err != nil {
		t.Fatalf("InOperator query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "IN") {
		t.Errorf("InOperator query should contain 'IN'")
	}
}

func TestStringOperations(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Where(
			Or(
				node.Property("name").StartsWith("Jo"),
				node.Property("email").Contains("@example.com"),
			),
		).
		Returning(node).
		Build()

	if err != nil {
		t.Fatalf("StringOperations query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "STARTS WITH") && !strings.Contains(cypher, "CONTAINS") {
		t.Errorf("StringOperations query should contain string operations")
	}
}

func TestAggregationFunctions(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Returning(
			As(Count(node), "total"),
			As(Sum(node.Property("age")), "totalAge"),
			As(Avg(node.Property("age")), "avgAge"),
			As(Min(node.Property("age")), "minAge"),
			As(Max(node.Property("age")), "maxAge"),
		).
		Build()

	if err != nil {
		t.Fatalf("AggregationFunctions query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	required := []string{"count", "sum", "avg", "min", "max"}
	for _, req := range required {
		if !strings.Contains(strings.ToLower(cypher), req) {
			t.Errorf("AggregationFunctions query should contain %s", req)
		}
	}
}

func TestDistinct(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Returning(Distinct(node.Property("name"))).
		Build()

	if err != nil {
		t.Fatalf("Distinct query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "DISTINCT") {
		t.Errorf("Distinct query should contain 'DISTINCT'")
	}
}

func TestPagination(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := Match(node).
		Returning(node).
		OrderBy(node.Property("name")).
		Skip(20).
		Limit(10).
		Build()

	if err != nil {
		t.Fatalf("Pagination query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SKIP") || !strings.Contains(cypher, "LIMIT") {
		t.Errorf("Pagination query should contain SKIP and LIMIT")
	}
}

func TestOptionalMatch(t *testing.T) {
	node1 := ast.Node("Person").Named("p")
	node2 := ast.Node("Address").Named("a")
	rel := node1.RelationshipTo(node2, "HAS_ADDRESS")

	stmt, err := Match(node1).
		OptionalMatch(ast.Pattern(node1, rel, node2)).
		Returning(node1.Property("name"), node2.Property("city")).
		Build()

	if err != nil {
		t.Fatalf("OptionalMatch query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "OPTIONAL MATCH") {
		t.Errorf("OptionalMatch query should contain 'OPTIONAL MATCH'")
	}
}

