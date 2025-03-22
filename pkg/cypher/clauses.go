package cypher

import (
	"fmt"
)

// ClauseBuilder is the base interface for all clause builders
type ClauseBuilder interface {
	Buildable
	// Error returns any error that occurred
	Error() error
	// HasError checks if an error occurred
	HasError() bool
}

// Clause represents a Cypher clause
type Clause interface {
	Expression
	// AsString returns the clause as a string
	AsString() string
}

// MatchClause represents a MATCH clause
type MatchClause struct {
	patterns []Expression
	optional bool
}

// Match creates a new MATCH clause with the given patterns
func Match(patterns ...Expression) *MatchBuilder {
	return &MatchBuilder{
		clause: &MatchClause{
			patterns: patterns,
		},
	}
}

// OptionalMatch creates a new OPTIONAL MATCH clause with the given patterns
func OptionalMatch(patterns ...Expression) *MatchBuilder {
	return &MatchBuilder{
		clause: &MatchClause{
			patterns: patterns,
			optional: true,
		},
	}
}

// Accept implements the Expression interface
func (m *MatchClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(m)
}

// String returns a string representation of this MATCH clause
func (m *MatchClause) String() string {
	return m.AsString()
}

// AsString returns the clause as a string
func (m *MatchClause) AsString() string {
	if m.optional {
		return "OPTIONAL MATCH"
	}
	return "MATCH"
}

// Patterns returns the patterns in this MATCH clause
func (m *MatchClause) Patterns() []Expression {
	return m.patterns
}

// IsOptional returns whether this is an OPTIONAL MATCH clause
func (m *MatchClause) IsOptional() bool {
	return m.optional
}

// MatchBuilder builds a query starting with a MATCH clause
type MatchBuilder struct {
	prev   ClauseBuilder // Previous builder in the chain
	clause *MatchClause
	where  *WhereClause
	Return *ReturnClause
	with   *WithClause
	err    error // Error during building
}

// Error returns the error that occurred during building
func (m *MatchBuilder) Error() error {
	return m.err
}

// HasError checks if an error occurred during building
func (m *MatchBuilder) HasError() bool {
	return m.err != nil
}

// Where adds a WHERE clause to this builder
func (m *MatchBuilder) Where(condition BooleanExpression) *MatchBuilder {
	if m.err != nil {
		return m
	}

	if condition == nil {
		m.err = fmt.Errorf("where condition cannot be nil")
		return m
	}

	m.where = &WhereClause{
		condition: condition,
	}
	return m
}

// With adds a WITH clause to this builder
func (m *MatchBuilder) With(expressions ...Expression) *WithBuilder {
	if m.err != nil {
		return &WithBuilder{prev: m, err: m.err}
	}

	if len(expressions) == 0 {
		m.err = fmt.Errorf("with clause must have at least one expression")
		return &WithBuilder{prev: m, err: m.err}
	}

	withClause := &WithClause{
		expressions: expressions,
	}
	builder := &WithBuilder{
		prev:   m,
		clause: withClause,
	}
	m.with = withClause
	return builder
}

// Returning adds a RETURN clause to this builder
func (m *MatchBuilder) Returning(expressions ...Expression) *ReturnBuilder {
	if m.err != nil {
		return &ReturnBuilder{prev: m, err: m.err}
	}

	if len(expressions) == 0 {
		m.err = fmt.Errorf("return clause must have at least one expression")
		return &ReturnBuilder{prev: m, err: m.err}
	}

	returnClause := &ReturnClause{
		expressions: expressions,
	}
	builder := &ReturnBuilder{
		prev:   m,
		clause: returnClause,
	}
	m.Return = returnClause
	return builder
}

// WhereClause represents a WHERE clause
type WhereClause struct {
	condition BooleanExpression
}

// Accept implements the Expression interface
func (w *WhereClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(w)
}

// String returns a string representation of this WHERE clause
func (w *WhereClause) String() string {
	return w.AsString()
}

