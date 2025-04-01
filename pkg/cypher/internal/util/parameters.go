package util

import (
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// ExtractParameters extracts parameters from expressions recursively into the given parameters map
func ExtractParameters(expr core.Expression, paramsMap map[string]any) {
	if expr == nil {
		return
	}

	// Handle direct parameter expressions
	if paramExpr, ok := expr.(interface {
		Name() string
		Value() any
	}); ok {
		paramsMap[paramExpr.Name()] = paramExpr.Value()
		return
	}

	// Handle expression containers
	if container, ok := expr.(interface{ Expressions() []core.Expression }); ok {
		for _, subExpr := range container.Expressions() {
			ExtractParameters(subExpr, paramsMap)
		}
	}

	// Handle binary expressions (like comparisons)
	if binaryExpr, ok := expr.(interface {
		Left() core.Expression
		Right() core.Expression
	}); ok {
		if binaryExpr.Left() != nil {
			ExtractParameters(binaryExpr.Left(), paramsMap)
		}
		if binaryExpr.Right() != nil {
			ExtractParameters(binaryExpr.Right(), paramsMap)
		}
	}
}
