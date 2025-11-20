package expr

import (
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

func TestConcat(t *testing.T) {
	tests := []struct {
		name     string
		exprs    []interface{}
		contains string
	}{
		{"two strings", []interface{}{"hello", "world"}, "+"},
		{"three strings", []interface{}{"a", "b", "c"}, "+"},
		{"single string", []interface{}{"hello"}, ""}, // Should return as-is
		{"empty", []interface{}{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exprs := make([]core.Expression, len(tt.exprs))
			for i, v := range tt.exprs {
				exprs[i] = LiteralFromValue(v)
			}
			concatExpr := Concat(exprs...)
			result := concatExpr.String()

			if tt.contains != "" && !containsString(result, tt.contains) {
				t.Errorf("Concat(...).String() = %q, should contain %q", result, tt.contains)
			}
		})
	}
}

func TestSubstring(t *testing.T) {
	expr := Property("n", "name")
	start := Integer(1)
	length := Integer(5)
	substrExpr := Substring(expr, start, length)
	result := substrExpr.String()
	// Should contain substring and the property
	if !containsString(result, "substring") {
		t.Errorf("Substring(...).String() = %q, should contain 'substring'", result)
	}

	// Test without length
	substrExpr2 := Substring(expr, start)
	result2 := substrExpr2.String()
	// Should contain substring
	if !containsString(result2, "substring") {
		t.Errorf("Substring(..., start).String() = %q, should contain 'substring'", result2)
	}
}

func TestReplace(t *testing.T) {
	expr := Property("n", "text")
	search := String("old")
	replace := String("new")
	replaceExpr := Replace(expr, search, replace)
	result := replaceExpr.String()
	// Should contain replace
	if !containsString(result, "replace") {
		t.Errorf("Replace(...).String() = %q, should contain 'replace'", result)
	}
}

func TestSplit(t *testing.T) {
	expr := Property("n", "text")
	delimiter := String(",")
	splitExpr := Split(expr, delimiter)
	result := splitExpr.String()
	// Should contain split
	if !containsString(result, "split") {
		t.Errorf("Split(...).String() = %q, should contain 'split'", result)
	}
}

func TestToLower(t *testing.T) {
	expr := Property("n", "name")
	toLowerExpr := ToLower(expr)
	result := toLowerExpr.String()
	// Should contain toLower
	if !containsString(result, "toLower") {
		t.Errorf("ToLower(...).String() = %q, should contain 'toLower'", result)
	}
}

func TestToUpper(t *testing.T) {
	expr := Property("n", "name")
	toUpperExpr := ToUpper(expr)
	result := toUpperExpr.String()
	// Should contain toUpper and the property
	if !containsString(result, "toUpper") {
		t.Errorf("ToUpper(...).String() = %q, should contain 'toUpper'", result)
	}
}

func TestTrim(t *testing.T) {
	expr := Property("n", "text")
	trimExpr := Trim(expr)
	result := trimExpr.String()
	// Should contain trim
	if !containsString(result, "trim") {
		t.Errorf("Trim(...).String() = %q, should contain 'trim'", result)
	}
}

func TestLTrim(t *testing.T) {
	expr := Property("n", "text")
	lTrimExpr := LTrim(expr)
	result := lTrimExpr.String()
	// Should contain lTrim
	if !containsString(result, "lTrim") {
		t.Errorf("LTrim(...).String() = %q, should contain 'lTrim'", result)
	}
}

func TestRTrim(t *testing.T) {
	expr := Property("n", "text")
	rTrimExpr := RTrim(expr)
	result := rTrimExpr.String()
	// Should contain rTrim
	if !containsString(result, "rTrim") {
		t.Errorf("RTrim(...).String() = %q, should contain 'rTrim'", result)
	}
}

func TestStringFunctionChaining(t *testing.T) {
	expr := Property("n", "name")
	// Chain: toUpper(trim(n.name))
	trimmed := Trim(expr)
	upper := ToUpper(trimmed)
	result := upper.String()
	// Should contain both toUpper and trim
	if !containsString(result, "toUpper") || !containsString(result, "trim") {
		t.Errorf("Chained string functions = %q, should contain 'toUpper' and 'trim'", result)
	}
}