// AsString returns the clause as a string
func (w *WhereClause) AsString() string {
	return "WHERE"
}

// Condition returns the condition in this WHERE clause
func (w *WhereClause) Condition() BooleanExpression {
	return w.condition
}

// ReturnClause represents a RETURN clause
type ReturnClause struct {
	expressions []Expression
	distinct    bool
	limit       Expression
	skip        Expression
	orderBy     []*OrderByItem
}

// Distinct marks this RETURN clause as DISTINCT
func (r *ReturnClause) Distinct() *ReturnClause {
	r.distinct = true
	return r
}

// OrderByItem represents an ORDER BY item
type OrderByItem struct {
	expression Expression
	descending bool
}

// Accept implements the Expression interface
func (r *ReturnClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(r)
}

// String returns a string representation of this RETURN clause
func (r *ReturnClause) String() string {
	return r.AsString()
}

// AsString returns the clause as a string
func (r *ReturnClause) AsString() string {
	if r.distinct {
		return "RETURN DISTINCT"
	}
	return "RETURN"
}

// Expressions returns the expressions in this RETURN clause
func (r *ReturnClause) Expressions() []Expression {
	return r.expressions
}

// IsDistinct returns whether this is a DISTINCT RETURN clause
func (r *ReturnClause) IsDistinct() bool {
	return r.distinct
}

// Limit returns the limit expression for this RETURN clause
func (r *ReturnClause) Limit() Expression {
	return r.limit
}

// Skip returns the skip expression for this RETURN clause
func (r *ReturnClause) Skip() Expression {
	return r.skip
}

// OrderBy returns the ORDER BY items for this RETURN clause
func (r *ReturnClause) OrderBy() []*OrderByItem {
	return r.orderBy
}

// ReturnBuilder builds a query ending with a RETURN clause
type ReturnBuilder struct {
	prev   ClauseBuilder
	clause *ReturnClause
	err    error // Error during building
}

// Error returns the error that occurred during building
func (r *ReturnBuilder) Error() error {
	if r.prev != nil {
		if prevWithErr, ok := r.prev.(interface{ Error() error }); ok && prevWithErr.Error() != nil {
			return prevWithErr.Error()
		}
	}
	return r.err
}

// HasError checks if an error occurred during building
func (r *ReturnBuilder) HasError() bool {
	if r.prev != nil {
		if prevWithErr, ok := r.prev.(interface{ HasError() bool }); ok && prevWithErr.HasError() {
			return true
		}
	}
	return r.err != nil
}

// Distinct marks this RETURN clause as DISTINCT
func (r *ReturnBuilder) Distinct() *ReturnBuilder {
	r.clause.distinct = true
	return r
}

// OrderBy adds an ORDER BY clause to this builder
func (r *ReturnBuilder) OrderBy(expressions ...Expression) *ReturnBuilder {
	for _, expr := range expressions {
		r.clause.orderBy = append(r.clause.orderBy, &OrderByItem{
			expression: expr,
			descending: false,
		})
	}
	return r
}

// OrderByDesc adds a descending ORDER BY clause to this builder
func (r *ReturnBuilder) OrderByDesc(expressions ...Expression) *ReturnBuilder {
	for _, expr := range expressions {
		r.clause.orderBy = append(r.clause.orderBy, &OrderByItem{
			expression: expr,
			descending: true,
		})
	}
	return r
}

// Skip adds a SKIP clause to this builder
func (r *ReturnBuilder) Skip(count Expression) *ReturnBuilder {
	r.clause.skip = count
	return r
}

// Limit adds a LIMIT clause to this builder
func (r *ReturnBuilder) Limit(count Expression) *ReturnBuilder {
	r.clause.limit = count
	return r
}

// Build constructs the final Statement
func (r *ReturnBuilder) Build() (Statement, error) {
	if r.HasError() {
		return nil, r.Error()
	}

	return &cypherStatement{
		builder: r,
	}, nil
}

