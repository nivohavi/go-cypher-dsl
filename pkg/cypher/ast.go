package cypher

import (
	"fmt"
	"strconv"
	"strings"
)

// nodePattern represents a node pattern in Cypher (e.g., (n:Person))
type nodePattern struct {
	labels     []string
	alias      string
	properties map[string]Expression
}

// Node creates a new node pattern with the given labels
func Node(labels ...string) NodeExpression {
	return &nodePattern{
		labels:     labels,
		properties: make(map[string]Expression),
	}
}

// Named sets the alias for this node pattern
func (n *nodePattern) Named(alias string) NodeExpression {
	clone := *n
	clone.alias = alias
	return &clone
}

// WithLabels adds labels to this node pattern
func (n *nodePattern) WithLabels(labels ...string) NodeExpression {
	clone := *n
	clone.labels = append(clone.labels, labels...)
	return &clone
}

// WithProperties adds properties to this node pattern
func (n *nodePattern) WithProperties(properties map[string]Expression) NodeExpression {
	clone := *n
	for k, v := range properties {
		clone.properties[k] = v
	}
	return &clone
}

// WithProps adds properties with automatic conversion to expressions
func (n *nodePattern) WithProps(properties map[string]interface{}) Expression {
	clone := *n
	for k, v := range properties {
		clone.properties[k] = autoExpression(v)
	}
	return &clone
}

// Props is an alias for WithProps
func (n *nodePattern) Props(properties map[string]interface{}) Expression {
	return n.WithProps(properties)
}

// Property returns a property access expression for this node
func (n *nodePattern) Property(propertyName string) PropertyExpression {
	return &propertyExpression{
		subject:      n,
		propertyName: propertyName,
	}
}

// Prop is an alias for Property
func (n *nodePattern) Prop(name string) PropertyExpression {
	return n.Property(name)
}

// RelationshipTo creates a relationship from this node to another
func (n *nodePattern) RelationshipTo(other NodeExpression, types ...string) RelationshipPattern {
	rel := &relationshipPattern{
		direction: OUTGOING,
		startNode: n,
		endNode:   other,
		types:     types,
	}
	return rel
}

// RelationshipFrom creates a relationship from another node to this one
func (n *nodePattern) RelationshipFrom(other NodeExpression, types ...string) RelationshipPattern {
	rel := &relationshipPattern{
		direction: INCOMING,
		startNode: other,
		endNode:   n,
		types:     types,
	}
	return rel
}

// RelationshipBetween creates an undirected relationship between this node and another
func (n *nodePattern) RelationshipBetween(other NodeExpression, types ...string) RelationshipPattern {
	rel := &relationshipPattern{
		direction: BIDIRECTIONAL,
		startNode: n,
		endNode:   other,
		types:     types,
	}
	return rel
}

// Accept implements the Expression interface
func (n *nodePattern) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(n)
}

// String returns a string representation of this node pattern
func (n *nodePattern) String() string {
	var builder strings.Builder

	builder.WriteString("(")
	if n.alias != "" {
		builder.WriteString(n.alias)
	}

	for _, label := range n.labels {
		builder.WriteString(":`")
		builder.WriteString(label)
		builder.WriteString("`")
	}

	if len(n.properties) > 0 {
		builder.WriteString(" {")
		propStrings := make([]string, 0, len(n.properties))
		for key, value := range n.properties {
			propStrings = append(propStrings, fmt.Sprintf("%s: %s", key, value.String()))
		}
		builder.WriteString(strings.Join(propStrings, ", "))
		builder.WriteString("}")
	}

	builder.WriteString(")")
	return builder.String()
}

// SymbolicName returns the alias of this node pattern
func (n *nodePattern) SymbolicName() string {
	return n.alias
}

