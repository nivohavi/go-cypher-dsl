package builder

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/expr"
)

func TestComplexQuery(t *testing.T) {
	person := ast.Node("Person").Named("p")
	movie := ast.Node("Movie").Named("m")
	actedIn := person.RelationshipTo(movie, "ACTED_IN")
	pattern := ast.Pattern(person, actedIn, movie)

	stmt, err := Match(pattern).
		Where(person.Property("name").Eq("Tom Hanks")).
		Returning(movie.Property("title"), movie.Property("released")).
		OrderBy(expr.Desc(movie.Property("released"))).
		Limit(10).
		Build()

	if err != nil {
		t.Fatalf("Complex query Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	required := []string{"MATCH", "WHERE", "RETURN", "ORDER BY", "DESC", "LIMIT"}
	for _, req := range required {
		if !strings.Contains(cypher, req) {
			t.Errorf("Cypher() = %q, should contain %q", cypher, req)
		}
	}
}

func TestQueryWithWith(t *testing.T) {
	person := ast.Node("Person").Named("p")
	movie := ast.Node("Movie").Named("m")
	actedIn := person.RelationshipTo(movie, "ACTED_IN")
	pattern := ast.Pattern(person, actedIn, movie)

	stmt, err := Match(person).
		Where(person.Property("name").Eq("Tom Hanks")).
		With(person).
		Match(pattern).
		Returning(movie).
		Build()

	if err != nil {
		t.Fatalf("Query with WITH Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "WITH") {
		t.Errorf("Cypher() = %q, should contain 'WITH'", cypher)
	}
}

func TestQueryWithAggregation(t *testing.T) {
	person := ast.Node("Person").Named("p")
	movie := ast.Node("Movie").Named("m")
	actedIn := person.RelationshipTo(movie, "ACTED_IN")
	pattern := ast.Pattern(person, actedIn, movie)

	movieCount := expr.As(expr.Count(movie), "movieCount")
	stmt, err := Match(pattern).
		Returning(
			person.Property("name"),
			movieCount,
		).
		OrderBy(expr.Desc(movieCount)).
		Limit(10).
		Build()

	if err != nil {
		t.Fatalf("Query with aggregation Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "count") {
		t.Errorf("Cypher() = %q, should contain 'count'", cypher)
	}
}

func TestQueryWithUnwind(t *testing.T) {
	names := core.NewParameter("", []string{"Tom Hanks", "Tom Cruise"})
	person := ast.Node("Person").Named("p")

	stmt, err := Unwind(names, "name").
		Match(person).
		Where(person.Property("name").Eq(expr.NewVariableExpression("name"))).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Query with UNWIND Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "UNWIND") {
		t.Errorf("Cypher() = %q, should contain 'UNWIND'", cypher)
	}
}

func TestQueryWithCreateAndReturn(t *testing.T) {
	person := ast.Node("Person").Named("p").WithProps(map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	})

	stmt, err := Create(person).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Query with CREATE Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "CREATE") {
		t.Errorf("Cypher() = %q, should contain 'CREATE'", cypher)
	}
	if !strings.Contains(cypher, "RETURN") {
		t.Errorf("Cypher() = %q, should contain 'RETURN'", cypher)
	}
}

func TestQueryWithMerge(t *testing.T) {
	person := ast.Node("Person").Named("p").WithProps(map[string]interface{}{
		"name": "John Doe",
	})

	stmt, err := Merge(person).
		OnCreate(expr.Equals(expr.Property("p", "created"), core.NewParameter("", 2023))).
		OnMatch(expr.Equals(expr.Property("p", "updated"), core.NewParameter("", 2023))).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Query with MERGE Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "MERGE") {
		t.Errorf("Cypher() = %q, should contain 'MERGE'", cypher)
	}
}