// WithClause represents a WITH clause
type WithClause struct {
	expressions []Expression
	distinct    bool
	where       *WhereClause
	orderBy     []*OrderByItem
	limit       Expression
	skip        Expression
}

// Accept implements the Expression interface
func (w *WithClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(w)
}

// String returns a string representation of this WITH clause
func (w *WithClause) String() string {
	return w.AsString()
}

// AsString returns the clause as a string
func (w *WithClause) AsString() string {
	if w.distinct {
		return "WITH DISTINCT"
	}
	return "WITH"
}

// Expressions returns the expressions in this WITH clause
func (w *WithClause) Expressions() []Expression {
	return w.expressions
}

// IsDistinct returns whether this is a DISTINCT WITH clause
func (w *WithClause) IsDistinct() bool {
	return w.distinct
}

// Where returns the WHERE clause for this WITH clause
func (w *WithClause) Where() *WhereClause {
	return w.where
}

// OrderBy returns the ORDER BY items for this WITH clause
func (w *WithClause) OrderBy() []*OrderByItem {
	return w.orderBy
}

// Limit returns the limit expression for this WITH clause
func (w *WithClause) Limit() Expression {
	return w.limit
}

// Skip returns the skip expression for this WITH clause
func (w *WithClause) Skip() Expression {
	return w.skip
}

// WithBuilder builds a query continuing with a WITH clause
type WithBuilder struct {
	prev   ClauseBuilder
	clause *WithClause
	err    error // Error during building
}

// Error returns the error that occurred during building
func (w *WithBuilder) Error() error {
	if w.prev != nil {
		if prevWithErr, ok := w.prev.(interface{ Error() error }); ok && prevWithErr.Error() != nil {
			return prevWithErr.Error()
		}
	}
	return w.err
}

// HasError checks if an error occurred during building
func (w *WithBuilder) HasError() bool {
	if w.prev != nil {
		if prevWithErr, ok := w.prev.(interface{ HasError() bool }); ok && prevWithErr.HasError() {
			return true
		}
	}
	return w.err != nil
}

// Where adds a WHERE clause to this builder
func (w *WithBuilder) Where(condition BooleanExpression) *WithBuilder {
	if w.HasError() {
		return w
	}

	if condition == nil {
		w.err = fmt.Errorf("where condition cannot be nil")
		return w
	}

	w.clause.where = &WhereClause{
		condition: condition,
	}
	return w
}

// OrderBy adds an ORDER BY clause to this builder
func (w *WithBuilder) OrderBy(expressions ...Expression) *WithBuilder {
	for _, expr := range expressions {
		w.clause.orderBy = append(w.clause.orderBy, &OrderByItem{
			expression: expr,
			descending: false,
		})
	}
	return w
}

// OrderByDesc adds a descending ORDER BY clause to this builder
func (w *WithBuilder) OrderByDesc(expressions ...Expression) *WithBuilder {
	for _, expr := range expressions {
		w.clause.orderBy = append(w.clause.orderBy, &OrderByItem{
			expression: expr,
			descending: true,
		})
	}
	return w
}

// Skip adds a SKIP clause to this builder
func (w *WithBuilder) Skip(count Expression) *WithBuilder {
	w.clause.skip = count
	return w
}

// Limit adds a LIMIT clause to this builder
func (w *WithBuilder) Limit(count Expression) *WithBuilder {
	w.clause.limit = count
	return w
}

// Distinct marks this WITH clause as DISTINCT
func (w *WithBuilder) Distinct() *WithBuilder {
	w.clause.distinct = true
	return w
}

