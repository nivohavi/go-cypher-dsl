package expr

import (
	"testing"
)

func TestAsc(t *testing.T) {
	expr := Property("n", "name")
	ascExpr := Asc(expr)
	result := ascExpr.String()
	// Should contain ASC
	if !containsString(result, "ASC") {
		t.Errorf("Asc(...).String() = %q, should contain 'ASC'", result)
	}
}

func TestDesc(t *testing.T) {
	expr := Property("n", "age")
	descExpr := Desc(expr)
	result := descExpr.String()
	// Should contain DESC
	if !containsString(result, "DESC") {
		t.Errorf("Desc(...).String() = %q, should contain 'DESC'", result)
	}
}

func TestOrderByWithFunction(t *testing.T) {
	expr := Count(Property("n", "id"))
	descExpr := Desc(expr)
	result := descExpr.String()
	// Should contain DESC and count
	if !containsString(result, "DESC") || !containsString(result, "count") {
		t.Errorf("Desc(Count(...)).String() = %q, should contain 'DESC' and 'count'", result)
	}
}

func TestOrderByExpressionLogicalOps(t *testing.T) {
	expr := Property("n", "age")
	orderExpr := Desc(expr)
	other := Property("n", "score")

	// Test And
	andExpr := orderExpr.And(other)
	if andExpr == nil {
		t.Error("orderExpr.And(other) returned nil")
	}

	// Test Or
	orExpr := orderExpr.Or(other)
	if orExpr == nil {
		t.Error("orderExpr.Or(other) returned nil")
	}

	// Test Not
	notExpr := orderExpr.Not()
	if notExpr == nil {
		t.Error("orderExpr.Not() returned nil")
	}
}


