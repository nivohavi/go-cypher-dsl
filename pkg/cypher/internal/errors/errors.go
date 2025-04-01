package errors

import (
	"fmt"
)

// ErrType represents different categories of errors that can occur during query building
type ErrType int

const (
	// ErrTypeSyntax indicates a syntax error in the query
	ErrTypeSyntax ErrType = iota

	// ErrTypeValidation indicates a validation error in the query
	ErrTypeValidation

	// ErrTypeParameter indicates an error with parameters
	ErrTypeParameter

	// ErrTypeUnsupported indicates an unsupported operation
	ErrTypeUnsupported

	// ErrTypeIncomplete indicates an incomplete query state
	ErrTypeIncomplete
)

// Error represents a cypher DSL error with context
type Error struct {
	Type    ErrType
	Message string
	Context string
	Cause   error
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (context: %s, caused by: %v)", e.Type.String(), e.Message, e.Context, e.Cause)
	}
	return fmt.Sprintf("%s: %s (context: %s)", e.Type.String(), e.Message, e.Context)
}

// Unwrap returns the underlying cause error
func (e *Error) Unwrap() error {
	return e.Cause
}

// String returns a string representation of the error type
func (t ErrType) String() string {
	switch t {
	case ErrTypeSyntax:
		return "SYNTAX_ERROR"
	case ErrTypeValidation:
		return "VALIDATION_ERROR"
	case ErrTypeParameter:
		return "PARAMETER_ERROR"
	case ErrTypeUnsupported:
		return "UNSUPPORTED_ERROR"
	case ErrTypeIncomplete:
		return "INCOMPLETE_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

// NewSyntaxError creates a new syntax error
func NewSyntaxError(message, context string, cause error) error {
	return &Error{
		Type:    ErrTypeSyntax,
		Message: message,
		Context: context,
		Cause:   cause,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message, context string, cause error) error {
	return &Error{
		Type:    ErrTypeValidation,
		Message: message,
		Context: context,
		Cause:   cause,
	}
}

// NewParameterError creates a new parameter error
func NewParameterError(message, context string, cause error) error {
	return &Error{
		Type:    ErrTypeParameter,
		Message: message,
		Context: context,
		Cause:   cause,
	}
}

// NewUnsupportedError creates a new unsupported operation error
func NewUnsupportedError(message, context string, cause error) error {
	return &Error{
		Type:    ErrTypeUnsupported,
		Message: message,
		Context: context,
		Cause:   cause,
	}
}

// NewIncompleteError creates a new incomplete query error
func NewIncompleteError(message, context string, cause error) error {
	return &Error{
		Type:    ErrTypeIncomplete,
		Message: message,
		Context: context,
		Cause:   cause,
	}
}

// IsSyntaxError checks if the error is a syntax error
func IsSyntaxError(err error) bool {
	e, ok := err.(*Error)
	return ok && e.Type == ErrTypeSyntax
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	e, ok := err.(*Error)
	return ok && e.Type == ErrTypeValidation
}

// IsParameterError checks if the error is a parameter error
func IsParameterError(err error) bool {
	e, ok := err.(*Error)
	return ok && e.Type == ErrTypeParameter
}

// IsUnsupportedError checks if the error is an unsupported operation error
func IsUnsupportedError(err error) bool {
	e, ok := err.(*Error)
	return ok && e.Type == ErrTypeUnsupported
}

// IsIncompleteError checks if the error is an incomplete query error
func IsIncompleteError(err error) bool {
	e, ok := err.(*Error)
	return ok && e.Type == ErrTypeIncomplete
}