// Match adds a MATCH clause after this WITH clause
func (w *WithBuilder) Match(patterns ...Expression) *MatchBuilder {
	if w.HasError() {
		return &MatchBuilder{err: w.Error()}
	}

	if len(patterns) == 0 {
		w.err = fmt.Errorf("match clause must have at least one pattern")
		return &MatchBuilder{err: w.Error()}
	}

	matchBuilder := &MatchBuilder{
		clause: &MatchClause{
			patterns: patterns,
		},
	}

	// Update the chain to properly connect WITH and MATCH
	matchBuilder.prev = w

	return matchBuilder
}

// OptionalMatch adds an OPTIONAL MATCH clause after this WITH clause
func (w *WithBuilder) OptionalMatch(patterns ...Expression) *MatchBuilder {
	if w.HasError() {
		return &MatchBuilder{err: w.Error()}
	}

	if len(patterns) == 0 {
		w.err = fmt.Errorf("optional match clause must have at least one pattern")
		return &MatchBuilder{err: w.Error()}
	}

	matchBuilder := &MatchBuilder{
		clause: &MatchClause{
			patterns: patterns,
			optional: true,
		},
	}

	// Update the chain to properly connect WITH and MATCH
	matchBuilder.prev = w

	return matchBuilder
}

// Returning adds a RETURN clause after this WITH clause
func (w *WithBuilder) Returning(expressions ...Expression) *ReturnBuilder {
	if w.HasError() {
		return &ReturnBuilder{err: w.Error()}
	}

	if len(expressions) == 0 {
		w.err = fmt.Errorf("return clause must have at least one expression")
		return &ReturnBuilder{err: w.Error()}
	}

	returnClause := &ReturnClause{
		expressions: expressions,
	}
	builder := &ReturnBuilder{
		prev:   w,
		clause: returnClause,
	}
	return builder
}

// CreateClause represents a CREATE clause
type CreateClause struct {
	patterns []Expression
}

// Create creates a new CREATE clause with the given patterns
func Create(patterns ...Expression) *CreateBuilder {
	return &CreateBuilder{
		clause: &CreateClause{
			patterns: patterns,
		},
	}
}

// Accept implements the Expression interface
func (c *CreateClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(c)
}

// String returns a string representation of this CREATE clause
func (c *CreateClause) String() string {
	return c.AsString()
}

// AsString returns the clause as a string
func (c *CreateClause) AsString() string {
	return "CREATE"
}

// Patterns returns the patterns in this CREATE clause
func (c *CreateClause) Patterns() []Expression {
	return c.patterns
}

// CreateBuilder builds a query starting with a CREATE clause
type CreateBuilder struct {
	clause *CreateClause
	Return *ReturnClause
	prev   ClauseBuilder
	err    error // Error during building
}

// Error returns the error that occurred during building
func (c *CreateBuilder) Error() error {
	if c.prev != nil {
		if prevWithErr, ok := c.prev.(interface{ Error() error }); ok && prevWithErr.Error() != nil {
			return prevWithErr.Error()
		}
	}
	return c.err
}

// HasError checks if an error occurred during building
func (c *CreateBuilder) HasError() bool {
	if c.prev != nil {
		if prevWithErr, ok := c.prev.(interface{ HasError() bool }); ok && prevWithErr.HasError() {
			return true
		}
	}
	return c.err != nil
}

// Returning adds a RETURN clause to this builder
func (c *CreateBuilder) Returning(expressions ...Expression) *ReturnBuilder {
	if c.HasError() {
		return &ReturnBuilder{err: c.Error()}
	}

	if len(expressions) == 0 {
		c.err = fmt.Errorf("return clause must have at least one expression")
		return &ReturnBuilder{err: c.Error()}
	}

	returnClause := &ReturnClause{
		expressions: expressions,
	}
	builder := &ReturnBuilder{
		prev:   c,
		clause: returnClause,
	}
	c.Return = returnClause
	return builder
}

// Build constructs the final Statement
func (c *CreateBuilder) Build() (Statement, error) {
	if c.HasError() {
		return nil, c.Error()
	}

	return &cypherStatement{
		builder: c,
	}, nil
}