// autoExpression automatically converts Go values to appropriate Expression types
func autoExpression(value interface{}) Expression {
	if value == nil {
		return Null()
	}

	switch v := value.(type) {
	case Expression:
		return v
	case string:
		return Literal(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return Literal(v)
	case []interface{}:
		elements := make([]Expression, len(v))
		for i, e := range v {
			elements[i] = autoExpression(e)
		}
		return ListOf(elements...)
	case map[string]interface{}:
		properties := make(map[string]Expression)
		for k, val := range v {
			properties[k] = autoExpression(val)
		}
		return Map(properties)
	default:
		// As a fallback, convert to a string representation
		return Literal(fmt.Sprintf("%v", v))
	}
}

// relationshipPattern represents a relationship pattern in Cypher (e.g., -[r:ACTED_IN]->)
type relationshipPattern struct {
	direction  Direction
	types      []string
	alias      string
	properties map[string]Expression
	minHops    *int
	maxHops    *int
	startNode  NodeExpression
	endNode    NodeExpression
}

// Direction constants are now defined in expressions.go

// Relationship creates a new relationship pattern with the given types
func Relationship(types ...string) *relationshipBuilder {
	return &relationshipBuilder{
		rel: &relationshipPattern{
			types:      types,
			properties: make(map[string]Expression),
		},
	}
}

// relationshipBuilder helps build relationship patterns
type relationshipBuilder struct {
	rel *relationshipPattern
}

// Direction sets the direction of this relationship pattern
func (r *relationshipBuilder) Direction(dir Direction) *relationshipBuilder {
	r.rel.direction = dir
	return r
}

// Named sets the alias for this relationship pattern
func (r *relationshipBuilder) Named(alias string) *relationshipBuilder {
	r.rel.alias = alias
	return r
}

// WithProperties adds properties to this relationship pattern
func (r *relationshipBuilder) WithProperties(properties map[string]Expression) *relationshipBuilder {
	for k, v := range properties {
		r.rel.properties[k] = v
	}
	return r
}

// WithProps adds properties with automatic conversion to expressions
func (r *relationshipBuilder) WithProps(properties map[string]interface{}) *relationshipBuilder {
	for k, v := range properties {
		r.rel.properties[k] = autoExpression(v)
	}
	return r
}

// Props is an alias for WithProps
func (r *relationshipBuilder) Props(properties map[string]interface{}) *relationshipBuilder {
	return r.WithProps(properties)
}

// Length sets the length range for this relationship pattern
func (r *relationshipBuilder) Length(min, max int) *relationshipBuilder {
	r.rel.minHops = &min
	r.rel.maxHops = &max
	return r
}

// Unbounded sets this relationship pattern to unbounded length
func (r *relationshipBuilder) Unbounded() *relationshipBuilder {
	// For unbounded, we don't set any limits
	r.rel.minHops = nil
	r.rel.maxHops = nil
	return r
}

// Between connects this relationship between two nodes
func (r *relationshipBuilder) Between(start NodeExpression, end NodeExpression) RelationshipPattern {
	r.rel.startNode = start
	r.rel.endNode = end
	return r.rel
}

// Build builds and returns the relationship pattern
func (r *relationshipBuilder) Build() RelationshipPattern {
	return r.rel
}

// Direction returns the direction of this relationship
func (r *relationshipPattern) Direction() Direction {
	return r.direction
}

// SetDirection sets the direction of this relationship
func (r *relationshipPattern) SetDirection(direction Direction) RelationshipPattern {
	clone := *r
	clone.direction = direction
	return &clone
}

// StartNode returns the start node of this relationship
func (r *relationshipPattern) StartNode() NodeExpression {
	return r.startNode
}

// EndNode returns the end node of this relationship
func (r *relationshipPattern) EndNode() NodeExpression {
	return r.endNode
}

// Named sets the alias for this relationship pattern
func (r *relationshipPattern) Named(alias string) RelationshipPattern {
	clone := *r
	clone.alias = alias
	return &clone
}

// WithType sets the type for this relationship pattern
func (r *relationshipPattern) WithType(relType string) RelationshipPattern {
	clone := *r
	clone.types = []string{relType}
	return &clone
}

// WithTypes sets the types for this relationship pattern
func (r *relationshipPattern) WithTypes(types ...string) RelationshipPattern {
	clone := *r
	clone.types = types
	return &clone
}

// WithProperties adds properties to this relationship pattern
func (r *relationshipPattern) WithProperties(properties map[string]Expression) RelationshipPattern {
	clone := *r
	if clone.properties == nil {
		clone.properties = make(map[string]Expression)
	}
	for k, v := range properties {
		clone.properties[k] = v
	}
	return &clone
}

// RelationshipTo continues the pattern with a new relationship
func (r *relationshipPattern) RelationshipTo(node NodeExpression, types ...string) PatternElement {
	// Implementation for chaining relationships
	newRel := &relationshipPattern{
		direction: OUTGOING,
		types:     types,
		startNode: r.endNode,
		endNode:   node,
	}

	return &patternPath{
		elements: []PatternElement{
			r.startNode.(PatternElement),
			r,
			r.endNode.(PatternElement),
			newRel,
			node.(PatternElement),
		},
	}
}

// RelationshipFrom continues the pattern with a new relationship
func (r *relationshipPattern) RelationshipFrom(node NodeExpression, types ...string) PatternElement {
	// Implementation for chaining relationships
	newRel := &relationshipPattern{
		direction: INCOMING,
		types:     types,
		startNode: node,
		endNode:   r.endNode,
	}

	return &patternPath{
		elements: []PatternElement{
			r.startNode.(PatternElement),
			r,
			r.endNode.(PatternElement),
			newRel,
			node.(PatternElement),
		},
	}
}

// Property returns a property access expression for this relationship
func (r *relationshipPattern) Property(propertyName string) PropertyExpression {
	return &propertyExpression{
		subject:      r,
		propertyName: propertyName,
	}
}

// Prop is an alias for Property
func (r *relationshipPattern) Prop(name string) PropertyExpression {
	return r.Property(name)
}

// WithProps adds properties with automatic conversion to expressions
func (r *relationshipPattern) WithProps(properties map[string]interface{}) Expression {
	clone := *r
	if clone.properties == nil {
		clone.properties = make(map[string]Expression)
	}
	for k, v := range properties {
		clone.properties[k] = autoExpression(v)
	}
	return &clone
}

// Props is an alias for WithProps
func (r *relationshipPattern) Props(properties map[string]interface{}) Expression {
	return r.WithProps(properties)
}

// Accept implements the Expression interface
func (r *relationshipPattern) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(r)
}

