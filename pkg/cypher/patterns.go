package cypher

// FindNodesByProperties creates a query to find nodes by their properties
func FindNodesByProperties(label string, properties map[string]interface{}) *CypherBuilder {
	// Create a node with the given label
	node := Node(label).Named("n")

	// Prepare the WHERE clause conditions
	var conditions []BooleanExpression
	for key, value := range properties {
		propExpr := node.(*nodePattern).Property(key).Eq(value)
		conditions = append(conditions, propExpr)
	}

	// Combine conditions with AND
	var whereCondition BooleanExpression
	if len(conditions) > 0 {
		whereCondition = conditions[0]
		for i := 1; i < len(conditions); i++ {
			whereCondition = whereCondition.And(conditions[i])
		}
	}

	// Build the query
	matchBuilder := Match(node)
	if whereCondition != nil {
		matchBuilder = matchBuilder.Where(whereCondition)
	}

	stmt, err := matchBuilder.Returning(node).Build()
	if err != nil {
		return &CypherBuilder{lastError: err}
	}

	return &CypherBuilder{statement: stmt}
}

// CreateNodeWithProperties creates a query to create a node with properties
func CreateNodeWithProperties(label string, properties map[string]interface{}) *CypherBuilder {
	// Create a node with the given label
	node := Node(label).Named("n")

	// Convert properties to expressions
	propsExpr := make(map[string]Expression)
	for k, v := range properties {
		propsExpr[k] = autoExpression(v)
	}

	// Apply properties to the node
	nodeWithProps := node.(*nodePattern).WithProperties(propsExpr)

	// Build the query
	stmt, err := Create(nodeWithProps).Returning(nodeWithProps).Build()
	if err != nil {
		return &CypherBuilder{lastError: err}
	}

	return &CypherBuilder{statement: stmt}
}

// RelateNodes creates a query to create a relationship between two nodes
func RelateNodes(fromNode NodeExpression, relType string, toNode NodeExpression, properties map[string]interface{}) *CypherBuilder {
	// Create the relationship pattern
	rel := fromNode.RelationshipTo(toNode, relType)

	// Convert properties to expressions
	if properties != nil && len(properties) > 0 {
		propsExpr := make(map[string]Expression)
		for k, v := range properties {
			propsExpr[k] = autoExpression(v)
		}

		// Apply properties to the relationship
		rel = rel.(RelationshipPattern).WithProperties(propsExpr)
	}

	// Create a pattern with the relationship
	pattern := Pattern(fromNode, rel, toNode)

	// Build the query
	stmt, err := Create(pattern).Returning(fromNode, rel, toNode).Build()
	if err != nil {
		return &CypherBuilder{lastError: err}
	}

	return &CypherBuilder{statement: stmt}
}

// FindRelatedNodes creates a query to find nodes related to a source node
func FindRelatedNodes(sourceNode NodeExpression, relType string, targetLabel string) *CypherBuilder {
	// Create the target node
	targetNode := Node(targetLabel).Named("target")

	// Create the relationship pattern
	rel := sourceNode.RelationshipTo(targetNode, relType)

	// Create a pattern with the relationship
	pattern := Pattern(sourceNode, rel, targetNode)

	// Build the query
	stmt, err := Match(pattern).Returning(targetNode).Build()
	if err != nil {
		return &CypherBuilder{lastError: err}
	}

	return &CypherBuilder{statement: stmt}
}

// DeleteNode creates a query to delete a node (with optional detach)
func DeleteNode(node NodeExpression, detach bool) *CypherBuilder {
	// Build the match part
	matchBuilder := Match(node)

	// Build the delete part
	var deleteBuilder *DeleteBuilder
	if detach {
		deleteBuilder = DetachDelete(node)
	} else {
		deleteBuilder = Delete(node)
	}

	// Combine into a transaction
	tx := NewTransaction()

	// Add the match and delete statements
	matchStmt, err := matchBuilder.Build()
	if err != nil {
		return &CypherBuilder{lastError: err}
	}

	deleteStmt, err := deleteBuilder.Build()
	if err != nil {
		return &CypherBuilder{lastError: err}
	}

	tx.Add(matchStmt)
	tx.Add(deleteStmt)

	// Return a combined builder
	builder := &CypherBuilder{}

	// Use HasError/Error method instead of directly accessing lastError
	if tx.HasError() {
		builder.lastError = tx.Error()
	}
	return builder
}
