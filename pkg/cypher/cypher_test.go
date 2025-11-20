package cypher

import (
	"strings"
	"testing"
)

func TestComplexPath(t *testing.T) {
	user := Node("User").Named("u")
	company := Node("Company").Named("c")
	city := Node("City").Named("city")

	path := ComplexPath(user, "WORKS_AT", company, "LOCATED_IN", city)
	result := path.String()

	if !strings.Contains(result, "WORKS_AT") || !strings.Contains(result, "LOCATED_IN") {
		t.Errorf("ComplexPath() = %q, should contain both relationship types", result)
	}
}

func TestCompareProperty(t *testing.T) {
	expr := CompareProperty("n", "age", ">", 30)
	result := expr.String()

	if !strings.Contains(result, "age") {
		t.Errorf("CompareProperty() = %q, should contain 'age'", result)
	}
}

func TestNamedCompareProperty(t *testing.T) {
	expr := NamedCompareProperty("n", "age", ">=", "minAge", 30)
	result := expr.String()

	if !strings.Contains(result, "age") {
		t.Errorf("NamedCompareProperty() = %q, should contain 'age'", result)
	}
}

func TestComparePropertyOperators(t *testing.T) {
	tests := []struct {
		op       string
		contains string
	}{
		{"=", "="},
		{">", ">"},
		{">=", ">="},
		{"<", "<"},
		{"<=", "<="},
		{"IN", "IN"},
		{"CONTAINS", "CONTAINS"},
		{"STARTS WITH", "STARTS WITH"},
		{"ENDS WITH", "ENDS WITH"},
	}

	for _, tt := range tests {
		t.Run(tt.op, func(t *testing.T) {
			if tt.op == "IN" {
				expr := CompareProperty("n", "status", tt.op, []string{"active", "pending"})
				result := expr.String()
				if !strings.Contains(result, tt.contains) {
					t.Errorf("CompareProperty() with %s = %q, should contain %q", tt.op, result, tt.contains)
				}
			} else if strings.Contains(tt.op, "CONTAINS") || strings.Contains(tt.op, "STARTS") || strings.Contains(tt.op, "ENDS") {
				expr := CompareProperty("n", "name", tt.op, "test")
				result := expr.String()
				if !strings.Contains(result, tt.contains) {
					t.Errorf("CompareProperty() with %s = %q, should contain %q", tt.op, result, tt.contains)
				}
			} else {
				expr := CompareProperty("n", "age", tt.op, 30)
				result := expr.String()
				if !strings.Contains(result, tt.contains) {
					t.Errorf("CompareProperty() with %s = %q, should contain %q", tt.op, result, tt.contains)
				}
			}
		})
	}
}

func TestPublicAPIFunctions(t *testing.T) {
	// Test various public API functions exist and work
	node := Node("Person").Named("p")
	if node == nil {
		t.Error("Node() returned nil")
	}

	prop := Property("p", "name")
	if prop == nil {
		t.Error("Property() returned nil")
	}

	param := Param("value")
	if param == nil {
		t.Error("Param() returned nil")
	}

	namedParam := NamedParam("name", "value")
	if namedParam == nil {
		t.Error("NamedParam() returned nil")
	}

	lit := Literal("test")
	if lit == nil {
		t.Error("Literal() returned nil")
	}
}

func TestPrettyPrint(t *testing.T) {
	node := Node("Person").Named("p")
	stmt, err := Match(node).
		Returning(node).
		Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	formatted := PrettyPrint(stmt)
	if !strings.Contains(formatted, "MATCH") {
		t.Errorf("PrettyPrint() = %q, should contain 'MATCH'", formatted)
	}
}

