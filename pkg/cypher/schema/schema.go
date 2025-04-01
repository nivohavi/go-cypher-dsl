// Package schema provides helper functions for working with Neo4j schema operations
package schema

import (
	"fmt"
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/core"
)

// CreateNodeKeyConstraint generates a Cypher statement to create a node key constraint
func CreateNodeKeyConstraint(constraintName string, label string, properties ...string) (core.Statement, error) {
	if len(properties) == 0 {
		return nil, fmt.Errorf("at least one property is required for a node key constraint")
	}

	var propsList strings.Builder
	for i, prop := range properties {
		if i > 0 {
			propsList.WriteString(", ")
		}
		propsList.WriteString(prop)
	}

	query := fmt.Sprintf("CREATE CONSTRAINT %s IF NOT EXISTS FOR (n:%s) REQUIRE (n.%s) IS NODE KEY",
		constraintName, label, propsList.String())

	return core.NewStatement(query, nil), nil
}

// CreateUniqueConstraint generates a Cypher statement to create a uniqueness constraint
func CreateUniqueConstraint(constraintName string, label string, property string) (core.Statement, error) {
	query := fmt.Sprintf("CREATE CONSTRAINT %s IF NOT EXISTS FOR (n:%s) REQUIRE n.%s IS UNIQUE",
		constraintName, label, property)

	return core.NewStatement(query, nil), nil
}

// CreateExistsConstraint generates a Cypher statement to create a property existence constraint
func CreateExistsConstraint(constraintName string, label string, property string) (core.Statement, error) {
	query := fmt.Sprintf("CREATE CONSTRAINT %s IF NOT EXISTS FOR (n:%s) REQUIRE n.%s IS NOT NULL",
		constraintName, label, property)

	return core.NewStatement(query, nil), nil
}

// CreateRelationshipConstraint generates a Cypher statement to create a relationship constraint
func CreateRelationshipConstraint(constraintName string, relType string, property string) (core.Statement, error) {
	query := fmt.Sprintf("CREATE CONSTRAINT %s IF NOT EXISTS FOR ()-[r:%s]-() REQUIRE r.%s IS NOT NULL",
		constraintName, relType, property)

	return core.NewStatement(query, nil), nil
}

// CreateIndex generates a Cypher statement to create an index
func CreateIndex(indexName string, label string, properties ...string) (core.Statement, error) {
	if len(properties) == 0 {
		return nil, fmt.Errorf("at least one property is required for an index")
	}

	var propsList strings.Builder
	for i, prop := range properties {
		if i > 0 {
			propsList.WriteString(", ")
		}
		propsList.WriteString("n." + prop)
	}

	query := fmt.Sprintf("CREATE INDEX %s IF NOT EXISTS FOR (n:%s) ON (%s)",
		indexName, label, propsList.String())

	return core.NewStatement(query, nil), nil
}

// CreateFullTextIndex generates a Cypher statement to create a full-text search index
func CreateFullTextIndex(indexName string, labels []string, properties []string) (core.Statement, error) {
	if len(labels) == 0 {
		return nil, fmt.Errorf("at least one label is required for a full-text index")
	}

	if len(properties) == 0 {
		return nil, fmt.Errorf("at least one property is required for a full-text index")
	}

	var labelsList strings.Builder
	for i, label := range labels {
		if i > 0 {
			labelsList.WriteString(", ")
		}
		labelsList.WriteString(fmt.Sprintf("'%s'", label))
	}

	var propsList strings.Builder
	for i, prop := range properties {
		if i > 0 {
			propsList.WriteString(", ")
		}
		propsList.WriteString(fmt.Sprintf("'%s'", prop))
	}

	query := fmt.Sprintf("CALL db.index.fulltext.createNodeIndex('%s', [%s], [%s])",
		indexName, labelsList.String(), propsList.String())

	return core.NewStatement(query, nil), nil
}

// DropConstraint generates a Cypher statement to drop a constraint
func DropConstraint(constraintName string) (core.Statement, error) {
	query := fmt.Sprintf("DROP CONSTRAINT %s IF EXISTS", constraintName)

	return core.NewStatement(query, nil), nil
}

// DropIndex generates a Cypher statement to drop an index
func DropIndex(indexName string) (core.Statement, error) {
	query := fmt.Sprintf("DROP INDEX %s IF EXISTS", indexName)

	return core.NewStatement(query, nil), nil
}

// ShowConstraints generates a Cypher statement to show all constraints
func ShowConstraints() (core.Statement, error) {
	return core.NewStatement("SHOW CONSTRAINTS", nil), nil
}

// ShowIndexes generates a Cypher statement to show all indexes
func ShowIndexes() (core.Statement, error) {
	return core.NewStatement("SHOW INDEXES", nil), nil
}
