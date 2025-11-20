package complex

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

// Test 15: Co-Actors
func TestCoActors(t *testing.T) {
	actor1 := ast.Node("Person").Named("actor1").WithProps(map[string]interface{}{"name": "Tom Hanks"})
	actor2 := ast.Node("Person").Named("actor2")
	movie := ast.Node("Movie").Named("m")

	rel1 := actor1.RelationshipTo(movie, "ACTED_IN")
	rel2 := actor2.RelationshipTo(movie, "ACTED_IN")

	stmt, err := cypher.Match(ast.Pattern(actor1, rel1, movie)).
		Match(ast.Pattern(actor2, rel2, movie)).
		Where(
			cypher.Ne(cypher.Var("actor1"), cypher.Var("actor2")),
		).
		Returning(
			cypher.As(actor2.Property("name"), "CoActor"),
			cypher.As(cypher.Count(movie), "MoviesTogether"),
		).
		OrderBy(cypher.Desc(cypher.Var("MoviesTogether"))).
		Build()

	if err != nil {
		t.Fatalf("CoActors query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "ACTED_IN") {
		t.Errorf("CoActors query should contain ACTED_IN")
	}
	if !strings.Contains(cypherStr, "count") {
		t.Errorf("CoActors query should contain count")
	}
}

// Test 16: Actor-Director Network
func TestActorDirectorNetwork(t *testing.T) {
	actor := ast.Node("Person").Named("actor")
	director := ast.Node("Person").Named("director")
	movie := ast.Node("Movie").Named("m")

	rel1 := actor.RelationshipTo(movie, "ACTED_IN")
	rel2 := movie.RelationshipTo(director, "DIRECTED")

	stmt, err := cypher.Match(ast.Pattern(actor, rel1, movie, rel2, director)).
		Returning(
			cypher.As(actor.Property("name"), "Actor"),
			cypher.As(director.Property("name"), "Director"),
			cypher.As(cypher.Count(movie), "collaborations"),
		).
		OrderBy(cypher.Desc(cypher.Var("collaborations"))).
		Limit(20).
		Build()

	if err != nil {
		t.Fatalf("ActorDirectorNetwork query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "ACTED_IN") {
		t.Errorf("ActorDirectorNetwork query should contain ACTED_IN")
	}
	if !strings.Contains(cypherStr, "DIRECTED") {
		t.Errorf("ActorDirectorNetwork query should contain DIRECTED")
	}
}

// Test 17: Movie Recommendations
func TestMovieRecommendations(t *testing.T) {
	user := ast.Node("User").Named("user")
	other := ast.Node("User").Named("other")
	m1 := ast.Node("Movie").Named("m1")
	m2 := ast.Node("Movie").Named("m2")

	rel1 := user.RelationshipTo(m1, "RATED").WithProps(map[string]interface{}{"rating": 5})
	rel2 := other.RelationshipTo(m1, "RATED").WithProps(map[string]interface{}{"rating": 5})
	rel3 := other.RelationshipTo(m2, "RATED").WithProps(map[string]interface{}{"rating": 5})

	stmt, err := cypher.Match(ast.Pattern(user, rel1, m1)).
		Match(ast.Pattern(other, rel2, m1)).
		Match(ast.Pattern(other, rel3, m2)).
		Where(
			cypher.Not(cypher.Pattern(user, user.RelationshipTo(m2, "RATED"), m2)),
		).
		Returning(
			cypher.As(m2.Property("title"), "RecommendedMovie"),
			cypher.As(cypher.Count(cypher.Var("other")), "recommendationScore"),
		).
		OrderBy(cypher.Desc(cypher.Var("recommendationScore"))).
		Limit(10).
		Build()

	if err != nil {
		t.Fatalf("MovieRecommendations query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "RATED") {
		t.Errorf("MovieRecommendations query should contain RATED")
	}
}

