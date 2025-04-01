package core

// StatementImpl implements the Statement interface
type StatementImpl struct {
	cypher     string
	params     map[string]any
	parameters *Parameters
}

// NewStatement creates a new statement with the given Cypher and params
func NewStatement(cypher string, params map[string]any) *StatementImpl {
	if params == nil {
		params = make(map[string]any)
	}
	return &StatementImpl{
		cypher: cypher,
		params: params,
	}
}

// NewStatementWithParams creates a new statement with the given Cypher and Parameters object
func NewStatementWithParams(cypher string, parameters *Parameters) *StatementImpl {
	var params map[string]any
	if parameters != nil {
		params = parameters.Get()
	} else {
		params = make(map[string]any)
	}
	return &StatementImpl{
		cypher:     cypher,
		params:     params,
		parameters: parameters,
	}
}

// Cypher returns the Cypher query string
func (s *StatementImpl) Cypher() string {
	return s.cypher
}

// Params returns the parameters for this statement
func (s *StatementImpl) Params() map[string]any {
	return s.params
}

// Parameters returns the Parameters object for this statement
func (s *StatementImpl) Parameters() *Parameters {
	return s.parameters
}

// Accept applies a visitor to this statement
func (s *StatementImpl) Accept(visitor StatementVisitor) any {
	return visitor.Visit(s)
}

// WithCypher creates a new statement with the given Cypher string
func (s *StatementImpl) WithCypher(cypher string) *StatementImpl {
	return &StatementImpl{
		cypher:     cypher,
		params:     s.params,
		parameters: s.parameters,
	}
}

// WithParams creates a new statement with the given parameters map
func (s *StatementImpl) WithParams(params map[string]any) *StatementImpl {
	return &StatementImpl{
		cypher:     s.cypher,
		params:     params,
		parameters: s.parameters,
	}
}

// Merge combines this statement with another one
func (s *StatementImpl) Merge(other *StatementImpl) *StatementImpl {
	if other == nil {
		return s
	}

	// Merge parameters
	mergedParams := make(map[string]any)
	for k, v := range s.params {
		mergedParams[k] = v
	}
	for k, v := range other.params {
		mergedParams[k] = v
	}

	// Merge Parameters objects if available
	var mergedParameters *Parameters
	if s.parameters != nil {
		mergedParameters = NewParameters()
		mergedParameters.Merge(s.parameters)
		if other.parameters != nil {
			mergedParameters.Merge(other.parameters)
		}
	} else if other.parameters != nil {
		mergedParameters = NewParameters()
		mergedParameters.Merge(other.parameters)
	}

	// Combine Cypher strings
	combinedCypher := s.cypher
	if other.cypher != "" {
		if combinedCypher != "" {
			combinedCypher += " "
		}
		combinedCypher += other.cypher
	}

	return &StatementImpl{
		cypher:     combinedCypher,
		params:     mergedParams,
		parameters: mergedParameters,
	}
}