// MergeClause represents a MERGE clause
type MergeClause struct {
	pattern Expression
}

// Merge creates a new MERGE clause with the given pattern
func Merge(pattern Expression) *MergeBuilder {
	return &MergeBuilder{
		clause: &MergeClause{
			pattern: pattern,
		},
	}
}

// Accept implements the Expression interface
func (m *MergeClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(m)
}

// String returns a string representation of this MERGE clause
func (m *MergeClause) String() string {
	return m.AsString()
}

// AsString returns the clause as a string
func (m *MergeClause) AsString() string {
	return "MERGE"
}

// Pattern returns the pattern in this MERGE clause
func (m *MergeClause) Pattern() Expression {
	return m.pattern
}

// MergeBuilder builds a query starting with a MERGE clause
type MergeBuilder struct {
	clause *MergeClause
	Return *ReturnClause
	err    error // Error during building
}

// Error returns the error that occurred during building
func (m *MergeBuilder) Error() error {
	return m.err
}

// HasError checks if an error occurred during building
func (m *MergeBuilder) HasError() bool {
	return m.err != nil
}

// Returning adds a RETURN clause to this builder
func (m *MergeBuilder) Returning(expressions ...Expression) *ReturnBuilder {
	if m.HasError() {
		return &ReturnBuilder{err: m.Error()}
	}

	if len(expressions) == 0 {
		m.err = fmt.Errorf("return clause must have at least one expression")
		return &ReturnBuilder{err: m.Error()}
	}

	returnClause := &ReturnClause{
		expressions: expressions,
	}
	builder := &ReturnBuilder{
		prev:   m,
		clause: returnClause,
	}
	m.Return = returnClause
	return builder
}

// Build constructs the final Statement
func (m *MergeBuilder) Build() (Statement, error) {
	if m.HasError() {
		return nil, m.Error()
	}

	return &cypherStatement{
		builder: m,
	}, nil
}

// DeleteClause represents a DELETE clause
type DeleteClause struct {
	expressions []Expression
	detach      bool
}

// Delete creates a new DELETE clause with the given expressions
func Delete(expressions ...Expression) *DeleteBuilder {
	return &DeleteBuilder{
		clause: &DeleteClause{
			expressions: expressions,
		},
	}
}

// DetachDelete creates a new DETACH DELETE clause with the given expressions
func DetachDelete(expressions ...Expression) *DeleteBuilder {
	return &DeleteBuilder{
		clause: &DeleteClause{
			expressions: expressions,
			detach:      true,
		},
	}
}

// Accept implements the Expression interface
func (d *DeleteClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(d)
}

// String returns a string representation of this DELETE clause
func (d *DeleteClause) String() string {
	return d.AsString()
}

// AsString returns the clause as a string
func (d *DeleteClause) AsString() string {
	if d.detach {
		return "DETACH DELETE"
	}
	return "DELETE"
}

// Expressions returns the expressions in this DELETE clause
func (d *DeleteClause) Expressions() []Expression {
	return d.expressions
}

// IsDetach returns whether this is a DETACH DELETE clause
func (d *DeleteClause) IsDetach() bool {
	return d.detach
}

// DeleteBuilder builds a query ending with a DELETE clause
type DeleteBuilder struct {
	clause *DeleteClause
	err    error // Error during building
}

// Error returns the error that occurred during building
func (d *DeleteBuilder) Error() error {
	return d.err
}

// HasError checks if an error occurred during building
func (d *DeleteBuilder) HasError() bool {
	return d.err != nil
}

// Build constructs the final Statement
func (d *DeleteBuilder) Build() (Statement, error) {
	if d.HasError() {
		return nil, d.Error()
	}

	return &cypherStatement{
		builder: d,
	}, nil
}

// SetClause represents a SET clause
type SetClause struct {
	items []*SetItem
}

