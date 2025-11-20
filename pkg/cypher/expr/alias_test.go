package expr

import (
	"testing"
)

func TestAs(t *testing.T) {
	expr := Property("n", "name")
	aliased := As(expr, "personName")
	result := aliased.String()
	// Should contain AS and the alias
	if !containsString(result, "AS") || !containsString(result, "personName") {
		t.Errorf("As(...).String() = %q, should contain 'AS' and 'personName'", result)
	}
}

func TestAsWithFunction(t *testing.T) {
	expr := Count(Property("n", "id"))
	aliased := As(expr, "total")
	result := aliased.String()
	// Should contain AS and the alias
	if !containsString(result, "AS") || !containsString(result, "total") {
		t.Errorf("As(Count(...)).String() = %q, should contain 'AS' and 'total'", result)
	}
}

func TestAsWithSpecialCharacters(t *testing.T) {
	expr := Property("n", "name")
	// Alias with spaces should be quoted
	aliased := As(expr, "person name")
	result := aliased.String()
	// Should contain backticks around the alias
	if !containsString(result, "`person name`") {
		t.Errorf("As(...) with spaces = %q, should contain backticks", result)
	}
}

func TestAsWithNumericStart(t *testing.T) {
	expr := Property("n", "value")
	aliased := As(expr, "2ndValue")
	result := aliased.String()
	// Should be quoted because it starts with a digit
	if !containsString(result, "`2ndValue`") {
		t.Errorf("As(...) with numeric start = %q, should contain backticks", result)
	}
}

func TestAliasExpressionLogicalOps(t *testing.T) {
	expr := Property("n", "age")
	aliased := As(expr, "age")
	other := Property("n", "score")

	// Test And
	andExpr := aliased.And(other)
	if andExpr == nil {
		t.Error("aliased.And(other) returned nil")
	}

	// Test Or
	orExpr := aliased.Or(other)
	if orExpr == nil {
		t.Error("aliased.Or(other) returned nil")
	}

	// Test Not
	notExpr := aliased.Not()
	if notExpr == nil {
		t.Error("aliased.Not() returned nil")
	}
}