// String returns a string representation of this relationship pattern
func (r *relationshipPattern) String() string {
	var sb strings.Builder

	// Render the appropriate arrow start
	switch r.direction {
	case OUTGOING:
		sb.WriteString("-")
	case INCOMING:
		sb.WriteString("<-")
	case BIDIRECTIONAL:
		sb.WriteString("-")
	}

	// Render the relationship details
	sb.WriteString("[")
	if r.alias != "" {
		sb.WriteString(r.alias)
	}

	// Add relationship types
	for _, t := range r.types {
		sb.WriteString(":")
		sb.WriteString(t)
	}

	// Add hops if specified
	if r.minHops != nil || r.maxHops != nil {
		sb.WriteString("*")
		if r.minHops != nil {
			sb.WriteString(strconv.Itoa(*r.minHops))
		}
		sb.WriteString("..")
		if r.maxHops != nil {
			sb.WriteString(strconv.Itoa(*r.maxHops))
		}
	}

	// Add properties if any
	if len(r.properties) > 0 {
		sb.WriteString(" {")
		propStrings := make([]string, 0, len(r.properties))
		for key, value := range r.properties {
			propStrings = append(propStrings, fmt.Sprintf("%s: %s", key, value.String()))
		}
		sb.WriteString(strings.Join(propStrings, ", "))
		sb.WriteString("}")
	}

	sb.WriteString("]")

	// Render the appropriate arrow end
	switch r.direction {
	case OUTGOING:
		sb.WriteString("->")
	case INCOMING:
		sb.WriteString("-")
	case BIDIRECTIONAL:
		sb.WriteString("-")
	}

	return sb.String()
}

