package expr

import (
	"testing"
)

func TestVariableExpression(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		expected string
	}{
		{"simple variable", "n", "n"},
		{"descriptive variable", "person", "person"},
		{"variable with underscore", "user_name", "user_name"},
		{"single letter", "a", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariableExpression(tt.varName)
			result := v.String()
			if result != tt.expected {
				t.Errorf("NewVariableExpression(%q).String() = %q, want %q", tt.varName, result, tt.expected)
			}
			if v.Name() != tt.varName {
				t.Errorf("NewVariableExpression(%q).Name() = %q, want %q", tt.varName, v.Name(), tt.varName)
			}
		})
	}
}

func TestVariableExpressionLogicalOps(t *testing.T) {
	v1 := NewVariableExpression("a")
	v2 := NewVariableExpression("b")

	// Test And
	andExpr := v1.And(v2)
	if andExpr == nil {
		t.Error("v1.And(v2) returned nil")
	}

	// Test Or
	orExpr := v1.Or(v2)
	if orExpr == nil {
		t.Error("v1.Or(v2) returned nil")
	}

	// Test Not
	notExpr := v1.Not()
	if notExpr == nil {
		t.Error("v1.Not() returned nil")
	}
}

