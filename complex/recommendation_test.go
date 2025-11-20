package complex

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

// Test 6: Collaborative Filtering - Product Recommendations
func TestCollaborativeFiltering(t *testing.T) {
	u1 := ast.Node("User").Named("u1")
	u2 := ast.Node("User").Named("u2")
	p := ast.Node("Product").Named("p")
	rec := ast.Node("Product").Named("rec")

	rel1 := u1.RelationshipTo(p, "PURCHASED")
	rel2 := u2.RelationshipTo(p, "PURCHASED")
	rel3 := u2.RelationshipTo(rec, "PURCHASED")

	stmt, err := cypher.Match(ast.Pattern(u1, rel1, p)).
		Match(ast.Pattern(u2, rel2, p)).
		Match(ast.Pattern(u2, rel3, rec)).
		Where(
			cypher.Not(cypher.Pattern(u1, u1.RelationshipTo(rec, "PURCHASED"), rec)),
		).
		Returning(
			cypher.As(rec.Property("name"), "RecommendedProduct"),
			cypher.As(cypher.Count(cypher.Var("u2")), "RecommendationScore"),
		).
		OrderBy(cypher.Desc(cypher.Var("RecommendationScore"))).
		Limit(5).
		Build()

	if err != nil {
		t.Fatalf("CollaborativeFiltering query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "PURCHASED") {
		t.Errorf("CollaborativeFiltering query should contain PURCHASED")
	}
	if !strings.Contains(cypherStr, "count") {
		t.Errorf("CollaborativeFiltering query should contain count")
	}
}

// Test 7: Viewing-Based Recommendations
func TestViewingBasedRecommendations(t *testing.T) {
	u1 := ast.Node("User").Named("u1")
	u2 := ast.Node("User").Named("u2")
	p := ast.Node("Product").Named("p")
	rec := ast.Node("Product").Named("rec")

	rel1 := u1.RelationshipTo(p, "VIEWED")
	rel2 := u2.RelationshipTo(p, "VIEWED")
	rel3 := u2.RelationshipTo(rec, "PURCHASED")

	stmt, err := cypher.Match(ast.Pattern(u1, rel1, p)).
		Match(ast.Pattern(u2, rel2, p)).
		Match(ast.Pattern(u2, rel3, rec)).
		Where(
			cypher.Not(cypher.Pattern(u1, u1.RelationshipTo(rec, "PURCHASED"), rec)),
		).
		Returning(
			cypher.As(rec.Property("name"), "RecommendedProduct"),
			cypher.As(cypher.Count(cypher.Var("u2")), "RecommendationScore"),
		).
		OrderBy(cypher.Desc(cypher.Var("RecommendationScore"))).
		Limit(5).
		Build()

	if err != nil {
		t.Fatalf("ViewingBasedRecommendations query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "VIEWED") {
		t.Errorf("ViewingBasedRecommendations query should contain VIEWED")
	}
}

// Test 8: Page Recommendations
func TestPageRecommendations(t *testing.T) {
	u := ast.Node("User").Named("u")
	other := ast.Node("User").Named("other")
	p := ast.Node("Page").Named("p")
	rec := ast.Node("Page").Named("rec")

	rel1 := u.RelationshipTo(p, "FOLLOW")
	rel2 := other.RelationshipTo(p, "FOLLOW")
	rel3 := other.RelationshipTo(rec, "FOLLOW")

	stmt, err := cypher.Match(ast.Pattern(u, rel1, p)).
		Match(ast.Pattern(other, rel2, p)).
		Match(ast.Pattern(other, rel3, rec)).
		Where(
			cypher.Not(cypher.Pattern(u, u.RelationshipTo(rec, "FOLLOW"), rec)),
		).
		Returning(
			cypher.As(rec.Property("name"), "RecommendedPage"),
			cypher.As(cypher.Count(cypher.Var("other")), "CommonFollowers"),
		).
		OrderBy(cypher.Desc(cypher.Var("CommonFollowers"))).
		Limit(10).
		Build()

	if err != nil {
		t.Fatalf("PageRecommendations query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "FOLLOW") {
		t.Errorf("PageRecommendations query should contain FOLLOW")
	}
}

