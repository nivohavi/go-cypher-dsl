package core

import (
	"errors"
	"fmt"
)

// Error types for Cypher operations
var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrInvalidParameter  = errors.New("invalid parameter")
	ErrInvalidPattern    = errors.New("invalid pattern")
	ErrInvalidQuery      = errors.New("invalid query")
	ErrInvalidProperty   = errors.New("invalid property")
	ErrMissingAlias      = errors.New("missing alias")
	ErrNodeRequired      = errors.New("node required")
	ErrBuildFailed       = errors.New("build failed")
)

// CypherError represents an error that occurred during Cypher query construction
type CypherError struct {
	Err       error
	Message   string
	Component string
	Query     string
}

// Error returns the error message
func (e *CypherError) Error() string {
	if e.Message != "" {
		if e.Component != "" {
			return fmt.Sprintf("%s: %s", e.Component, e.Message)
		}
		return e.Message
	}
	if e.Err != nil {
		if e.Component != "" {
			return fmt.Sprintf("%s: %s", e.Component, e.Err.Error())
		}
		return e.Err.Error()
	}
	return "unknown error"
}

// Unwrap returns the underlying error
func (e *CypherError) Unwrap() error {
	return e.Err
}

// NewError creates a new CypherError
func NewError(err error, message string) *CypherError {
	return &CypherError{
		Err:     err,
		Message: message,
	}
}

// NewComponentError creates a new CypherError with a component
func NewComponentError(err error, component, message string) *CypherError {
	return &CypherError{
		Err:       err,
		Message:   message,
		Component: component,
	}
}

// WithQuery adds query information to the error
func (e *CypherError) WithQuery(query string) *CypherError {
	e.Query = query
	return e
}
