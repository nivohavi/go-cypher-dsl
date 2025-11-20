package complex

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

// Test 1: Mutual Friends Query
func TestMutualFriends(t *testing.T) {
	alice := ast.Node("Person").Named("a").WithProps(map[string]interface{}{"name": "Alice"})
	bob := ast.Node("Person").Named("b").WithProps(map[string]interface{}{"name": "Bob"})
	mf := ast.Node("Person").Named("mf")

	rel1 := alice.RelationshipBetween(mf, "FRIENDS_WITH")
	rel2 := mf.RelationshipBetween(bob, "FRIENDS_WITH")

	stmt, err := cypher.Match(ast.Pattern(alice, rel1, mf, rel2, bob)).
		Returning(cypher.As(mf.Property("name"), "MutualFriend")).
		Build()

	if err != nil {
		t.Fatalf("MutualFriends query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "FRIENDS_WITH") {
		t.Errorf("MutualFriends query should contain FRIENDS_WITH")
	}
	if !strings.Contains(cypherStr, "RETURN") {
		t.Errorf("MutualFriends query should contain RETURN")
	}
}

// Test 2: Friend Recommendations
func TestFriendRecommendations(t *testing.T) {
	user := ast.Node("User").Named("user")
	friend := ast.Node("User").Named("friend")
	recommended := ast.Node("User").Named("recommended")

	rel1 := user.RelationshipTo(friend, "FOLLOWS")
	rel2 := friend.RelationshipTo(recommended, "FOLLOWS")

	stmt, err := cypher.Match(ast.Pattern(user, rel1, friend, rel2, recommended)).
		Where(
			cypher.And(
				cypher.Not(cypher.Pattern(user, user.RelationshipTo(recommended, "FOLLOWS"), recommended)),
				cypher.Ne(cypher.Var("user"), cypher.Var("recommended")),
			),
		).
		Returning(
			cypher.As(recommended.Property("name"), "RecommendedUser"),
			cypher.As(cypher.Count(cypher.Var("friend")), "mutualFriends"),
		).
		OrderBy(cypher.Desc(cypher.Var("mutualFriends"))).
		Limit(10).
		Build()

	if err != nil {
		t.Fatalf("FriendRecommendations query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "FOLLOWS") {
		t.Errorf("FriendRecommendations query should contain FOLLOWS")
	}
	if !strings.Contains(cypherStr, "ORDER BY") {
		t.Errorf("FriendRecommendations query should contain ORDER BY")
	}
}

// Test 3: Friends of Friends
func TestFriendsOfFriends(t *testing.T) {
	p := ast.Node("Person").Named("p").WithProps(map[string]interface{}{"name": "Alice"})
	fof := ast.Node("Person").Named("fof")

	rel1 := p.RelationshipBetween(ast.Node("Person").Named("_"), "FRIENDS_WITH")
	rel2 := ast.Node("Person").Named("_").RelationshipBetween(fof, "FRIENDS_WITH")

	stmt, err := cypher.Match(ast.Pattern(p, rel1, ast.Node("Person").Named("_"), rel2, fof)).
		Where(
			cypher.And(
				cypher.Not(cypher.Pattern(p, p.RelationshipBetween(fof, "FRIENDS_WITH"), fof)),
				cypher.Ne(cypher.Var("p"), cypher.Var("fof")),
			),
		).
		Returning(cypher.Distinct(cypher.As(fof.Property("name"), "FriendOfFriend"))).
		Build()

	if err != nil {
		t.Fatalf("FriendsOfFriends query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "FRIENDS_WITH") {
		t.Errorf("FriendsOfFriends query should contain FRIENDS_WITH")
	}
}

// Test 4: Network Analysis - Count Friends
func TestNetworkAnalysis(t *testing.T) {
	p := ast.Node("Person").Named("p")
	rel := p.RelationshipBetween(ast.Node("Person").Named("_"), "FRIENDS_WITH")

	stmt, err := cypher.Match(ast.Pattern(p, rel, ast.Node("Person").Named("_"))).
		Returning(
			cypher.As(p.Property("name"), "Person"),
			cypher.As(cypher.Count(rel), "numberOfFriends"),
		).
		OrderBy(cypher.Desc(cypher.Var("numberOfFriends"))).
		Build()

	if err != nil {
		t.Fatalf("NetworkAnalysis query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "count") {
		t.Errorf("NetworkAnalysis query should contain count")
	}
}