// SetItem represents a SET item
type SetItem struct {
	property  Expression
	value     Expression
	labels    []string
	nodeAlias string
	operator  string // "=" for assignment, "+=" for add, "-=" for remove
}

// Accept implements the Expression interface
func (s *SetClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(s)
}

// String returns a string representation of this SET clause
func (s *SetClause) String() string {
	return s.AsString()
}

// AsString returns the clause as a string
func (s *SetClause) AsString() string {
	return "SET"
}

// Items returns the items in this SET clause
func (s *SetClause) Items() []*SetItem {
	return s.items
}

// Set creates a new SET clause
func Set() *SetBuilder {
	return &SetBuilder{
		clause: &SetClause{},
	}
}

// SetProperty adds a property assignment to a SET clause
func SetProperty(property PropertyExpression, value Expression) *SetBuilder {
	return Set().Property(property, value)
}

// SetLabels adds labels to a node in a SET clause
func SetLabels(nodeAlias string, labels ...string) *SetBuilder {
	return Set().Labels(nodeAlias, labels...)
}

// SetBuilder builds a query with a SET clause
type SetBuilder struct {
	clause *SetClause
	err    error // Error during building
}

// Error returns the error that occurred during building
func (s *SetBuilder) Error() error {
	return s.err
}

// HasError checks if an error occurred during building
func (s *SetBuilder) HasError() bool {
	return s.err != nil
}

// Property adds a property assignment to this SET clause
func (s *SetBuilder) Property(property PropertyExpression, value Expression) *SetBuilder {
	s.clause.items = append(s.clause.items, &SetItem{
		property: property,
		value:    value,
		operator: "=",
	})
	return s
}

// PropertyAdd adds a property addition to this SET clause
func (s *SetBuilder) PropertyAdd(property PropertyExpression, value Expression) *SetBuilder {
	s.clause.items = append(s.clause.items, &SetItem{
		property: property,
		value:    value,
		operator: "+=",
	})
	return s
}

// Labels adds labels to a node in this SET clause
func (s *SetBuilder) Labels(nodeAlias string, labels ...string) *SetBuilder {
	s.clause.items = append(s.clause.items, &SetItem{
		nodeAlias: nodeAlias,
		labels:    labels,
	})
	return s
}

// Build constructs the final Statement
func (s *SetBuilder) Build() (Statement, error) {
	if s.HasError() {
		return nil, s.Error()
	}

	return &cypherStatement{
		builder: s,
	}, nil
}

// RemoveClause represents a REMOVE clause
type RemoveClause struct {
	items []*RemoveItem
}

// RemoveItem represents a REMOVE item
type RemoveItem struct {
	property  Expression
	nodeAlias string
	label     string
	isLabel   bool
}

// Accept implements the Expression interface
func (r *RemoveClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(r)
}

// String returns a string representation of this REMOVE clause
func (r *RemoveClause) String() string {
	return r.AsString()
}

// AsString returns the clause as a string
func (r *RemoveClause) AsString() string {
	return "REMOVE"
}

// Items returns the items in this REMOVE clause
func (r *RemoveClause) Items() []*RemoveItem {
	return r.items
}

// Remove creates a new REMOVE clause
func Remove() *RemoveBuilder {
	return &RemoveBuilder{
		clause: &RemoveClause{},
	}
}

// RemoveProperty removes a property in a REMOVE clause
func RemoveProperty(property PropertyExpression) *RemoveBuilder {
	return Remove().Property(property)
}

// RemoveLabel removes a label from a node in a REMOVE clause
func RemoveLabel(nodeAlias string, label string) *RemoveBuilder {
	return Remove().Label(nodeAlias, label)
}

// RemoveBuilder builds a query with a REMOVE clause
type RemoveBuilder struct {
	clause *RemoveClause
	err    error // Error during building
}

// Error returns the error that occurred during building
func (r *RemoveBuilder) Error() error {
	return r.err
}

// HasError checks if an error occurred during building
func (r *RemoveBuilder) HasError() bool {
	return r.err != nil
}

