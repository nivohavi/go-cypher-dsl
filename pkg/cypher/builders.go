package cypher

// Implementation of the Build method for MatchBuilder
func (m *MatchBuilder) Build() (Statement, error) {
	return &cypherStatement{
		builder: m,
	}, nil
}

// Implementation of the Build method for WithBuilder
func (w *WithBuilder) Build() (Statement, error) {
	return &cypherStatement{
		builder: w,
	}, nil
}

// Implementation of the Build method for UnwindBuilder
func (u *UnwindBuilder) Build() (Statement, error) {
	return &cypherStatement{
		builder: u,
	}, nil
}
