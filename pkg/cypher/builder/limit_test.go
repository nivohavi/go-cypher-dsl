package builder

import (
	"strings"
	"testing"
)

func TestLimit(t *testing.T) {
	builder := Limit(10)
	stmt, err := builder.Build()
	if err != nil {
		t.Fatalf("Limit().Build() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "LIMIT") {
		t.Errorf("Cypher() = %q, should contain 'LIMIT'", cypher)
	}
	if !strings.Contains(cypher, "10") {
		t.Errorf("Cypher() = %q, should contain '10'", cypher)
	}
}