// Property adds a property removal to this REMOVE clause
func (r *RemoveBuilder) Property(property PropertyExpression) *RemoveBuilder {
	r.clause.items = append(r.clause.items, &RemoveItem{
		property: property,
		isLabel:  false,
	})
	return r
}

// Label adds a label removal to this REMOVE clause
func (r *RemoveBuilder) Label(nodeAlias string, label string) *RemoveBuilder {
	r.clause.items = append(r.clause.items, &RemoveItem{
		nodeAlias: nodeAlias,
		label:     label,
		isLabel:   true,
	})
	return r
}

// Build constructs the final Statement
func (r *RemoveBuilder) Build() (Statement, error) {
	if r.HasError() {
		return nil, r.Error()
	}

	return &cypherStatement{
		builder: r,
	}, nil
}

// UnwindClause represents an UNWIND clause
type UnwindClause struct {
	expression Expression
	alias      string
}

// Accept implements the Expression interface
func (u *UnwindClause) Accept(visitor ExpressionVisitor) any {
	return visitor.Visit(u)
}

// String returns a string representation of this UNWIND clause
func (u *UnwindClause) String() string {
	return u.AsString()
}

// AsString returns the clause as a string
func (u *UnwindClause) AsString() string {
	return "UNWIND"
}

// Expression returns the expression in this UNWIND clause
func (u *UnwindClause) Expression() Expression {
	return u.expression
}

// Alias returns the alias for the unwound elements
func (u *UnwindClause) Alias() string {
	return u.alias
}

// Unwind creates a new UNWIND clause
func Unwind(expression Expression, alias string) *UnwindBuilder {
	return &UnwindBuilder{
		clause: &UnwindClause{
			expression: expression,
			alias:      alias,
		},
	}
}

// UnwindBuilder builds a query starting with an UNWIND clause
type UnwindBuilder struct {
	clause *UnwindClause
	err    error // Error during building
}

// Error returns the error that occurred during building
func (u *UnwindBuilder) Error() error {
	return u.err
}

// HasError checks if an error occurred during building
func (u *UnwindBuilder) HasError() bool {
	return u.err != nil
}

// Match adds a MATCH clause after this UNWIND clause
func (u *UnwindBuilder) Match(patterns ...Expression) *MatchBuilder {
	if u.HasError() {
		return &MatchBuilder{err: u.Error()}
	}

	if len(patterns) == 0 {
		u.err = fmt.Errorf("match clause must have at least one pattern")
		return &MatchBuilder{err: u.Error()}
	}

	matchBuilder := &MatchBuilder{
		clause: &MatchClause{
			patterns: patterns,
		},
	}

	// Connect UNWIND to MATCH
	matchBuilder.prev = u

	return matchBuilder
}

// Returning adds a RETURN clause after this UNWIND clause
func (u *UnwindBuilder) Returning(expressions ...Expression) *ReturnBuilder {
	if u.HasError() {
		return &ReturnBuilder{err: u.Error()}
	}

	if len(expressions) == 0 {
		u.err = fmt.Errorf("return clause must have at least one expression")
		return &ReturnBuilder{err: u.Error()}
	}

	returnClause := &ReturnClause{
		expressions: expressions,
	}
	builder := &ReturnBuilder{
		prev:   u,
		clause: returnClause,
	}
	return builder
}

// With adds a WITH clause after this UNWIND clause
func (u *UnwindBuilder) With(expressions ...Expression) *WithBuilder {
	if u.HasError() {
		return &WithBuilder{err: u.Error()}
	}

	if len(expressions) == 0 {
		u.err = fmt.Errorf("with clause must have at least one expression")
		return &WithBuilder{err: u.Error()}
	}

	withClause := &WithClause{
		expressions: expressions,
	}
	builder := &WithBuilder{
		prev:   u,
		clause: withClause,
	}
	return builder
}
