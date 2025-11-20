package complex

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

// Test 9: Suspicious IP Patterns
func TestSuspiciousIPPatterns(t *testing.T) {
	txn := ast.Node("Transaction").Named("t")
	account := ast.Node("Account").Named("a")
	ip := ast.Node("IPAddress").Named("ip")

	rel1 := txn.RelationshipTo(account, "FROM")
	rel2 := account.RelationshipTo(ip, "HAS_IP")

	stmt, err := cypher.Match(ast.Pattern(txn, rel1, account, rel2, ip)).
		With(ip, cypher.As(cypher.Count(cypher.Distinct(account)), "accountCount")).
		Where(
			cypher.Gt(cypher.Var("accountCount"), cypher.Param(1)),
		).
		Returning(
			cypher.As(ip.Property("address"), "SuspiciousIP"),
			cypher.Var("accountCount"),
		).
		OrderBy(cypher.Desc(cypher.Var("accountCount"))).
		Build()

	if err != nil {
		t.Fatalf("SuspiciousIPPatterns query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "HAS_IP") {
		t.Errorf("SuspiciousIPPatterns query should contain HAS_IP")
	}
	if !strings.Contains(cypherStr, "WITH") {
		t.Errorf("SuspiciousIPPatterns query should contain WITH")
	}
}

// Test 10: Multiple Accounts Same IP
func TestMultipleAccountsSameIP(t *testing.T) {
	a1 := ast.Node("Account").Named("a1")
	a2 := ast.Node("Account").Named("a2")
	ip := ast.Node("IPAddress").Named("ip")

	rel1 := a1.RelationshipTo(ip, "HAS_IP")
	rel2 := a2.RelationshipTo(ip, "HAS_IP")

	stmt, err := cypher.Match(ast.Pattern(a1, rel1, ip)).
		Match(ast.Pattern(a2, rel2, ip)).
		Where(
			cypher.Ne(cypher.Var("a1"), cypher.Var("a2")),
		).
		With(ip, cypher.As(cypher.Collect(cypher.Distinct(cypher.Var("a1"))), "accounts")).
		Where(
			cypher.Gt(cypher.Function("size", cypher.Var("accounts")), cypher.Param(2)),
		).
		Returning(
			ip.Property("address"),
			cypher.As(cypher.Function("size", cypher.Var("accounts")), "accountCount"),
		).
		Build()

	if err != nil {
		t.Fatalf("MultipleAccountsSameIP query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "HAS_IP") {
		t.Errorf("MultipleAccountsSameIP query should contain HAS_IP")
	}
}

