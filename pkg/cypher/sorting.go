package cypher

import "fmt"

// SortItem represents a sort item in a Cypher ORDER BY clause
type SortItem struct {
	expression Expression
	descending bool
}

// Accept implements the Expression interface for SortItem
func (s SortItem) Accept(visitor ExpressionVisitor) any {
	return s.expression.Accept(visitor)
}

// String implements the Expression interface for SortItem
func (s SortItem) String() string {
	if s.descending {
		return fmt.Sprintf("%s DESC", s.expression)
	}
	return fmt.Sprintf("%s ASC", s.expression)
}

// Desc creates a descending sort item from an expression
func Desc(expr Expression) Expression {
	return SortItem{
		expression: expr,
		descending: true,
	}
}

// Asc creates an ascending sort item from an expression
func Asc(expr Expression) Expression {
	return SortItem{
		expression: expr,
		descending: false,
	}
}
