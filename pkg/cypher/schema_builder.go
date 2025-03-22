package cypher

// Schema defines a domain-specific entity schema
// This provides better compile-time safety and discovery of domain entities
type Schema interface {
	// Label returns the primary label for this schema
	Label() string

	// AllLabels returns all labels for this schema
	AllLabels() []string

	// Node creates a node for this schema, optionally with an alias
	Node(alias ...string) NodeExpression
}

// SchemaProperty represents a property in a schema
type SchemaProperty interface {
	// Name returns the property name
	Name() string

	// Expression returns the property expression for a given node or relationship
	Expression(container PropertyContainer) PropertyExpression
}

// BaseSchema implements common Schema functionality
type BaseSchema struct {
	name       string
	label      string
	labels     []string
	properties map[string]SchemaProperty
}

// NewSchema creates a new schema with the given primary label
func NewSchema(label string, additionalLabels ...string) *BaseSchema {
	labels := append([]string{label}, additionalLabels...)
	return &BaseSchema{
		name:       label,
		label:      label,
		labels:     labels,
		properties: make(map[string]SchemaProperty),
	}
}

// Label returns the primary label for this schema
func (s *BaseSchema) Label() string {
	return s.label
}

// AllLabels returns all labels for this schema
func (s *BaseSchema) AllLabels() []string {
	return s.labels
}

// Node creates a node for this schema, optionally with an alias
func (s *BaseSchema) Node(alias ...string) NodeExpression {
	node := Node(s.labels...)
	if len(alias) > 0 && alias[0] != "" {
		return node.Named(alias[0])
	}
	return node
}

// AddProperty adds a property to this schema
func (s *BaseSchema) AddProperty(name string) SchemaProperty {
	prop := &schemaPropertyImpl{name: name}
	s.properties[name] = prop
	return prop
}

// Property gets a property by name
func (s *BaseSchema) Property(name string) SchemaProperty {
	return s.properties[name]
}

// schemaPropertyImpl implements SchemaProperty
type schemaPropertyImpl struct {
	name string
}

// Name returns the property name
func (p *schemaPropertyImpl) Name() string {
	return p.name
}

// Expression returns the property expression for a given property container
func (p *schemaPropertyImpl) Expression(container PropertyContainer) PropertyExpression {
	return container.Property(p.name)
}

// NewPersonSchema creates a schema for a Person entity
func NewPersonSchema() *BaseSchema {
	schema := NewSchema("Person")
	schema.AddProperty("name")
	schema.AddProperty("age")
	schema.AddProperty("email")
	return schema
}

// NewMovieSchema creates a schema for a Movie entity
func NewMovieSchema() *BaseSchema {
	schema := NewSchema("Movie")
	schema.AddProperty("title")
	schema.AddProperty("released")
	schema.AddProperty("tagline")
	return schema
}

// NewActedInSchema creates a schema for ACTED_IN relationships
func NewActedInSchema() *BaseSchema {
	schema := NewSchema("ACTED_IN")
	schema.AddProperty("role")
	return schema
}

// CreateTypedSchema creates helper functions to access schema properties
func CreateTypedSchema(schema *BaseSchema) map[string]func(PropertyContainer) PropertyExpression {
	propertyAccessors := make(map[string]func(PropertyContainer) PropertyExpression)

	for name, prop := range schema.properties {
		property := prop // Capture the property in a new variable
		accessor := func(container PropertyContainer) PropertyExpression {
			return property.Expression(container)
		}
		propertyAccessors[name] = accessor
	}

	return propertyAccessors
}

// Example usage:
// personSchema := NewPersonSchema()
// person := personSchema.Node("p")
// nameProperty := personSchema.Property("name")
// query := Match(person).Where(nameProperty.Expression(person).Eq("Tom Hanks")).Return(person)
