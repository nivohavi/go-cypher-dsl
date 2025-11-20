package expr

import (
	"testing"
)

func TestEquals(t *testing.T) {
	left := String("hello")
	right := String("world")
	expr := Equals(left, right)
	result := expr.String()
	// Should contain the equals operator
	if !containsString(result, "=") {
		t.Errorf("Equals(...).String() = %q, should contain '='", result)
	}
}

func TestNotEquals(t *testing.T) {
	left := Integer(10)
	right := Integer(20)
	expr := NotEquals(left, right)
	result := expr.String()
	expected := "(10 <> 20)"
	if result != expected {
		t.Errorf("NotEquals(...).String() = %q, want %q", result, expected)
	}
}

func TestGreaterThan(t *testing.T) {
	left := Integer(10)
	right := Integer(5)
	expr := GreaterThan(left, right)
	result := expr.String()
	expected := "(10 > 5)"
	if result != expected {
		t.Errorf("GreaterThan(...).String() = %q, want %q", result, expected)
	}
}

func TestLessThan(t *testing.T) {
	left := Integer(5)
	right := Integer(10)
	expr := LessThan(left, right)
	result := expr.String()
	expected := "(5 < 10)"
	if result != expected {
		t.Errorf("LessThan(...).String() = %q, want %q", result, expected)
	}
}

func TestGreaterThanEqual(t *testing.T) {
	left := Integer(10)
	right := Integer(10)
	expr := GreaterThanEqual(left, right)
	result := expr.String()
	expected := "(10 >= 10)"
	if result != expected {
		t.Errorf("GreaterThanEqual(...).String() = %q, want %q", result, expected)
	}
}

func TestLessThanEqual(t *testing.T) {
	left := Integer(5)
	right := Integer(10)
	expr := LessThanEqual(left, right)
	result := expr.String()
	expected := "(5 <= 10)"
	if result != expected {
		t.Errorf("LessThanEqual(...).String() = %q, want %q", result, expected)
	}
}

func TestIsNull(t *testing.T) {
	expr := Property("n", "name")
	nullExpr := IsNull(expr)
	result := nullExpr.String()
	// Should contain IS NULL
	if !containsString(result, "IS NULL") {
		t.Errorf("IsNull(...).String() = %q, should contain 'IS NULL'", result)
	}
}

func TestIsNotNull(t *testing.T) {
	expr := Property("n", "name")
	notNullExpr := IsNotNull(expr)
	result := notNullExpr.String()
	// Should contain IS NOT NULL
	if !containsString(result, "IS NOT NULL") {
		t.Errorf("IsNotNull(...).String() = %q, should contain 'IS NOT NULL'", result)
	}
}

func TestIn(t *testing.T) {
	expr := Property("n", "status")
	inExpr := In(expr, "active", "pending", "inactive")
	result := inExpr.String()
	// Should contain IN and the list
	if !containsString(result, "IN") {
		t.Errorf("In(...).String() = %q, should contain 'IN'", result)
	}
	if !containsString(result, "active") {
		t.Errorf("In(...).String() = %q, should contain 'active'", result)
	}
}

func TestContains(t *testing.T) {
	expr := Property("n", "name")
	containsExpr := Contains(expr, "test")
	result := containsExpr.String()
	// Should contain CONTAINS
	if !containsString(result, "CONTAINS") {
		t.Errorf("Contains(...).String() = %q, should contain 'CONTAINS'", result)
	}
}

func TestStartsWith(t *testing.T) {
	expr := Property("n", "name")
	startsWithExpr := StartsWith(expr, "test")
	result := startsWithExpr.String()
	// Should contain STARTS WITH
	if !containsString(result, "STARTS WITH") {
		t.Errorf("StartsWith(...).String() = %q, should contain 'STARTS WITH'", result)
	}
}

func TestEndsWith(t *testing.T) {
	expr := Property("n", "name")
	endsWithExpr := EndsWith(expr, "test")
	result := endsWithExpr.String()
	// Should contain ENDS WITH
	if !containsString(result, "ENDS WITH") {
		t.Errorf("EndsWith(...).String() = %q, should contain 'ENDS WITH'", result)
	}
}

func TestRegularExpression(t *testing.T) {
	expr := Property("n", "email")
	regexExpr := RegularExpression(expr, ".*@.*")
	result := regexExpr.String()
	// Should contain =~
	if !containsString(result, "=~") {
		t.Errorf("RegularExpression(...).String() = %q, should contain '=~'", result)
	}
}

func TestComparisonLogicalOps(t *testing.T) {
	left := Property("n", "age")
	right := Integer(30)
	comp := Equals(left, right)
	other := GreaterThan(left, Integer(20))

	// Test And
	andExpr := comp.And(other)
	if andExpr == nil {
		t.Error("comp.And(other) returned nil")
	}

	// Test Or
	orExpr := comp.Or(other)
	if orExpr == nil {
		t.Error("comp.Or(other) returned nil")
	}

	// Test Not
	notExpr := comp.Not()
	if notExpr == nil {
		t.Error("comp.Not() returned nil")
	}

	// Test Xor - only works on ComparisonExpression
	if compExpr, ok := comp.(*ComparisonExpression); ok {
		xorExpr := compExpr.Xor(other)
		if xorExpr == nil {
			t.Error("comp.Xor(other) returned nil")
		}
	}
}


