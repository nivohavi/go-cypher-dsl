package expr

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

func TestFunction(t *testing.T) {
	tests := []struct {
		name     string
		fnName   string
		args     []interface{}
		expected string
	}{
		{"no args", "myFunc", []interface{}{}, "myFunc()"},
		{"single arg", "abs", []interface{}{-10}, "abs(-10)"},
		{"multiple args", "coalesce", []interface{}{"a", "b", "c"}, "coalesce('a', 'b', 'c')"},
		{"mixed types", "func", []interface{}{"str", 42, true}, "func('str', 42, true)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]core.Expression, len(tt.args))
			for i, v := range tt.args {
				args[i] = LiteralFromValue(v)
			}
			fn := Function(tt.fnName, args...)
			result := fn.String()
			if result != tt.expected {
				t.Errorf("Function(%q, ...).String() = %q, want %q", tt.fnName, result, tt.expected)
			}
		})
	}
}

func TestCount(t *testing.T) {
	expr := Property("n", "id")
	countExpr := Count(expr)
	result := countExpr.String()
	// Should contain count
	if !containsString(result, "count") {
		t.Errorf("Count(...).String() = %q, should contain 'count'", result)
	}
}

func TestCountStar(t *testing.T) {
	countStarExpr := CountStar()
	result := countStarExpr.String()
	// Should contain count and *
	if !containsString(result, "count") {
		t.Errorf("CountStar().String() = %q, should contain 'count'", result)
	}
	if !containsString(result, "*") {
		t.Errorf("CountStar().String() = %q, should contain '*'", result)
	}
}

func TestSum(t *testing.T) {
	expr := Property("n", "price")
	sumExpr := Sum(expr)
	result := sumExpr.String()
	// Should contain sum
	if !containsString(result, "sum") {
		t.Errorf("Sum(...).String() = %q, should contain 'sum'", result)
	}
}

func TestAvg(t *testing.T) {
	expr := Property("n", "score")
	avgExpr := Avg(expr)
	result := avgExpr.String()
	// Should contain avg
	if !containsString(result, "avg") {
		t.Errorf("Avg(...).String() = %q, should contain 'avg'", result)
	}
}

func TestMin(t *testing.T) {
	expr := Property("n", "age")
	minExpr := Min(expr)
	result := minExpr.String()
	// Should contain min
	if !containsString(result, "min") {
		t.Errorf("Min(...).String() = %q, should contain 'min'", result)
	}
}

func TestMax(t *testing.T) {
	expr := Property("n", "age")
	maxExpr := Max(expr)
	result := maxExpr.String()
	// Should contain max
	if !containsString(result, "max") {
		t.Errorf("Max(...).String() = %q, should contain 'max'", result)
	}
}

func TestCollect(t *testing.T) {
	expr := Property("n", "name")
	collectExpr := Collect(expr)
	result := collectExpr.String()
	// Should contain collect
	if !containsString(result, "collect") {
		t.Errorf("Collect(...).String() = %q, should contain 'collect'", result)
	}
}

func TestDistinct(t *testing.T) {
	expr := Property("n", "name")
	distinctExpr := Distinct(expr)
	result := distinctExpr.String()
	// Should contain DISTINCT
	if !containsString(result, "DISTINCT") {
		t.Errorf("Distinct(...).String() = %q, should contain 'DISTINCT'", result)
	}
}

func TestFunctionDistinct(t *testing.T) {
	expr := Property("n", "name")
	countExpr := Count(expr)
	distinctCount := countExpr.(*FunctionExpression).Distinct()
	result := distinctCount.String()
	// Should contain count and DISTINCT
	if !containsString(result, "count") {
		t.Errorf("Count(...).Distinct().String() = %q, should contain 'count'", result)
	}
	if !containsString(result, "DISTINCT") {
		t.Errorf("Count(...).Distinct().String() = %q, should contain 'DISTINCT'", result)
	}
}

func TestFunctionAs(t *testing.T) {
	expr := Property("n", "count")
	countExpr := Count(expr)
	aliased := countExpr.(*FunctionExpression).As("total")
	result := aliased.String()
	// Should contain AS and the alias
	if !containsString(result, "AS") || !containsString(result, "total") {
		t.Errorf("Count(...).As(...).String() = %q, should contain 'AS' and 'total'", result)
	}
}

func TestFunctionLogicalOps(t *testing.T) {
	expr1 := Count(Property("n", "id"))
	expr2 := Sum(Property("n", "value"))

	// Test And
	andExpr := expr1.And(expr2)
	if andExpr == nil {
		t.Error("expr1.And(expr2) returned nil")
	}

	// Test Or
	orExpr := expr1.Or(expr2)
	if orExpr == nil {
		t.Error("expr1.Or(expr2) returned nil")
	}

	// Test Not
	notExpr := expr1.Not()
	if notExpr == nil {
		t.Error("expr1.Not() returned nil")
	}
}


