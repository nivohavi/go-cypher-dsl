package cypher

import (
	"testing"
)

func TestSimpleMatchReturn(t *testing.T) {
	// Define a node
	person := Node("Person").Named("p")

	// Build a MATCH ... RETURN query
	stmt, err := Match(person).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (p:`Person`) RETURN p"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}

	if len(stmt.Params()) != 0 {
		t.Errorf("Expected empty params map, got: %v", stmt.Params())
	}
}

func TestMatchWhereReturn(t *testing.T) {
	// Define a node
	person := Node("Person").Named("p").(*nodePattern)

	// Build a MATCH ... WHERE ... RETURN query
	stmt, err := Match(person).
		Where(person.Property("name").Eq("Tom Hanks")).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (p:`Person`) WHERE p.name = 'Tom Hanks' RETURN p"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestMatchWhereWithParameters(t *testing.T) {
	// Define a node
	person := Node("Person").Named("p").(*nodePattern)

	// Build a MATCH ... WHERE ... RETURN query with parameter
	stmt, err := Match(person).
		Where(person.Property("name").Eq(ParamWithValue("name", "Tom Hanks"))).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (p:`Person`) WHERE p.name = $name RETURN p"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}

	params := stmt.Params()
	if len(params) != 1 {
		t.Errorf("Expected 1 parameter, got: %v", params)
	}

	if params["name"] != "Tom Hanks" {
		t.Errorf("Expected parameter 'name' to be 'Tom Hanks', got: %v", params["name"])
	}
}

func TestCreateNodeWithProperties(t *testing.T) {
	// Define a node with properties
	person := Node("Person").(*nodePattern).WithProps(map[string]interface{}{
		"name": "Keanu Reeves",
		"born": 1964,
	})

	// Build a CREATE query
	stmt, err := Create(person).Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	// Note that map iteration order is non-deterministic, so we need to check for both possible outcomes
	expected1 := "CREATE (:`Person` {name: 'Keanu Reeves', born: 1964})"
	expected2 := "CREATE (:`Person` {born: 1964, name: 'Keanu Reeves'})"
	actual := stmt.Cypher()

	if actual != expected1 && actual != expected2 {
		t.Errorf("Expected: %s or %s, got: %s", expected1, expected2, actual)
	}
}

