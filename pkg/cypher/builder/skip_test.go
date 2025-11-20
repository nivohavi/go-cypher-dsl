package builder

import (
	"strings"
	"testing"
)

func TestSkip(t *testing.T) {
	builder := Skip(10)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Skip().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SKIP") {
		t.Errorf("Cypher() = %q, should contain 'SKIP'", cypher)
	}
	if !strings.Contains(cypher, "10") {
		t.Errorf("Cypher() = %q, should contain '10'", cypher)
	}
}

func TestSkipWithLimit(t *testing.T) {
	builder := Skip(10).
		Limit(5)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Skip().Limit().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SKIP") {
		t.Errorf("Cypher() = %q, should contain 'SKIP'", cypher)
	}
	if !strings.Contains(cypher, "LIMIT") {
		t.Errorf("Cypher() = %q, should contain 'LIMIT'", cypher)
	}
}

