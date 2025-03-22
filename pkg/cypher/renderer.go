package cypher

import (
	"fmt"
	"strconv"
	"strings"
)

// cypherStatement represents a complete Cypher statement
type cypherStatement struct {
	builder     Buildable
	queryString string
	params      map[string]any
}

// Cypher renders this statement as a Cypher query string
func (s *cypherStatement) Cypher() string {
	if s.queryString == "" {
		// Create a new renderer and use it to render this statement
		renderer := NewRenderer()
		s.Accept(renderer)
		s.queryString = renderer.String()
	}
	return s.queryString
}

// Params returns the parameters for this statement
func (s *cypherStatement) Params() map[string]any {
	if s.params == nil {
		s.params = make(map[string]any)

		// Extract parameters based on the type of builder
		switch b := s.builder.(type) {
		case *MatchBuilder:
			// Check WHERE clause for parameters
			if b.where != nil && b.where.condition != nil {
				s.extractParamsFromExpression(b.where.condition)
			}
		case *WithBuilder:
			// Check WHERE clause for parameters
			if b.clause.where != nil && b.clause.where.condition != nil {
				s.extractParamsFromExpression(b.clause.where.condition)
			}
		case *ReturnBuilder:
			// Traverse the builder chain to find WHERE clauses
			if b.prev != nil {
				s.extractParamsFromBuilder(b.prev)
			}
		}
	}
	return s.params
}

// extractParamsFromBuilder recursively extracts parameters from a builder chain
func (s *cypherStatement) extractParamsFromBuilder(builder ClauseBuilder) {
	switch b := builder.(type) {
	case *MatchBuilder:
		// Extract parameters from WHERE clause
		if b.where != nil && b.where.condition != nil {
			s.extractParamsFromExpression(b.where.condition)
		}
		// Continue traversing
		if b.prev != nil {
			s.extractParamsFromBuilder(b.prev)
		}
	case *WithBuilder:
		// Extract parameters from WHERE clause
		if b.clause.where != nil && b.clause.where.condition != nil {
			s.extractParamsFromExpression(b.clause.where.condition)
		}
		// Continue traversing
		if b.prev != nil {
			s.extractParamsFromBuilder(b.prev)
		}
	}
}

// extractParamsFromExpression extracts parameters from an expression
func (s *cypherStatement) extractParamsFromExpression(expr Expression) {
	// Handle parameter expressions directly
	if param, ok := expr.(*parameter); ok && param.value != nil {
		s.params[param.name] = param.value
		return
	}

	// Handle different expression types recursively
	switch e := expr.(type) {
	case *comparisonExpression:
		s.extractParamsFromExpression(e.left)
		s.extractParamsFromExpression(e.right)
	case *binaryExpression:
		s.extractParamsFromExpression(e.left)
		s.extractParamsFromExpression(e.right)
	case *logicalExpression:
		s.extractParamsFromExpression(e.left)
		s.extractParamsFromExpression(e.right)
	case *stringOperation:
		s.extractParamsFromExpression(e.left)
		s.extractParamsFromExpression(e.right)
	case *notExpression:
		s.extractParamsFromExpression(e.expression)
	case *existsExpression:
		s.extractParamsFromExpression(e.expression)
	case *nullCheck:
		s.extractParamsFromExpression(e.expression)
	case *listExpression:
		for _, elem := range e.elements {
			s.extractParamsFromExpression(elem)
		}
	case *mapExpression:
		for _, value := range e.entries {
			s.extractParamsFromExpression(value)
		}
	}
}

// AddParam adds a parameter to this statement
func (s *cypherStatement) AddParam(name string, value any) {
	if s.params == nil {
		s.params = make(map[string]any)
	}
	s.params[name] = value
}

// Accept implements the Statement interface
func (s *cypherStatement) Accept(visitor StatementVisitor) any {
	return visitor.Visit(s.builder)
}

// Renderer is used to render Cypher statements to strings
type Renderer struct {
	sb strings.Builder
}

// NewRenderer creates a new Renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// String returns the rendered query string
func (r *Renderer) String() string {
	return r.sb.String()
}

// Visit implements the StatementVisitor interface
func (r *Renderer) Visit(element any) any {
	visitor := &cypherVisitor{renderer: r}
	return visitor.Visit(element)
}