func TestReturnWithOrderBy(t *testing.T) {
	// Define a node
	movie := Node("Movie").Named("m").(*nodePattern)

	// Build a MATCH ... RETURN ... ORDER BY query
	stmt, err := Match(movie).
		Returning(movie.Property("title")).
		OrderBy(movie.Property("released")).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (m:`Movie`) RETURN m.title ORDER BY m.released"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestRelationshipPattern(t *testing.T) {
	// Define nodes and relationship
	person := Node("Person").Named("p")
	movie := Node("Movie").Named("m")
	acted := Relationship("ACTED_IN").Named("r").Build()

	// Build a pattern path (person)-[acted]->(movie)
	path := Pattern(person, acted, movie)

	// Build a MATCH query with the path
	stmt, err := Match(path).
		Returning(person, acted, movie).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (p:`Person`)-[r:`ACTED_IN`]->(m:`Movie`) RETURN p, r, m"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestLogicalExpressions(t *testing.T) {
	// Define a node
	movie := Node("Movie").Named("m").(*nodePattern)

	// Build a query with AND, OR, NOT
	year := movie.Property("released")
	title := movie.Property("title")

	condition := year.Gt(2000).And(
		year.Lt(2010).Or(
			title.Contains("Matrix").Not(),
		),
	)

	stmt, err := Match(movie).
		Where(condition).
		Returning(movie).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (m:`Movie`) WHERE (m.released > 2000 AND (m.released < 2010 OR NOT (m.title CONTAINS 'Matrix'))) RETURN m"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestWithClause(t *testing.T) {
	// Define nodes
	person := Node("Person").Named("p").(*nodePattern)
	movie := Node("Movie").Named("m").(*nodePattern)

	// Build a query with WITH
	stmt, err := Match(person).
		Where(person.Property("name").Eq("Tom Hanks")).
		With(person).
		Match(movie).
		Where(movie.Property("title").Contains(person.Property("name"))).
		Returning(movie).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (p:`Person`) WHERE p.name = 'Tom Hanks' WITH p MATCH (m:`Movie`) WHERE m.title CONTAINS p.name RETURN m"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestFunctions(t *testing.T) {
	// Define a node
	movie := Node("Movie").Named("m")
	movieNode := movie.(*nodePattern)

	// Build a query with functions
	stmt, err := Match(movie).
		Returning(
			Count(movie).(*functionExpression),
			Avg(movieNode.Property("rating")).(*functionExpression),
		).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (m:`Movie`) RETURN count(m), avg(m.rating)"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestStaticSchema(t *testing.T) {
	// Use the static schema
	person := Person()
	personNode := person.Named("p")
	movie := Movie()
	movieNode := movie.Named("m")
	actedIn := ActedIn()
	actedRel := actedIn.Direction(OUTGOING).Named("r").Build()

	// Build a query with the static schema
	stmt, err := Match(Pattern(personNode, actedRel, movieNode)).
		Where(person.Name().Eq(Literal("Keanu Reeves"))).
		Returning(movie.Title(), actedIn.Role()).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (p:`Person`)-[r:`ACTED_IN`]->(m:`Movie`) WHERE p.name = 'Keanu Reeves' RETURN m.title, r.role"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestUnwindQuery(t *testing.T) {
	// Define a node and a list
	names := List(Literal("Keanu Reeves"), Literal("Tom Hanks"), Literal("Meryl Streep"))

	// Build an UNWIND query
	stmt, err := Unwind(names, "name").
		Match(Node("Person").Named("p")).
		Where(Node("Person").Named("p").(*nodePattern).Property("name").Eq(Parameter("name"))).
		Returning(Parameter("name"), Count(Node("Person").Named("p")).(*functionExpression)).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "UNWIND ['Keanu Reeves', 'Tom Hanks', 'Meryl Streep'] AS name MATCH (p:`Person`) WHERE p.name = $name RETURN $name, count(p)"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestMergeQuery(t *testing.T) {
	// Define a node
	person := Node("Person").Named("p")

	// Build a MERGE query
	stmt, err := Merge(person).
		Returning(person).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MERGE (p:`Person`) RETURN p"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}

func TestSetQuery(t *testing.T) {
	// Define a node
	person := Node("Person").Named("p").(*nodePattern)

	// Build a MATCH ... SET query
	matchStmt, err := Match(person).
		Where(person.Property("name").Eq("Keanu Reeves")).
		Build()

	if err != nil {
		t.Fatalf("Error building match query: %v", err)
	}

	setStmt, err := Set().
		Property(person.Property("age"), Literal(56)).
		Build()

	if err != nil {
		t.Fatalf("Error building set query: %v", err)
	}

	matchExpected := "MATCH (p:`Person`) WHERE p.name = 'Keanu Reeves'"
	matchActual := matchStmt.Cypher()

	if matchActual != matchExpected {
		t.Errorf("Expected: %s, got: %s", matchExpected, matchActual)
	}

	setExpected := "SET p.age = 56"
	setActual := setStmt.Cypher()

	if setActual != setExpected {
		t.Errorf("Expected: %s, got: %s", setExpected, setActual)
	}
}

func TestDeleteQuery(t *testing.T) {
	// Define a node
	person := Node("Person").Named("p")

	// Build a MATCH ... DELETE query
	matchStmt, err := Match(person).
		Where(person.(*nodePattern).Property("name").Eq("John")).
		Build()

	if err != nil {
		t.Fatalf("Error building match query: %v", err)
	}

	deleteStmt, err := Delete(person).Build()

	if err != nil {
		t.Fatalf("Error building delete query: %v", err)
	}

	matchExpected := "MATCH (p:`Person`) WHERE p.name = 'John'"
	matchActual := matchStmt.Cypher()

	if matchActual != matchExpected {
		t.Errorf("Expected: %s, got: %s", matchExpected, matchActual)
	}

	deleteExpected := "DELETE p"
	deleteActual := deleteStmt.Cypher()

	if deleteActual != deleteExpected {
		t.Errorf("Expected: %s, got: %s", deleteExpected, deleteActual)
	}
}

func TestDetachDeleteQuery(t *testing.T) {
	// Define a node
	person := Node("Person").Named("p")

	// Build a MATCH ... DETACH DELETE query
	matchStmt, err := Match(person).
		Where(person.(*nodePattern).Property("name").Eq(Literal("John"))).
		Build()

	if err != nil {
		t.Fatalf("Error building match query: %v", err)
	}

	detachDeleteStmt, err := DetachDelete(person).Build()

	if err != nil {
		t.Fatalf("Error building detach delete query: %v", err)
	}

	matchExpected := "MATCH (p:`Person`) WHERE p.name = 'John'"
	matchActual := matchStmt.Cypher()

	if matchActual != matchExpected {
		t.Errorf("Expected: %s, got: %s", matchExpected, matchActual)
	}

	detachDeleteExpected := "DETACH DELETE p"
	detachDeleteActual := detachDeleteStmt.Cypher()

	if detachDeleteActual != detachDeleteExpected {
		t.Errorf("Expected: %s, got: %s", detachDeleteExpected, detachDeleteActual)
	}
}

func TestAliasedExpression(t *testing.T) {
	// Define a node
	movie := Node("Movie").Named("m").(*nodePattern)

	// Build a MATCH ... RETURN with aliased expressions
	stmt, err := Match(movie).
		Returning(
			As(movie.Property("title"), "movieTitle"),
			As(movie.Property("released"), "year"),
		).
		Build()

	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}

	expected := "MATCH (m:`Movie`) RETURN m.title AS movieTitle, m.released AS year"
	actual := stmt.Cypher()

	if actual != expected {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}