// SymbolicName returns the alias of this relationship pattern
func (r *relationshipPattern) SymbolicName() string {
	return r.alias
}

// patternPath represents a path pattern in Cypher (e.g., (a)-[r]->(b))
type patternPath struct {
	elements []PatternElement
	alias    string
}

// Pattern creates a new pattern path with the given elements
func Pattern(elements ...PatternElement) PatternElement {
	return &patternPath{
		elements: elements,
	}
}

// Named sets the alias for this pattern path
func (p *patternPath) Named(alias string) PatternElement {
	clone := *p
	clone.alias = alias
	return &clone
}

// Accept implements the Expression interface
func (p *patternPath) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(p)
}

// String returns a string representation of this pattern path
func (p *patternPath) String() string {
	var sb strings.Builder

	if p.alias != "" {
		sb.WriteString(p.alias)
		sb.WriteString(" = ")
	}

	for _, element := range p.elements {
		sb.WriteString(element.String())
	}

	return sb.String()
}

// SymbolicName returns the alias of this pattern path
func (p *patternPath) SymbolicName() string {
	return p.alias
}

// propertyExpression represents a property access expression in Cypher (e.g., n.name)
type propertyExpression struct {
	subject      Expression
	propertyName string
}

// PropertyName returns the name of the property being accessed
func (p *propertyExpression) PropertyName() string {
	return p.propertyName
}

// Subject returns the expression on which the property is accessed
func (p *propertyExpression) Subject() Expression {
	return p.subject
}

// Accept implements the Expression interface
func (p *propertyExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(p)
}

// String returns a string representation of this property expression
func (p *propertyExpression) String() string {
	return fmt.Sprintf("%s.%s", p.subject.String(), p.propertyName)
}

// Eq creates an equality comparison with another expression
func (p *propertyExpression) Eq(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    convertToExpression(other),
		operator: "=",
	}
}

// Ne creates an inequality comparison with another expression
func (p *propertyExpression) Ne(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    convertToExpression(other),
		operator: "<>",
	}
}

// Gt creates a greater-than comparison with another expression
func (p *propertyExpression) Gt(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    convertToExpression(other),
		operator: ">",
	}
}

// Gte creates a greater-than-or-equal comparison with another expression
func (p *propertyExpression) Gte(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    convertToExpression(other),
		operator: ">=",
	}
}

// Lt creates a less-than comparison with another expression
func (p *propertyExpression) Lt(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    convertToExpression(other),
		operator: "<",
	}
}

// Lte creates a less-than-or-equal comparison with another expression
func (p *propertyExpression) Lte(other interface{}) ComparisonExpression {
	return &comparisonExpression{
		left:     p,
		right:    convertToExpression(other),
		operator: "<=",
	}
}

// Contains creates a string contains operation
func (p *propertyExpression) Contains(value interface{}) BooleanExpression {
	return &stringOperation{
		left:     p,
		right:    convertToExpression(value),
		operator: "CONTAINS",
	}
}

// StartsWith creates a string starts-with operation
func (p *propertyExpression) StartsWith(value interface{}) BooleanExpression {
	return &stringOperation{
		left:     p,
		right:    convertToExpression(value),
		operator: "STARTS WITH",
	}
}

// EndsWith creates a string ends-with operation
func (p *propertyExpression) EndsWith(value interface{}) BooleanExpression {
	return &stringOperation{
		left:     p,
		right:    convertToExpression(value),
		operator: "ENDS WITH",
	}
}

// In creates a collection contains operation
func (p *propertyExpression) In(list interface{}) BooleanExpression {
	return &comparisonExpression{
		left:     p,
		right:    convertToExpression(list),
		operator: "IN",
	}
}

