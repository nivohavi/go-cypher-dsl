package expr

import (
	"testing"
)

func TestAnd(t *testing.T) {
	left := Property("n", "age")
	right := Property("n", "active")
	andExpr := And(left, right)
	result := andExpr.String()
	// Should contain AND
	if !containsString(result, "AND") {
		t.Errorf("And(...).String() = %q, should contain 'AND'", result)
	}
}

func TestOr(t *testing.T) {
	left := Property("n", "status")
	right := String("active")
	orExpr := Or(left, right)
	result := orExpr.String()
	// Should contain OR
	if !containsString(result, "OR") {
		t.Errorf("Or(...).String() = %q, should contain 'OR'", result)
	}
}

func TestXor(t *testing.T) {
	left := Boolean(true)
	right := Boolean(false)
	xorExpr := Xor(left, right)
	result := xorExpr.String()
	// Should contain XOR
	if !containsString(result, "XOR") {
		t.Errorf("Xor(...).String() = %q, should contain 'XOR'", result)
	}
}

func TestNot(t *testing.T) {
	expr := Property("n", "deleted")
	notExpr := Not(expr)
	result := notExpr.String()
	// Should contain NOT
	if !containsString(result, "NOT") {
		t.Errorf("Not(...).String() = %q, should contain 'NOT'", result)
	}
}

func TestNotDoubleNegation(t *testing.T) {
	expr := Property("n", "active")
	notExpr := Not(expr)
	notNotExpr := notExpr.Not()
	// Double negation should cancel out
	if notNotExpr.String() != expr.String() {
		t.Errorf("Not(Not(...)) should equal original expression")
	}
}

func TestComplexLogicalExpressions(t *testing.T) {
	age := Property("n", "age")
	active := Property("n", "active")
	name := Property("n", "name")

	// (age > 30) AND (active = true)
	expr1 := GreaterThan(age, Integer(30))
	expr2 := Equals(active, Boolean(true))
	andExpr := And(expr1, expr2)
	result1 := andExpr.String()
	if !containsString(result1, "AND") {
		t.Error("Complex AND expression should contain 'AND'")
	}

	// (name = 'John') OR (name = 'Jane')
	expr3 := Equals(name, String("John"))
	expr4 := Equals(name, String("Jane"))
	orExpr := Or(expr3, expr4)
	result2 := orExpr.String()
	if !containsString(result2, "OR") {
		t.Error("Complex OR expression should contain 'OR'")
	}

	// NOT (age < 18)
	expr5 := LessThan(age, Integer(18))
	notExpr := Not(expr5)
	result3 := notExpr.String()
	if !containsString(result3, "NOT") {
		t.Error("NOT expression should contain 'NOT'")
	}
}

func TestNestedLogicalExpressions(t *testing.T) {
	age := Property("n", "age")
	active := Property("n", "active")
	name := Property("n", "name")

	// ((age > 30) AND (active = true)) OR (name = 'Admin')
	innerAnd := And(
		GreaterThan(age, Integer(30)),
		Equals(active, Boolean(true)),
	)
	outerOr := Or(innerAnd, Equals(name, String("Admin")))
	result := outerOr.String()

	if !containsString(result, "OR") {
		t.Error("Nested expression should contain 'OR'")
	}
	if !containsString(result, "AND") {
		t.Error("Nested expression should contain 'AND'")
	}
}

func TestLogicalExpressionChaining(t *testing.T) {
	age := Property("n", "age")
	active := Property("n", "active")
	name := Property("n", "name")

	// Chain: age > 30 AND active = true AND name = 'John'
	expr := GreaterThan(age, Integer(30)).
		And(Equals(active, Boolean(true))).
		And(Equals(name, String("John")))

	result := expr.String()
	if !containsString(result, "AND") {
		t.Error("Chained expression should contain 'AND'")
	}
}

func TestLogicalExpressionMethods(t *testing.T) {
	left := Property("n", "age")
	right := Property("n", "score")
	andExpr := And(left, right)

	// Test that logical expressions can be chained
	orExpr := andExpr.Or(left)
	if orExpr == nil {
		t.Error("andExpr.Or() returned nil")
	}

	// Xor only works on LogicalExpression
	if logExpr, ok := andExpr.(*LogicalExpression); ok {
		xorExpr := logExpr.Xor(right)
		if xorExpr == nil {
			t.Error("andExpr.Xor() returned nil")
		}
	}

	notExpr := andExpr.Not()
	if notExpr == nil {
		t.Error("andExpr.Not() returned nil")
	}
}


