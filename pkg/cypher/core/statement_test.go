package core

import (
	"testing"
)

func TestNewStatement(t *testing.T) {
	cypher := "MATCH (n) RETURN n"
	params := map[string]any{"key": "value"}
	stmt := NewStatement(cypher, params)

	if stmt.Cypher() != cypher {
		t.Errorf("Cypher() = %q, want %q", stmt.Cypher(), cypher)
	}

	resultParams := stmt.Params()
	if resultParams["key"] != "value" {
		t.Errorf("Params() = %v, want map with key 'value'", resultParams)
	}
}

func TestNewStatementWithNilParams(t *testing.T) {
	cypher := "MATCH (n) RETURN n"
	stmt := NewStatement(cypher, nil)

	if stmt.Cypher() != cypher {
		t.Errorf("Cypher() = %q, want %q", stmt.Cypher(), cypher)
	}

	resultParams := stmt.Params()
	if resultParams == nil {
		t.Error("Params() should not be nil")
	}
	if len(resultParams) != 0 {
		t.Errorf("Params() should be empty, got %v", resultParams)
	}
}

func TestStatementWithCypher(t *testing.T) {
	stmt := NewStatement("MATCH (n) RETURN n", nil)
	newStmt := stmt.WithCypher("MATCH (m) RETURN m")

	if newStmt.Cypher() != "MATCH (m) RETURN m" {
		t.Errorf("WithCypher() = %q, want 'MATCH (m) RETURN m'", newStmt.Cypher())
	}
}

func TestStatementWithParams(t *testing.T) {
	stmt := NewStatement("MATCH (n) RETURN n", map[string]any{"old": "value"})
	newParams := map[string]any{"new": "value"}
	newStmt := stmt.WithParams(newParams)

	if newStmt.Params()["new"] != "value" {
		t.Errorf("WithParams() did not set new params correctly")
	}
}

func TestStatementMerge(t *testing.T) {
	stmt1 := NewStatement("MATCH (n)", map[string]any{"key1": "value1"})
	stmt2 := NewStatement("RETURN n", map[string]any{"key2": "value2"})

	merged := stmt1.Merge(stmt2)
	cypher := merged.Cypher()

	if !contains(cypher, "MATCH") || !contains(cypher, "RETURN") {
		t.Errorf("Merge() cypher = %q, should contain both MATCH and RETURN", cypher)
	}

	params := merged.Params()
	if params["key1"] != "value1" || params["key2"] != "value2" {
		t.Errorf("Merge() params = %v, should contain both keys", params)
	}
}

func TestStatementMergeWithNil(t *testing.T) {
	stmt := NewStatement("MATCH (n)", nil)
	merged := stmt.Merge(nil)

	if merged.Cypher() != stmt.Cypher() {
		t.Errorf("Merge(nil) should return original statement")
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

