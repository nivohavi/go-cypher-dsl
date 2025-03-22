package cypher

// TypedSchema represents a defined schema (node, relationship type)
type TypedSchema struct {
	label      string
	properties map[string]*TypedProperty
}

// TypedProperty represents a defined property in a schema
type TypedProperty struct {
	name   string
	schema *TypedSchema
}

// TypedNode represents a node with a specific schema
type TypedNode struct {
	schema *TypedSchema
	alias  string
}

// NewTypedSchema creates a new schema definition
func NewTypedSchema(label string) *TypedSchema {
	return &TypedSchema{
		label:      label,
		properties: make(map[string]*TypedProperty),
	}
}

// AddProperty adds a property to the schema
func (s *TypedSchema) AddProperty(name string) *TypedProperty {
	prop := &TypedProperty{
		name:   name,
		schema: s,
	}
	s.properties[name] = prop
	return prop
}

// Node creates a node with this schema
func (s *TypedSchema) Node(alias string) PatternElement {
	nodePattern := Node(s.label).Named(alias)

	// Create a typed node but we don't need to return it
	// it's just for reference tracking if needed later
	_ = &TypedNode{
		schema: s,
		alias:  alias,
	}

	return nodePattern
}

// Property returns a property from this schema by name
func (s *TypedSchema) Property(name string) *TypedProperty {
	return s.properties[name]
}

// Of applies this property to a specific node instance
func (p *TypedProperty) Of(node PatternElement) PropertyExpression {
	// Convert to nodePattern to get the alias
	var alias string

	if namedExpr, ok := node.(NamedExpression); ok {
		alias = namedExpr.SymbolicName()
	}

	if alias == "" {
		// Fallback to a generic alias if none specified
		alias = "n"
	}

	// Return a schema property expression with the node alias
	return &schemaPropertyExpression{
		nodeType:     &NodeType{Label: p.schema.label},
		propertyName: p.name,
		nodeAlias:    alias,
	}
}
