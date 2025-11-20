package schema

import (
	"strings"
	"testing"
)

func TestCreateUniqueConstraint(t *testing.T) {
	stmt, err := CreateUniqueConstraint("user_email_unique", "User", "email")
	if err != nil {
		t.Fatalf("CreateUniqueConstraint() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "CREATE CONSTRAINT") {
		t.Errorf("Cypher() = %q, should contain 'CREATE CONSTRAINT'", cypher)
	}
	if !strings.Contains(cypher, "IS UNIQUE") {
		t.Errorf("Cypher() = %q, should contain 'IS UNIQUE'", cypher)
	}
}

func TestCreateNodeKeyConstraint(t *testing.T) {
	stmt, err := CreateNodeKeyConstraint("user_id_key", "User", "id")
	if err != nil {
		t.Fatalf("CreateNodeKeyConstraint() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "IS NODE KEY") {
		t.Errorf("Cypher() = %q, should contain 'IS NODE KEY'", cypher)
	}
}

func TestCreateNodeKeyConstraintMultipleProperties(t *testing.T) {
	stmt, err := CreateNodeKeyConstraint("user_composite_key", "User", "id", "email")
	if err != nil {
		t.Fatalf("CreateNodeKeyConstraint() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "id") || !strings.Contains(cypher, "email") {
		t.Errorf("Cypher() = %q, should contain both properties", cypher)
	}
}

func TestCreateExistsConstraint(t *testing.T) {
	stmt, err := CreateExistsConstraint("user_name_exists", "User", "name")
	if err != nil {
		t.Fatalf("CreateExistsConstraint() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "IS NOT NULL") {
		t.Errorf("Cypher() = %q, should contain 'IS NOT NULL'", cypher)
	}
}

func TestCreateRelationshipConstraint(t *testing.T) {
	stmt, err := CreateRelationshipConstraint("rel_amount_exists", "TRANSFERRED", "amount")
	if err != nil {
		t.Fatalf("CreateRelationshipConstraint() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "TRANSFERRED") {
		t.Errorf("Cypher() = %q, should contain 'TRANSFERRED'", cypher)
	}
}

func TestCreateIndex(t *testing.T) {
	stmt, err := CreateIndex("user_name_idx", "User", "firstName", "lastName")
	if err != nil {
		t.Fatalf("CreateIndex() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "CREATE INDEX") {
		t.Errorf("Cypher() = %q, should contain 'CREATE INDEX'", cypher)
	}
}

func TestCreateFullTextIndex(t *testing.T) {
	stmt, err := CreateFullTextIndex("content_search", []string{"Post", "Comment"}, []string{"title", "content"})
	if err != nil {
		t.Fatalf("CreateFullTextIndex() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "db.index.fulltext.createNodeIndex") {
		t.Errorf("Cypher() = %q, should contain 'db.index.fulltext.createNodeIndex'", cypher)
	}
}

func TestDropConstraint(t *testing.T) {
	stmt, err := DropConstraint("user_email_unique")
	if err != nil {
		t.Fatalf("DropConstraint() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "DROP CONSTRAINT") {
		t.Errorf("Cypher() = %q, should contain 'DROP CONSTRAINT'", cypher)
	}
}

func TestDropIndex(t *testing.T) {
	stmt, err := DropIndex("user_name_idx")
	if err != nil {
		t.Fatalf("DropIndex() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "DROP INDEX") {
		t.Errorf("Cypher() = %q, should contain 'DROP INDEX'", cypher)
	}
}

func TestShowConstraints(t *testing.T) {
	stmt, err := ShowConstraints()
	if err != nil {
		t.Fatalf("ShowConstraints() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SHOW CONSTRAINTS") {
		t.Errorf("Cypher() = %q, should contain 'SHOW CONSTRAINTS'", cypher)
	}
}

func TestShowIndexes(t *testing.T) {
	stmt, err := ShowIndexes()
	if err != nil {
		t.Fatalf("ShowIndexes() error = %v", err)
	}

	cypher := stmt.Cypher()
	if !strings.Contains(cypher, "SHOW INDEXES") {
		t.Errorf("Cypher() = %q, should contain 'SHOW INDEXES'", cypher)
	}
}

func TestCreateNodeKeyConstraintNoProperties(t *testing.T) {
	_, err := CreateNodeKeyConstraint("test", "User")
	if err == nil {
		t.Error("CreateNodeKeyConstraint() with no properties should return error")
	}
}

func TestCreateIndexNoProperties(t *testing.T) {
	_, err := CreateIndex("test", "User")
	if err == nil {
		t.Error("CreateIndex() with no properties should return error")
	}
}

func TestCreateFullTextIndexNoLabels(t *testing.T) {
	_, err := CreateFullTextIndex("test", []string{}, []string{"prop"})
	if err == nil {
		t.Error("CreateFullTextIndex() with no labels should return error")
	}
}

func TestCreateFullTextIndexNoProperties(t *testing.T) {
	_, err := CreateFullTextIndex("test", []string{"User"}, []string{})
	if err == nil {
		t.Error("CreateFullTextIndex() with no properties should return error")
	}
}