// parameterCollector is a visitor that collects parameters from expressions
type parameterCollector struct {
	params map[string]any
}

// Visit implements the ExpressionVisitor interface
func (p *parameterCollector) Visit(element any) any {
	// Collect parameter expressions
	if pe, ok := element.(*parameter); ok && pe.value != nil {
		p.params[pe.name] = pe.value
		return nil
	}

	// Recursively visit child expressions
	switch expr := element.(type) {
	case *propertyExpression:
		if expr.subject != nil {
			expr.subject.Accept(p)
		}
	case *comparisonExpression:
		expr.left.Accept(p)
		expr.right.Accept(p)
	case *logicalExpression:
		expr.left.Accept(p)
		expr.right.Accept(p)
	case *notExpression:
		expr.expression.Accept(p)
	case *stringOperation:
		expr.left.Accept(p)
		expr.right.Accept(p)
	case *aliasedExpression:
		expr.expression.Accept(p)
	case *nodePattern:
		for _, propExpr := range expr.properties {
			propExpr.Accept(p)
		}
	case *relationshipPattern:
		for _, propExpr := range expr.properties {
			propExpr.Accept(p)
		}
	case *functionExpression:
		for _, arg := range expr.arguments {
			arg.Accept(p)
		}
	case *WhereClause:
		if expr.condition != nil {
			expr.condition.Accept(p)
		}
	case *MatchClause:
		for _, pattern := range expr.patterns {
			pattern.Accept(p)
		}
	case *ReturnClause:
		for _, returnExpr := range expr.expressions {
			returnExpr.Accept(p)
		}
	case *WithClause:
		for _, withExpr := range expr.expressions {
			withExpr.Accept(p)
		}
		if expr.where != nil && expr.where.condition != nil {
			expr.where.condition.Accept(p)
		}
	}

	return nil
}

// cypherVisitor implements the ExpressionVisitor interface to render Cypher
type cypherVisitor struct {
	renderer *Renderer
}

// Visit visits an expression
func (v *cypherVisitor) Visit(element any) any {
	switch e := element.(type) {
	case *MatchClause:
		return v.visitMatchClause(e)
	case *WhereClause:
		return v.visitWhereClause(e)
	case *ReturnClause:
		return v.visitReturnClause(e)
	case *WithClause:
		return v.visitWithClause(e)
	case *CreateClause:
		return v.visitCreateClause(e)
	case *MergeClause:
		return v.visitMergeClause(e)
	case *DeleteClause:
		return v.visitDeleteClause(e)
	case *SetClause:
		return v.visitSetClause(e)
	case *RemoveClause:
		return v.visitRemoveClause(e)
	case *UnwindClause:
		return v.visitUnwindClause(e)
	case *nodePattern:
		return v.visitNodePattern(e)
	case *relationshipPattern:
		return v.visitRelationshipPattern(e)
	case *patternPath:
		return v.visitPatternPath(e)
	case *propertyExpression:
		return v.visitPropertyExpression(e)
	case *schemaPropertyExpression:
		return v.visitSchemaPropertyExpression(e)
	case *comparisonExpression:
		return v.visitComparisonExpression(e)
	case *logicalExpression:
		return v.visitLogicalExpression(e)
	case *notExpression:
		return v.visitNotExpression(e)
	case *stringOperation:
		return v.visitStringOperation(e)
	case *literal:
		return v.visitLiteralExpression(e)
	case *parameter:
		return v.visitParameterExpression(e)
	case *listExpression:
		return v.visitListExpression(e)
	case *functionExpression:
		return v.visitFunctionExpression(e)
	case *aliasedExpression:
		return v.visitAliasedExpression(e)
	case ClauseBuilder:
		v.handleBuilderChain(e)
		return nil
	default:
		return nil
	}
}

// visitMatchClause renders a MATCH clause
func (v *cypherVisitor) visitMatchClause(clause *MatchClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")

	for i, pattern := range clause.Patterns() {
		if i > 0 {
			v.renderer.sb.WriteString(", ")
		}
		pattern.Accept(v)
	}
	return ""
}

