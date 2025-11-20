package expr

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"simple string", "hello", "'hello'"},
		{"string with quotes", "it's", "'it\\'s'"},
		{"empty string", "", "''"},
		{"string with spaces", "hello world", "'hello world'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lit := String(tt.value)
			result := lit.String()
			if result != tt.expected {
				t.Errorf("String(%q).String() = %q, want %q", tt.value, result, tt.expected)
			}
		})
	}
}

func TestIntegerLiteral(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected string
	}{
		{"positive integer", 42, "42"},
		{"zero", 0, "0"},
		{"negative integer", -10, "-10"},
		{"large integer", 9223372036854775807, "9223372036854775807"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lit := Integer(tt.value)
			result := lit.String()
			if result != tt.expected {
				t.Errorf("Integer(%d).String() = %q, want %q", tt.value, result, tt.expected)
			}
		})
	}
}

func TestFloatLiteral(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{"positive float", 3.14, "3.14"},
		{"zero float", 0.0, "0"},
		{"negative float", -2.5, "-2.5"},
		{"small float", 0.001, "0.001"},
		{"large float", 123456.789, "123456.789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lit := Float(tt.value)
			result := lit.String()
			if result != tt.expected {
				t.Errorf("Float(%f).String() = %q, want %q", tt.value, result, tt.expected)
			}
		})
	}
}

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected string
	}{
		{"true", true, "true"},
		{"false", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lit := Boolean(tt.value)
			result := lit.String()
			if result != tt.expected {
				t.Errorf("Boolean(%v).String() = %q, want %q", tt.value, result, tt.expected)
			}
		})
	}
}

func TestNullLiteral(t *testing.T) {
	lit := Null()
	result := lit.String()
	expected := "NULL"
	if result != expected {
		t.Errorf("Null().String() = %q, want %q", result, expected)
	}
}

func TestListExpression(t *testing.T) {
	tests := []struct {
		name     string
		elements []interface{}
		expected string
	}{
		{"empty list", []interface{}{}, "[]"},
		{"single element", []interface{}{"hello"}, "['hello']"},
		{"multiple elements", []interface{}{1, 2, 3}, "[1, 2, 3]"},
		{"mixed types", []interface{}{"hello", 42, true}, "['hello', 42, true]"},
		{"nested lists", []interface{}{[]interface{}{1, 2}}, "[[1, 2]]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := make([]core.Expression, len(tt.elements))
			for i, v := range tt.elements {
				elements[i] = LiteralFromValue(v)
			}
			lit := List(elements...)
			result := lit.String()
			if result != tt.expected {
				t.Errorf("List(...).String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestMapLiteralExpression(t *testing.T) {
	tests := []struct {
		name     string
		entries  map[string]interface{}
		expected string // Note: map order is non-deterministic, so we check for key presence
		keys     []string
	}{
		{
			"simple map",
			map[string]interface{}{"name": "John", "age": 30},
			"",
			[]string{"name", "age"},
		},
		{
			"empty map",
			map[string]interface{}{},
			"{}",
			[]string{},
		},
		{
			"single entry",
			map[string]interface{}{"key": "value"},
			"{key: 'value'}",
			[]string{"key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries := make(map[string]core.Expression)
			for k, v := range tt.entries {
				entries[k] = LiteralFromValue(v)
			}
			lit := Map(entries)
			result := lit.String()

			if tt.expected != "" {
				if result != tt.expected {
					t.Errorf("Map(...).String() = %q, want %q", result, tt.expected)
				}
			} else {
				// Check that all keys are present
				for _, key := range tt.keys {
					if !containsString(result, key) {
						t.Errorf("Map(...).String() = %q, should contain key %q", result, key)
					}
				}
				// Check basic structure
				if result[0] != '{' || result[len(result)-1] != '}' {
					t.Errorf("Map(...).String() = %q, should start with { and end with }", result)
				}
			}
		})
	}
}

func TestLiteralFromValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"string", "hello", "'hello'"},
		{"int", 42, "42"},
		{"int64", int64(100), "100"},
		{"float64", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"nil", nil, "NULL"},
		{"slice", []interface{}{1, 2, 3}, "[1, 2, 3]"},
		{"map", map[string]interface{}{"key": "value"}, ""}, // Order non-deterministic
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lit := LiteralFromValue(tt.value)
			result := lit.String()
			if tt.expected != "" && result != tt.expected {
				t.Errorf("LiteralFromValue(%v).String() = %q, want %q", tt.value, result, tt.expected)
			}
			if tt.expected == "" && result == "" {
				t.Errorf("LiteralFromValue(%v).String() should not be empty", tt.value)
			}
		})
	}
}

func TestParameterExpression(t *testing.T) {
	tests := []struct {
		name     string
		param    string
		value    interface{}
		expected string
	}{
		{"named parameter", "name", "John", "$name"},
		{"empty name", "", 42, "$"},
		{"numeric name", "123", true, "$123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := Param(tt.param, tt.value)
			result := param.String()
			if result != tt.expected {
				t.Errorf("Param(%q, %v).String() = %q, want %q", tt.param, tt.value, result, tt.expected)
			}
		})
	}
}