// IsNull creates a null check operation
func (p *propertyExpression) IsNull() BooleanExpression {
	return &comparisonExpression{
		left:     p,
		right:    Null(),
		operator: "IS",
	}
}

// IsNotNull creates a not-null check operation
func (p *propertyExpression) IsNotNull() BooleanExpression {
	return &comparisonExpression{
		left:     p,
		right:    Null(),
		operator: "IS NOT",
	}
}

// Aliases for improved readability

// Equals is an alias for Eq
func (p *propertyExpression) Equals(other interface{}) ComparisonExpression {
	return p.Eq(other)
}

// NotEquals is an alias for Ne
func (p *propertyExpression) NotEquals(other interface{}) ComparisonExpression {
	return p.Ne(other)
}

// GreaterThan is an alias for Gt
func (p *propertyExpression) GreaterThan(other interface{}) ComparisonExpression {
	return p.Gt(other)
}

// GreaterThanOrEqual is an alias for Gte
func (p *propertyExpression) GreaterThanOrEqual(other interface{}) ComparisonExpression {
	return p.Gte(other)
}

// LessThan is an alias for Lt
func (p *propertyExpression) LessThan(other interface{}) ComparisonExpression {
	return p.Lt(other)
}

// LessThanOrEqual is an alias for Lte
func (p *propertyExpression) LessThanOrEqual(other interface{}) ComparisonExpression {
	return p.Lte(other)
}

// Helper function to convert various types to Expression
func convertToExpression(value interface{}) Expression {
	if value == nil {
		return Null()
	}

	// If value is already an Expression, return it
	if expr, ok := value.(Expression); ok {
		return expr
	}

	// Otherwise, convert to a literal
	return Literal(value)
}

// stringOperation represents a string operation in Cypher (e.g., a.name CONTAINS 'Tom')
type stringOperation struct {
	left     Expression
	right    Expression
	operator string
}

// Accept implements the Expression interface
func (s *stringOperation) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(s)
}

// String returns a string representation of this string operation
func (s *stringOperation) String() string {
	return fmt.Sprintf("%s %s %s", s.left.String(), s.operator, s.right.String())
}

// And combines this expression with another using AND
func (s *stringOperation) And(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     s,
		right:    other,
		operator: "AND",
	}
}

// Or combines this expression with another using OR
func (s *stringOperation) Or(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     s,
		right:    other,
		operator: "OR",
	}
}

// Not negates this boolean expression
func (s *stringOperation) Not() BooleanExpression {
	return &notExpression{expression: s}
}

// IsBooleanExpression marks this as a boolean expression
func (s *stringOperation) IsBooleanExpression() {}

// nullCheck represents a NULL check in Cypher (e.g., a.name IS NULL)
type nullCheck struct {
	expression Expression
	isNull     bool
}

// Accept implements the Expression interface
func (n *nullCheck) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(n)
}

// String returns a string representation of this NULL check
func (n *nullCheck) String() string {
	if n.isNull {
		return fmt.Sprintf("%s IS NULL", n.expression.String())
	}
	return fmt.Sprintf("%s IS NOT NULL", n.expression.String())
}

// And combines this expression with another using AND
func (n *nullCheck) And(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     n,
		right:    other,
		operator: "AND",
	}
}

// Or combines this expression with another using OR
func (n *nullCheck) Or(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     n,
		right:    other,
		operator: "OR",
	}
}

// Not negates this boolean expression
func (n *nullCheck) Not() BooleanExpression {
	return &notExpression{expression: n}
}

// IsBooleanExpression marks this as a boolean expression
func (n *nullCheck) IsBooleanExpression() {}

// comparisonExpression represents a comparison in Cypher (e.g., a.age > 30)
type comparisonExpression struct {
	left     Expression
	right    Expression
	operator string
}