// visitWhereClause renders a WHERE clause
func (v *cypherVisitor) visitWhereClause(clause *WhereClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")
	clause.Condition().Accept(v)
	return ""
}

// visitReturnClause renders a RETURN clause
func (v *cypherVisitor) visitReturnClause(clause *ReturnClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")

	for i, expr := range clause.Expressions() {
		if i > 0 {
			v.renderer.sb.WriteString(", ")
		}
		// For named nodes/relationships in a RETURN clause, just use the alias
		if node, ok := expr.(*nodePattern); ok && node.alias != "" {
			v.renderer.sb.WriteString(node.alias)
		} else if rel, ok := expr.(*relationshipPattern); ok && rel.alias != "" {
			v.renderer.sb.WriteString(rel.alias)
		} else {
			expr.Accept(v)
		}
	}

	// Render ORDER BY if present
	if len(clause.OrderBy()) > 0 {
		v.renderer.sb.WriteString(" ORDER BY ")
		for i, item := range clause.OrderBy() {
			if i > 0 {
				v.renderer.sb.WriteString(", ")
			}
			item.expression.Accept(v)
			if item.descending {
				v.renderer.sb.WriteString(" DESC")
			}
		}
	}

	// Render SKIP if present
	if clause.Skip() != nil {
		v.renderer.sb.WriteString(" SKIP ")
		clause.Skip().Accept(v)
	}

	// Render LIMIT if present
	if clause.Limit() != nil {
		v.renderer.sb.WriteString(" LIMIT ")
		clause.Limit().Accept(v)
	}
	return ""
}

// visitWithClause renders a WITH clause
func (v *cypherVisitor) visitWithClause(clause *WithClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")

	for i, expr := range clause.Expressions() {
		if i > 0 {
			v.renderer.sb.WriteString(", ")
		}
		// For named nodes/relationships in a WITH clause, just use the alias
		if node, ok := expr.(*nodePattern); ok && node.alias != "" {
			v.renderer.sb.WriteString(node.alias)
		} else if rel, ok := expr.(*relationshipPattern); ok && rel.alias != "" {
			v.renderer.sb.WriteString(rel.alias)
		} else {
			expr.Accept(v)
		}
	}

	// Render ORDER BY if present
	if len(clause.OrderBy()) > 0 {
		v.renderer.sb.WriteString(" ORDER BY ")
		for i, item := range clause.OrderBy() {
			if i > 0 {
				v.renderer.sb.WriteString(", ")
			}
			item.expression.Accept(v)
			if item.descending {
				v.renderer.sb.WriteString(" DESC")
			}
		}
	}

	// Render SKIP if present
	if clause.Skip() != nil {
		v.renderer.sb.WriteString(" SKIP ")
		clause.Skip().Accept(v)
	}

	// Render LIMIT if present
	if clause.Limit() != nil {
		v.renderer.sb.WriteString(" LIMIT ")
		clause.Limit().Accept(v)
	}

	// Render WHERE if present
	if clause.where != nil {
		v.renderer.sb.WriteString(" ")
		v.visitWhereClause(clause.where)
	}
	return ""
}

// visitCreateClause renders a CREATE clause
func (v *cypherVisitor) visitCreateClause(clause *CreateClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")

	for i, pattern := range clause.Patterns() {
		if i > 0 {
			v.renderer.sb.WriteString(", ")
		}
		pattern.Accept(v)
	}
	return ""
}

// visitMergeClause renders a MERGE clause
func (v *cypherVisitor) visitMergeClause(clause *MergeClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")
	clause.Pattern().Accept(v)
	return ""
}

// visitDeleteClause renders a DELETE clause
func (v *cypherVisitor) visitDeleteClause(clause *DeleteClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")

	for i, expr := range clause.Expressions() {
		if i > 0 {
			v.renderer.sb.WriteString(", ")
		}
		// For named nodes/relationships in a DELETE clause, just use the alias
		if node, ok := expr.(*nodePattern); ok && node.alias != "" {
			v.renderer.sb.WriteString(node.alias)
		} else if rel, ok := expr.(*relationshipPattern); ok && rel.alias != "" {
			v.renderer.sb.WriteString(rel.alias)
		} else {
			expr.Accept(v)
		}
	}
	return ""
}

