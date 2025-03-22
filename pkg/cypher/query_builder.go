package cypher

import (
	"errors"
	"fmt"
)

// CypherBuilder is the main entry point for building Cypher queries
// It includes error handling capabilities and simplified syntax
type CypherBuilder struct {
	statement Statement
	lastError error
}

// Builder creates a new CypherBuilder with the given statement
func Builder(statement Statement) *CypherBuilder {
	return &CypherBuilder{statement: statement}
}

// Error returns the last error that occurred during building
func (b *CypherBuilder) Error() error {
	return b.lastError
}

// HasError checks if an error occurred during building
func (b *CypherBuilder) HasError() bool {
	return b.lastError != nil
}

// Build finalizes the query and returns the statement and error
func (b *CypherBuilder) Build() (Statement, error) {
	if b.lastError != nil {
		return nil, b.lastError
	}
	return b.statement, nil
}

// BuildE finalizes the query and returns the statement, failing if an error occurred
func (b *CypherBuilder) BuildE() Statement {
	if b.lastError != nil {
		panic(b.lastError)
	}
	return b.statement
}

// MustBuild is an alias for BuildE for users who prefer that terminology
func (b *CypherBuilder) MustBuild() Statement {
	return b.BuildE()
}

// Cypher returns the Cypher string representation of the statement
// If an error occurred during building, it returns the error message
func (b *CypherBuilder) Cypher() string {
	if b.lastError != nil {
		return fmt.Sprintf("ERROR: %s", b.lastError.Error())
	}
	return b.statement.Cypher()
}

// Params returns the parameters of the statement
// If an error occurred during building, it returns nil
func (b *CypherBuilder) Params() map[string]interface{} {
	if b.lastError != nil {
		return nil
	}
	return b.statement.Params()
}

// Query starts a new query with the given statement
func Query(statement Statement) *CypherBuilder {
	if statement == nil {
		return &CypherBuilder{lastError: errors.New("statement cannot be nil")}
	}
	return &CypherBuilder{statement: statement}
}

// StatementsBuilder builds multiple statements as part of a transaction
type StatementsBuilder struct {
	statements []Statement
	lastError  error
}

// NewTransaction creates a new StatementsBuilder
func NewTransaction() *StatementsBuilder {
	return &StatementsBuilder{statements: make([]Statement, 0)}
}

// Add adds a statement to the builder
func (b *StatementsBuilder) Add(statement Statement) *StatementsBuilder {
	if b.lastError != nil {
		return b
	}
	if statement == nil {
		b.lastError = errors.New("statement cannot be nil")
		return b
	}
	b.statements = append(b.statements, statement)
	return b
}

// AddQuery adds a query to the builder
func (b *StatementsBuilder) AddQuery(builder *CypherBuilder) *StatementsBuilder {
	if b.lastError != nil {
		return b
	}
	if builder == nil {
		b.lastError = errors.New("builder cannot be nil")
		return b
	}
	if builder.HasError() {
		b.lastError = builder.Error()
		return b
	}
	return b.Add(builder.statement)
}

// Error returns the last error that occurred during building
func (b *StatementsBuilder) Error() error {
	return b.lastError
}

// HasError checks if an error occurred during building
func (b *StatementsBuilder) HasError() bool {
	return b.lastError != nil
}

// Build finalizes the query and returns the statements and error
func (b *StatementsBuilder) Build() ([]Statement, error) {
	if b.lastError != nil {
		return nil, b.lastError
	}
	return b.statements, nil
}

// Statements returns the statements
// If an error occurred during building, it returns nil
func (b *StatementsBuilder) Statements() []Statement {
	if b.lastError != nil {
		return nil
	}
	return b.statements
}

// Cypher returns the Cypher string representations of the statements
// If an error occurred during building, it returns the error message
func (b *StatementsBuilder) Cypher() []string {
	if b.lastError != nil {
		return []string{fmt.Sprintf("ERROR: %s", b.lastError.Error())}
	}
	cyphers := make([]string, len(b.statements))
	for i, stmt := range b.statements {
		cyphers[i] = stmt.Cypher()
	}
	return cyphers
}

// Params returns the parameters of the statements
// If an error occurred during building, it returns nil
func (b *StatementsBuilder) Params() []map[string]interface{} {
	if b.lastError != nil {
		return nil
	}
	params := make([]map[string]interface{}, len(b.statements))
	for i, stmt := range b.statements {
		params[i] = stmt.Params()
	}
	return params
}

// Higher-level query construction functions that give better developer experience

// FindRelated is a shorthand for creating a query that finds nodes related to a given node
func FindRelated(source NodeExpression, relType string, target NodeExpression) *CypherBuilder {
	stmt, err := Match(source.RelationshipTo(target, relType)).
		Returning(target).
		Build()

	if err != nil {
		return &CypherBuilder{lastError: err}
	}
	return &CypherBuilder{statement: stmt}
}

// FindWithProperties is a shorthand for creating a query that finds nodes with the given properties
func FindWithProperties(node NodeExpression, properties map[string]interface{}) *CypherBuilder {
	propsExpr := make(map[string]Expression)
	for k, v := range properties {
		propsExpr[k] = autoExpression(v)
	}

	nodeWithProps := node.WithProperties(propsExpr).(NodeExpression)

	stmt, err := Match(nodeWithProps).
		Returning(nodeWithProps).
		Build()

	if err != nil {
		return &CypherBuilder{lastError: err}
	}
	return &CypherBuilder{statement: stmt}
}

// CreateNode is a shorthand for creating a node with the given properties
func CreateNode(node NodeExpression, properties map[string]interface{}) *CypherBuilder {
	propsExpr := make(map[string]Expression)
	for k, v := range properties {
		propsExpr[k] = autoExpression(v)
	}

	nodeWithProps := node.WithProperties(propsExpr).(NodeExpression)

	stmt, err := Create(nodeWithProps).
		Returning(nodeWithProps).
		Build()

	if err != nil {
		return &CypherBuilder{lastError: err}
	}
	return &CypherBuilder{statement: stmt}
}

// CreateRelationship is a shorthand for creating a relationship between two nodes
func CreateRelationship(source NodeExpression, relType string, target NodeExpression, properties map[string]interface{}) *CypherBuilder {
	propsExpr := make(map[string]Expression)
	for k, v := range properties {
		propsExpr[k] = autoExpression(v)
	}

	rel := source.RelationshipTo(target, relType)

	if properties != nil && len(properties) > 0 {
		// This cast would need to be implemented in the real code
		// rel = rel.WithProperties(propsExpr)
	}

	// Then manually construct the CREATE statement without a preceding MATCH
	createStmt, err := Create(rel).
		Returning(source, rel, target).
		Build()

	if err != nil {
		return &CypherBuilder{lastError: err}
	}
	return &CypherBuilder{statement: createStmt}
}

// Match adds an additional MATCH clause to the query
func (b *MatchBuilder) Match(pattern Expression) *MatchBuilder {
	// Create a new MatchBuilder with the pattern and previous builder
	return &MatchBuilder{
		prev: b,
		clause: &MatchClause{
			patterns: []Expression{pattern},
		},
	}
}

// Create adds a CREATE clause to the query
func (b *MatchBuilder) Create(pattern Expression) *CreateBuilder {
	// Create a new CreateBuilder with the pattern and previous builder
	return &CreateBuilder{
		prev: b,
		clause: &CreateClause{
			patterns: []Expression{pattern},
		},
	}
}
