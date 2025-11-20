package driver

import (
	"context"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// Note: These tests require a running Neo4j instance or mocks
// They are structured to test the SessionManager and QueryHelper logic
// In a real scenario, you would use a test Neo4j instance or mocks

func TestNewSessionManager(t *testing.T) {
	// This test would require a Neo4j driver
	// For now, we just test that the function exists and can be called
	// In a real test environment, you would create a mock driver
	t.Skip("Requires Neo4j driver - skipping in unit tests")
}

func TestNewQueryHelper(t *testing.T) {
	helper := NewQueryHelper()
	if helper == nil {
		t.Error("NewQueryHelper() returned nil")
	}
}

func TestQueryHelperCollectSingle(t *testing.T) {
	helper := NewQueryHelper()
	handler := helper.CollectSingle("field")
	if handler == nil {
		t.Error("CollectSingle() returned nil handler")
	}
}

func TestQueryHelperCollectList(t *testing.T) {
	helper := NewQueryHelper()
	handler := helper.CollectList("field")
	if handler == nil {
		t.Error("CollectList() returned nil handler")
	}
}

func TestQueryHelperCollectMap(t *testing.T) {
	helper := NewQueryHelper()
	handler := helper.CollectMap("keyField", "valueField")
	if handler == nil {
		t.Error("CollectMap() returned nil handler")
	}
}

func TestQueryHelperCollectAll(t *testing.T) {
	helper := NewQueryHelper()
	handler := helper.CollectAll()
	if handler == nil {
		t.Error("CollectAll() returned nil handler")
	}
}

func TestQueryHelperCountResults(t *testing.T) {
	helper := NewQueryHelper()
	handler := helper.CountResults()
	if handler == nil {
		t.Error("CountResults() returned nil handler")
	}
}

// Test that queries can be built for use with SessionManager
func TestQueryBuildingForSessionManager(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := cypher.Match(node).
		Returning(node).
		Limit(5).
		Build()

	if err != nil {
		t.Fatalf("Query building for SessionManager error = %v", err)
	}

	// Verify the statement has Cypher and Params
	if stmt.Cypher() == "" {
		t.Error("Statement Cypher() should not be empty")
	}
	if stmt.Params() == nil {
		t.Error("Statement Params() should not be nil")
	}
}

// Test batch query building
func TestBatchQueryBuilding(t *testing.T) {
	node1 := ast.Node("Person").Named("p1").WithProps(map[string]interface{}{
		"name": "Alice",
	})
	node2 := ast.Node("Person").Named("p2").WithProps(map[string]interface{}{
		"name": "Bob",
	})

	stmt1, err1 := cypher.Create(node1).Build()
	stmt2, err2 := cypher.Create(node2).Build()

	if err1 != nil {
		t.Fatalf("Statement 1 build error = %v", err1)
	}
	if err2 != nil {
		t.Fatalf("Statement 2 build error = %v", err2)
	}

	statements := []core.Statement{stmt1, stmt2}
	if len(statements) != 2 {
		t.Errorf("Batch statements length = %d, want 2", len(statements))
	}
}

// Test query with parameters for SessionManager
func TestQueryWithParameters(t *testing.T) {
	node := ast.Node("Person").Named("p")
	stmt, err := cypher.Match(node).
		Where(node.Property("name").Eq(cypher.Param("John"))).
		Returning(node).
		Build()

	if err != nil {
		t.Fatalf("Query with parameters error = %v", err)
	}

	params := stmt.Params()
	if len(params) == 0 {
		t.Error("Query with parameters should have parameters")
	}
}

// Test context usage (even though we can't actually execute)
func TestContextUsage(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Error("context.Background() returned nil")
	}

	// Verify context can be used (it's just for demonstration)
	_ = ctx
}