// visitSetClause renders a SET clause
func (v *cypherVisitor) visitSetClause(clause *SetClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")

	for i, item := range clause.Items() {
		if i > 0 {
			v.renderer.sb.WriteString(", ")
		}

		if len(item.labels) > 0 {
			// Render label assignment
			v.renderer.sb.WriteString(item.nodeAlias)
			for _, label := range item.labels {
				v.renderer.sb.WriteString(":")
				v.renderer.sb.WriteString(label)
			}
		} else {
			// Render property assignment
			item.property.Accept(v)
			v.renderer.sb.WriteString(" ")
			v.renderer.sb.WriteString(item.operator)
			v.renderer.sb.WriteString(" ")
			item.value.Accept(v)
		}
	}
	return ""
}

// visitRemoveClause renders a REMOVE clause
func (v *cypherVisitor) visitRemoveClause(clause *RemoveClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")

	for i, item := range clause.Items() {
		if i > 0 {
			v.renderer.sb.WriteString(", ")
		}

		if item.isLabel {
			// Render label removal
			v.renderer.sb.WriteString(item.nodeAlias)
			v.renderer.sb.WriteString(":")
			v.renderer.sb.WriteString(item.label)
		} else {
			// Render property removal
			item.property.Accept(v)
		}
	}
	return ""
}

// visitUnwindClause renders an UNWIND clause
func (v *cypherVisitor) visitUnwindClause(clause *UnwindClause) string {
	v.renderer.sb.WriteString(clause.AsString())
	v.renderer.sb.WriteString(" ")
	clause.Expression().Accept(v)
	v.renderer.sb.WriteString(" AS ")
	v.renderer.sb.WriteString(clause.Alias())
	return ""
}

// visitNodePattern renders a node pattern
func (v *cypherVisitor) visitNodePattern(node *nodePattern) string {
	v.renderer.sb.WriteString("(")

	if node.alias != "" {
		v.renderer.sb.WriteString(node.alias)
	}

	for _, label := range node.labels {
		v.renderer.sb.WriteString(":`")
		v.renderer.sb.WriteString(label)
		v.renderer.sb.WriteString("`")
	}

	if len(node.properties) > 0 {
		v.renderer.sb.WriteString(" {")

		i := 0
		for key, value := range node.properties {
			if i > 0 {
				v.renderer.sb.WriteString(", ")
			}
			v.renderer.sb.WriteString(key)
			v.renderer.sb.WriteString(": ")
			value.Accept(v)
			i++
		}

		v.renderer.sb.WriteString("}")
	}

	v.renderer.sb.WriteString(")")
	return ""
}

// visitRelationshipPattern renders a relationship pattern
func (v *cypherVisitor) visitRelationshipPattern(rel *relationshipPattern) string {
	switch rel.direction {
	case OUTGOING:
		v.renderer.sb.WriteString("-[")
	case INCOMING:
		v.renderer.sb.WriteString("<-[")
	case BIDIRECTIONAL:
		v.renderer.sb.WriteString("-[")
	}

	if rel.alias != "" {
		v.renderer.sb.WriteString(rel.alias)
	}

	for i, typ := range rel.types {
		if i > 0 {
			v.renderer.sb.WriteString("|")
		}
		v.renderer.sb.WriteString(":`")
		v.renderer.sb.WriteString(typ)
		v.renderer.sb.WriteString("`")
	}

	if rel.minHops != nil || rel.maxHops != nil {
		v.renderer.sb.WriteString("*")
		if rel.minHops != nil {
			v.renderer.sb.WriteString(fmt.Sprintf("%d", *rel.minHops))
		}
		v.renderer.sb.WriteString("..")
		if rel.maxHops != nil {
			v.renderer.sb.WriteString(fmt.Sprintf("%d", *rel.maxHops))
		}
	}

	if len(rel.properties) > 0 {
		v.renderer.sb.WriteString(" {")

		i := 0
		for key, value := range rel.properties {
			if i > 0 {
				v.renderer.sb.WriteString(", ")
			}
			v.renderer.sb.WriteString(key)
			v.renderer.sb.WriteString(": ")
			value.Accept(v)
			i++
		}

		v.renderer.sb.WriteString("}")
	}

	v.renderer.sb.WriteString("]")

	switch rel.direction {
	case OUTGOING:
		v.renderer.sb.WriteString("->")
	case INCOMING:
		v.renderer.sb.WriteString("-")
	case BIDIRECTIONAL:
		v.renderer.sb.WriteString("-")
	}
	return ""
}

