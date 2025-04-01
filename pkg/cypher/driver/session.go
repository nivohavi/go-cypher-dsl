// Package driver provides helpers for working with the Neo4j driver
package driver

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// SessionManager simplifies working with Neo4j sessions
type SessionManager struct {
	driver neo4j.Driver
}

// NewSessionManager creates a new SessionManager
func NewSessionManager(driver neo4j.Driver) *SessionManager {
	return &SessionManager{
		driver: driver,
	}
}

// ExecuteRead executes a read query using the provided statement
func (sm *SessionManager) ExecuteRead(ctx context.Context, statement core.Statement,
	handler func(neo4j.Result) (any, error)) (any, error) {

	session := sm.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	return session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(statement.Cypher(), statement.Params())
		if err != nil {
			return nil, err
		}
		return handler(result)
	})
}

// ExecuteWrite executes a write query using the provided statement
func (sm *SessionManager) ExecuteWrite(ctx context.Context, statement core.Statement,
	handler func(neo4j.Result) (any, error)) (any, error) {

	session := sm.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	return session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(statement.Cypher(), statement.Params())
		if err != nil {
			return nil, err
		}
		return handler(result)
	})
}

// ExecuteBatchWrite executes multiple write statements in a single transaction
func (sm *SessionManager) ExecuteBatchWrite(ctx context.Context, statements []core.Statement,
	handler func([]neo4j.Result) (any, error)) (any, error) {

	session := sm.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	return session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		var results []neo4j.Result

		for _, stmt := range statements {
			result, err := tx.Run(stmt.Cypher(), stmt.Params())
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}

		return handler(results)
	})
}

// QueryHelper provides common handler functions for Neo4j results
type QueryHelper struct{}

// CollectSingle returns a handler function that collects a single record from the result
func (qh *QueryHelper) CollectSingle(field string) func(neo4j.Result) (any, error) {
	return func(result neo4j.Result) (any, error) {
		if result.Next() {
			record := result.Record()
			value, found := record.Get(field)
			if !found {
				return nil, nil
			}
			return value, nil
		}
		return nil, nil
	}
}

// CollectList returns a handler function that collects all records into a list
func (qh *QueryHelper) CollectList(field string) func(neo4j.Result) (any, error) {
	return func(result neo4j.Result) (any, error) {
		var list []any
		for result.Next() {
			record := result.Record()
			value, found := record.Get(field)
			if !found {
				continue
			}
			list = append(list, value)
		}
		return list, nil
	}
}

// CollectMap returns a handler function that collects all records with key-value pairs
func (qh *QueryHelper) CollectMap(keyField, valueField string) func(neo4j.Result) (any, error) {
	return func(result neo4j.Result) (any, error) {
		resultMap := make(map[string]any)
		for result.Next() {
			record := result.Record()
			key, found := record.Get(keyField)
			if !found {
				continue
			}
			value, found := record.Get(valueField)
			if !found {
				continue
			}
			keyStr, ok := key.(string)
			if !ok {
				continue
			}
			resultMap[keyStr] = value
		}
		return resultMap, nil
	}
}

// CollectAll returns a handler function that collects all records as-is
func (qh *QueryHelper) CollectAll() func(neo4j.Result) (any, error) {
	return func(result neo4j.Result) (any, error) {
		var records []map[string]any
		for result.Next() {
			record := result.Record()
			recordMap := make(map[string]any)
			for _, key := range record.Keys {
				value, found := record.Get(key)
				if !found {
					continue
				}
				recordMap[key] = value
			}
			records = append(records, recordMap)
		}
		return records, nil
	}
}

// CountResults returns a handler function that counts the number of records
func (qh *QueryHelper) CountResults() func(neo4j.Result) (any, error) {
	return func(result neo4j.Result) (any, error) {
		count := 0
		for result.Next() {
			count++
		}
		return count, nil
	}
}

// NewQueryHelper creates a new QueryHelper
func NewQueryHelper() *QueryHelper {
	return &QueryHelper{}
}
