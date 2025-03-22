package cypher

// QueryBuilder is a utility wrapper for a statement
type QueryBuilder struct {
	statement Statement
	lastError error
}

// Error returns the error that occurred during building
func (b *QueryBuilder) Error() error {
	return b.lastError
}

// HasError checks if an error occurred during building
func (b *QueryBuilder) HasError() bool {
	return b.lastError != nil
}

// Statement returns the built statement
func (b *QueryBuilder) Statement() Statement {
	return b.statement
}

// Cypher returns the query string
func (b *QueryBuilder) Cypher() string {
	if b.statement == nil {
		return ""
	}
	return b.statement.Cypher()
}

// Params returns the query parameters
func (b *QueryBuilder) Params() map[string]interface{} {
	if b.statement == nil {
		return make(map[string]interface{})
	}
	return b.statement.Params()
}

// TransactionBuilder represents a sequence of Cypher statements
type TransactionBuilder struct {
	statements []Statement
	err        error
}

// NewTransactionBuilder creates a new transaction
func NewTransactionBuilder() *TransactionBuilder {
	return &TransactionBuilder{
		statements: []Statement{},
	}
}

// Add adds a statement to the transaction
func (t *TransactionBuilder) Add(stmt Statement) {
	if t.err != nil {
		return
	}

	t.statements = append(t.statements, stmt)
}

// Statements returns all statements in the transaction
func (t *TransactionBuilder) Statements() []Statement {
	return t.statements
}

// Error returns any error that occurred
func (t *TransactionBuilder) Error() error {
	return t.err
}

// HasError checks if an error occurred
func (t *TransactionBuilder) HasError() bool {
	return t.err != nil
}