// visitPatternPath renders a pattern path
func (v *cypherVisitor) visitPatternPath(path *patternPath) string {
	if path.alias != "" {
		v.renderer.sb.WriteString(path.alias)
		v.renderer.sb.WriteString(" = ")
	}

	for _, element := range path.elements {
		element.Accept(v)
	}
	return ""
}

// visitPropertyExpression renders a property expression
func (v *cypherVisitor) visitPropertyExpression(property *propertyExpression) string {
	if property.subject != nil {
		// If the subject is a node or relationship pattern with an alias, use that alias
		if node, ok := property.subject.(*nodePattern); ok && node.alias != "" {
			v.renderer.sb.WriteString(node.alias)
		} else if rel, ok := property.subject.(*relationshipPattern); ok && rel.alias != "" {
			v.renderer.sb.WriteString(rel.alias)
		} else {
			property.subject.Accept(v)
		}
	}
	v.renderer.sb.WriteString(".")
	v.renderer.sb.WriteString(property.propertyName)
	return ""
}

// visitSchemaPropertyExpression renders a schema property expression
func (v *cypherVisitor) visitSchemaPropertyExpression(property *schemaPropertyExpression) string {
	if property.nodeType != nil {
		// Use the node alias from the schema type
		if property.nodeAlias != "" {
			v.renderer.sb.WriteString(property.nodeAlias)
		} else {
			// Fallback to a generic alias
			v.renderer.sb.WriteString("p")
		}
	} else if property.relType != nil {
		// Use the relationship alias from the schema type
		if property.relAlias != "" {
			v.renderer.sb.WriteString(property.relAlias)
		} else {
			// Fallback to a generic alias
			v.renderer.sb.WriteString("r")
		}
	}

	v.renderer.sb.WriteString(".")
	v.renderer.sb.WriteString(property.propertyName)
	return ""
}

// visitComparisonExpression renders a comparison expression
func (v *cypherVisitor) visitComparisonExpression(comparison *comparisonExpression) string {
	comparison.left.Accept(v)
	v.renderer.sb.WriteString(" ")
	v.renderer.sb.WriteString(comparison.operator)
	v.renderer.sb.WriteString(" ")
	comparison.right.Accept(v)
	return ""
}

// visitLogicalExpression renders a logical expression
func (v *cypherVisitor) visitLogicalExpression(logical *logicalExpression) string {
	v.renderer.sb.WriteString("(")
	logical.left.Accept(v)
	v.renderer.sb.WriteString(" ")
	v.renderer.sb.WriteString(logical.operator)
	v.renderer.sb.WriteString(" ")
	logical.right.Accept(v)
	v.renderer.sb.WriteString(")")
	return ""
}

// visitNotExpression renders a NOT expression
func (v *cypherVisitor) visitNotExpression(not *notExpression) string {
	v.renderer.sb.WriteString("NOT ")
	v.renderer.sb.WriteString("(")

	// Make sure we have a non-nil expression to avoid empty parentheses
	if not.expression != nil {
		not.expression.Accept(v)
	} else {
		// Provide a safe default to avoid empty parentheses
		v.renderer.sb.WriteString("true")
	}

	v.renderer.sb.WriteString(")")
	return ""
}

// visitStringOperation renders a string operation expression
func (v *cypherVisitor) visitStringOperation(op *stringOperation) string {
	op.left.Accept(v)
	v.renderer.sb.WriteString(" ")
	v.renderer.sb.WriteString(op.operator)
	v.renderer.sb.WriteString(" ")
	op.right.Accept(v)
	return ""
}

