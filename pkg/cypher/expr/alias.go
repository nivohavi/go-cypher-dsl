package expr

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// AliasExpression represents an expression with an alias (e.g., n AS person)
type AliasExpression struct {
	Expression core.Expression
	Alias      string
}

// Accept implements the Expression interface
func (a *AliasExpression) Accept(visitor core.ExpressionVisitor) any {
	return visitor.Visit(a)
}

// quoteIdentifier quotes an identifier with backticks if it contains special characters
// or if it's already quoted, returns it as-is
func quoteIdentifier(identifier string) string {
	// If already quoted, return as-is
	if len(identifier) >= 2 && identifier[0] == '`' && identifier[len(identifier)-1] == '`' {
		return identifier
	}

	// Check if identifier needs quoting (contains spaces, dots, or other special characters)
	needsQuoting := false
	for _, r := range identifier {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			needsQuoting = true
			break
		}
	}

	// Also quote if it starts with a digit
	if len(identifier) > 0 && unicode.IsDigit(rune(identifier[0])) {
		needsQuoting = true
	}

	if needsQuoting {
		// Escape backticks inside the identifier
		escaped := strings.ReplaceAll(identifier, "`", "``")
		return fmt.Sprintf("`%s`", escaped)
	}

	return identifier
}

// String returns a string representation of this alias expression
func (a *AliasExpression) String() string {
	quotedAlias := quoteIdentifier(a.Alias)
	return fmt.Sprintf("%s AS %s", a.Expression.String(), quotedAlias)
}

// And creates a logical AND with another expression
func (a *AliasExpression) And(other core.Expression) core.Expression {
	return And(a, other)
}

// Or creates a logical OR with another expression
func (a *AliasExpression) Or(other core.Expression) core.Expression {
	return Or(a, other)
}

// Not creates a logical NOT of this expression
func (a *AliasExpression) Not() core.Expression {
	return Not(a)
}

// As creates an alias expression
func As(expression core.Expression, alias string) core.Expression {
	return &AliasExpression{
		Expression: expression,
		Alias:      alias,
	}
}
