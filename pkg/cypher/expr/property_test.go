package expr

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

func TestPropertyExpression(t *testing.T) {
	tests := []struct {
		name     string
		entity   string
		property string
		contains string
	}{
		{"simple property", "n", "name", "name"},
		{"property with underscore", "user", "first_name", "first_name"},
		{"property with number", "p", "age2", "age2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prop := Property(tt.entity, tt.property)
			result := prop.String()
			if !containsString(result, tt.contains) {
				t.Errorf("Property(%q, %q).String() = %q, should contain %q", tt.entity, tt.property, result, tt.contains)
			}
		})
	}
}

func TestPropertyExpressionWithChain(t *testing.T) {
	prop := Property("user", "address", "city", "name")
	result := prop.String()
	// Should contain all property names in the chain
	if !containsString(result, "address") || !containsString(result, "city") || !containsString(result, "name") {
		t.Errorf("Property with chain = %q, should contain all property names", result)
	}
}

func TestNewProperty(t *testing.T) {
	subject := String("n")
	prop := NewProperty(subject, "name")
	result := prop.String()
	// Should contain the property name
	if !containsString(result, "name") {
		t.Errorf("NewProperty(...).String() = %q, should contain 'name'", result)
	}
}

func TestPropertyComparisons(t *testing.T) {
	propExpr := Property("n", "age")
	prop, ok := propExpr.(*PropertyExpression)
	if !ok {
		t.Fatal("Property() should return *PropertyExpression")
	}

	tests := []struct {
		name     string
		fn       func() core.Expression
		contains string
	}{
		{"Eq", func() core.Expression { return prop.Eq(30) }, "="},
		{"Gt", func() core.Expression { return prop.Gt(30) }, ">"},
		{"Lt", func() core.Expression { return prop.Lt(30) }, "<"},
		{"Gte", func() core.Expression { return prop.Gte(30) }, ">="},
		{"Lte", func() core.Expression { return prop.Lte(30) }, "<="},
		{"IsNull", func() core.Expression { return prop.IsNull() }, "IS"},
		{"IsNotNull", func() core.Expression { return prop.IsNotNull() }, "IS NOT"},
		{"In", func() core.Expression { return prop.In(1, 2, 3) }, "IN"},
		{"Contains", func() core.Expression { return prop.Contains("test") }, "CONTAINS"},
		{"StartsWith", func() core.Expression { return prop.StartsWith("test") }, "STARTS WITH"},
		{"EndsWith", func() core.Expression { return prop.EndsWith("test") }, "ENDS WITH"},
		{"RegularExpression", func() core.Expression { return prop.RegularExpression(".*") }, "=~"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := tt.fn()
			result := expr.String()
			if !containsString(result, tt.contains) {
				t.Errorf("%s() result = %q, should contain %q", tt.name, result, tt.contains)
			}
		})
	}
}

func TestPropertyLogicalOps(t *testing.T) {
	prop := Property("n", "age")
	other := Property("n", "score")

	// Test And
	andExpr := prop.And(other)
	if andExpr == nil {
		t.Error("prop.And(other) returned nil")
	}

	// Test Or
	orExpr := prop.Or(other)
	if orExpr == nil {
		t.Error("prop.Or(other) returned nil")
	}

	// Test Not
	notExpr := prop.Not()
	if notExpr == nil {
		t.Error("prop.Not() returned nil")
	}
}


