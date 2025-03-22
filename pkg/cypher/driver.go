package cypher

import (
	"context"
	"errors"
)

// ExecuteFunc is a function that can execute a Cypher query
type ExecuteFunc func(cypher string, params map[string]any) (any, error)

// Execute runs the statement using the provided execution function
func Execute(stmt Statement, execute ExecuteFunc) (any, error) {
	if stmt == nil {
		return nil, errors.New("statement cannot be nil")
	}

	cypher := stmt.Cypher()
	params := stmt.Params()

	return execute(cypher, params)
}

// Neo4jDriver interface abstracts the Neo4j driver operations to make testing easier
// and avoid direct dependency on the Neo4j driver package
type Neo4jDriver interface {
	// ExecuteQuery executes a query and returns a result
	ExecuteQuery(ctx context.Context, cypher string, params map[string]any, options ...any) (any, error)
}

// ExecuteWithDriver runs the statement using the provided Neo4j driver
func ExecuteWithDriver(ctx context.Context, driver Neo4jDriver, stmt Statement, options ...any) (any, error) {
	if stmt == nil {
		return nil, errors.New("statement cannot be nil")
	}
	if driver == nil {
		return nil, errors.New("driver cannot be nil")
	}

	cypher := stmt.Cypher()
	params := stmt.Params()

	return driver.ExecuteQuery(ctx, cypher, params, options...)
}

// Neo4jSession interface abstracts a Neo4j session
type Neo4jSession interface {
	// Run executes a query in a session and returns a result
	Run(cypher string, params map[string]any) (any, error)
}

// ExecuteWithSession runs the statement using the provided Neo4j session
func ExecuteWithSession(session Neo4jSession, stmt Statement) (any, error) {
	if stmt == nil {
		return nil, errors.New("statement cannot be nil")
	}
	if session == nil {
		return nil, errors.New("session cannot be nil")
	}

	cypher := stmt.Cypher()
	params := stmt.Params()

	return session.Run(cypher, params)
}

// Neo4jTransaction interface abstracts a Neo4j transaction
type Neo4jTransaction interface {
	// Run executes a query in a transaction and returns a result
	Run(cypher string, params map[string]any) (any, error)
}

// ExecuteWithTransaction runs the statement using the provided Neo4j transaction
func ExecuteWithTransaction(tx Neo4jTransaction, stmt Statement) (any, error) {
	if stmt == nil {
		return nil, errors.New("statement cannot be nil")
	}
	if tx == nil {
		return nil, errors.New("transaction cannot be nil")
	}

	cypher := stmt.Cypher()
	params := stmt.Params()

	return tx.Run(cypher, params)
}