// visitLiteralExpression renders a literal expression
func (v *cypherVisitor) visitLiteralExpression(literal *literal) string {
	switch value := literal.value.(type) {
	case string:
		v.renderer.sb.WriteString("'")
		v.renderer.sb.WriteString(value)
		v.renderer.sb.WriteString("'")
	case bool:
		v.renderer.sb.WriteString(strconv.FormatBool(value))
	case int:
		v.renderer.sb.WriteString(strconv.Itoa(value))
	case int64:
		v.renderer.sb.WriteString(strconv.FormatInt(value, 10))
	case float64:
		v.renderer.sb.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
	case nil:
		v.renderer.sb.WriteString("NULL")
	default:
		v.renderer.sb.WriteString(fmt.Sprintf("%v", value))
	}
	return ""
}

// visitParameterExpression renders a parameter expression
func (v *cypherVisitor) visitParameterExpression(param *parameter) string {
	v.renderer.sb.WriteString("$")
	v.renderer.sb.WriteString(param.name)
	return ""
}

// visitListExpression renders a list expression
func (v *cypherVisitor) visitListExpression(list *listExpression) string {
	v.renderer.sb.WriteString("[")

	for i, element := range list.elements {
		if i > 0 {
			v.renderer.sb.WriteString(", ")
		}
		element.Accept(v)
	}

	v.renderer.sb.WriteString("]")
	return ""
}

// visitFunctionExpression renders a function expression
func (v *cypherVisitor) visitFunctionExpression(function *functionExpression) string {
	v.renderer.sb.WriteString(function.name)
	v.renderer.sb.WriteString("(")

	for i, arg := range function.arguments {
		if i > 0 {
			v.renderer.sb.WriteString(", ")
		}
		// For named nodes/relationships in function arguments, just use the alias
		if node, ok := arg.(*nodePattern); ok && node.alias != "" {
			v.renderer.sb.WriteString(node.alias)
		} else if rel, ok := arg.(*relationshipPattern); ok && rel.alias != "" {
			v.renderer.sb.WriteString(rel.alias)
		} else {
			arg.Accept(v)
		}
	}

	v.renderer.sb.WriteString(")")
	return ""
}

// visitAliasedExpression renders an aliased expression
func (v *cypherVisitor) visitAliasedExpression(aliased *aliasedExpression) string {
	aliased.expression.Accept(v)
	v.renderer.sb.WriteString(" AS ")
	v.renderer.sb.WriteString(aliased.aliasName)
	return ""
}

// Handle the builder
func (v *cypherVisitor) handleBuilderChain(builder ClauseBuilder) {
	switch e := builder.(type) {
	case *MatchBuilder:
		// Handle previous clauses first
		if e.prev != nil {
			v.handleBuilderChain(e.prev)
			v.renderer.sb.WriteString(" ")
		}

		// Then render the MATCH clause
		v.visitMatchClause(e.clause)

		// And add the WHERE clause if present
		if e.where != nil {
			v.renderer.sb.WriteString(" ")
			v.visitWhereClause(e.where)
		}

	case *WithBuilder:
		// Handle previous clauses
		if e.prev != nil {
			v.handleBuilderChain(e.prev)
			v.renderer.sb.WriteString(" ")
		}

		// Render the WITH clause (which already handles its own WHERE clause)
		v.visitWithClause(e.clause)

	case *ReturnBuilder:
		// Handle previous clauses
		if e.prev != nil {
			v.handleBuilderChain(e.prev)
			v.renderer.sb.WriteString(" ")
		}

		// Render the RETURN clause
		v.visitReturnClause(e.clause)

	case *CreateBuilder:
		// CreateBuilder doesn't have a previous chain
		v.visitCreateClause(e.clause)

	case *MergeBuilder:
		// MergeBuilder doesn't have a previous chain
		v.visitMergeClause(e.clause)

	case *DeleteBuilder:
		// DeleteBuilder doesn't have a previous chain
		v.visitDeleteClause(e.clause)

	case *SetBuilder:
		// SetBuilder doesn't have a previous chain
		v.visitSetClause(e.clause)

	case *RemoveBuilder:
		// RemoveBuilder doesn't have a previous chain
		v.visitRemoveClause(e.clause)

	case *UnwindBuilder:
		// UnwindBuilder doesn't have a previous chain
		if e.clause != nil {
			v.visitUnwindClause(e.clause)
		}
	}
}

// renderCypher renders a statement to a Cypher string
func renderCypher(statement Statement) string {
	renderer := NewRenderer()
	statement.Accept(renderer)
	return renderer.String()
}