// Left returns the left side of the comparison
func (c *comparisonExpression) Left() Expression {
	return c.left
}

// Right returns the right side of the comparison
func (c *comparisonExpression) Right() Expression {
	return c.right
}

// Operator returns the comparison operator
func (c *comparisonExpression) Operator() string {
	return c.operator
}

// Accept implements the Expression interface
func (c *comparisonExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(c)
}

// String returns a string representation of this comparison
func (c *comparisonExpression) String() string {
	return fmt.Sprintf("%s %s %s", c.left.String(), c.operator, c.right.String())
}

// And combines this expression with another using AND
func (c *comparisonExpression) And(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     c,
		right:    other,
		operator: "AND",
	}
}

// Or combines this expression with another using OR
func (c *comparisonExpression) Or(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     c,
		right:    other,
		operator: "OR",
	}
}

// Not negates this comparison expression
func (c *comparisonExpression) Not() BooleanExpression {
	return &notExpression{
		expression: c,
	}
}

// IsBooleanExpression marks this as a boolean expression
func (c *comparisonExpression) IsBooleanExpression() {}

// logicalExpression represents a logical operation in Cypher (e.g., a AND b)
type logicalExpression struct {
	left     BooleanExpression
	right    BooleanExpression
	operator string
}

// Accept implements the Expression interface
func (l *logicalExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(l)
}

// String returns a string representation of this logical operation
func (l *logicalExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", l.left.String(), l.operator, l.right.String())
}

// And combines this expression with another using AND
func (l *logicalExpression) And(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     l,
		right:    other,
		operator: "AND",
	}
}

// Or combines this expression with another using OR
func (l *logicalExpression) Or(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     l,
		right:    other,
		operator: "OR",
	}
}

// Not negates this boolean expression
func (l *logicalExpression) Not() BooleanExpression {
	return &notExpression{expression: l}
}

// IsBooleanExpression marks this as a boolean expression
func (l *logicalExpression) IsBooleanExpression() {}

// notExpression represents a negation in Cypher (e.g., NOT a)
type notExpression struct {
	expression Expression
}

// Accept implements the Expression interface
func (n *notExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(n)
}

// String returns a string representation of this negation
func (n *notExpression) String() string {
	if n.expression == nil {
		return "NOT ()"
	}
	return fmt.Sprintf("NOT (%s)", n.expression.String())
}

// And combines this expression with another using AND
func (n *notExpression) And(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     n,
		right:    other,
		operator: "AND",
	}
}

// Or combines this expression with another using OR
func (n *notExpression) Or(other BooleanExpression) BooleanExpression {
	return &logicalExpression{
		left:     n,
		right:    other,
		operator: "OR",
	}
}

// Not negates this boolean expression
func (n *notExpression) Not() BooleanExpression {
	// If the inner expression can be cast to BooleanExpression, return it directly
	if boolExpr, ok := n.expression.(BooleanExpression); ok {
		return boolExpr
	}
	// Otherwise, wrap it in a notExpression again
	return &notExpression{expression: n.expression}
}

// IsBooleanExpression marks this as a boolean expression
func (n *notExpression) IsBooleanExpression() {}

// aliasedExpression represents an expression with an alias in Cypher (e.g., person AS p)
type aliasedExpression struct {
	expression Expression
	aliasName  string
}

// Accept implements the Expression interface
func (a *aliasedExpression) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(a)
}

// String returns a string representation of this aliased expression
func (a *aliasedExpression) String() string {
	return fmt.Sprintf("%s AS %s", a.expression.String(), a.aliasName)
}

// Inner returns the underlying expression
func (a *aliasedExpression) Inner() Expression {
	return a.expression
}

// Alias returns the alias for this expression
func (a *aliasedExpression) Alias() string {
	return a.aliasName
}

// As creates an aliased expression
func As(expression Expression, alias string) AliasedExpression {
	return &aliasedExpression{
		expression: expression,
		aliasName:  alias,
	}
}
