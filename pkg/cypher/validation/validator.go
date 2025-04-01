package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/internal/errors"
)

// ValidationLevel represents the strictness of validation
type ValidationLevel int

const (
	// ValidationLevelOff disables validation
	ValidationLevelOff ValidationLevel = iota

	// ValidationLevelBasic performs basic syntax checks
	ValidationLevelBasic

	// ValidationLevelStrict performs thorough validation
	ValidationLevelStrict
)

// ValidationRule represents a rule that can be checked against a Cypher query
type ValidationRule struct {
	ID          string
	Description string
	Level       ValidationLevel
	Check       func(query string) error
}

// DefaultRules returns the default set of validation rules
func DefaultRules() []ValidationRule {
	return []ValidationRule{
		{
			ID:          "UNMATCHED_PARENTHESES",
			Description: "Unmatched parentheses in query",
			Level:       ValidationLevelBasic,
			Check: func(query string) error {
				count := 0
				for _, c := range query {
					if c == '(' {
						count++
					} else if c == ')' {
						count--
						if count < 0 {
							return errors.NewSyntaxError(
								"unmatched closing parenthesis",
								"parentheses validation",
								nil,
							)
						}
					}
				}
				if count > 0 {
					return errors.NewSyntaxError(
						fmt.Sprintf("missing %d closing parentheses", count),
						"parentheses validation",
						nil,
					)
				}
				return nil
			},
		},
		{
			ID:          "UNMATCHED_BRACKETS",
			Description: "Unmatched square brackets in query",
			Level:       ValidationLevelBasic,
			Check: func(query string) error {
				count := 0
				for _, c := range query {
					if c == '[' {
						count++
					} else if c == ']' {
						count--
						if count < 0 {
							return errors.NewSyntaxError(
								"unmatched closing bracket",
								"brackets validation",
								nil,
							)
						}
					}
				}
				if count > 0 {
					return errors.NewSyntaxError(
						fmt.Sprintf("missing %d closing brackets", count),
						"brackets validation",
						nil,
					)
				}
				return nil
			},
		},
		{
			ID:          "UNMATCHED_BRACES",
			Description: "Unmatched curly braces in query",
			Level:       ValidationLevelBasic,
			Check: func(query string) error {
				count := 0
				for _, c := range query {
					if c == '{' {
						count++
					} else if c == '}' {
						count--
						if count < 0 {
							return errors.NewSyntaxError(
								"unmatched closing brace",
								"braces validation",
								nil,
							)
						}
					}
				}
				if count > 0 {
					return errors.NewSyntaxError(
						fmt.Sprintf("missing %d closing braces", count),
						"braces validation",
						nil,
					)
				}
				return nil
			},
		},
		{
			ID:          "UNQUOTED_PROPERTY",
			Description: "Property access looks incorrect (missing backticks or quotes)",
			Level:       ValidationLevelStrict,
			Check: func(query string) error {
				// Check for property access without proper quotes for labels with spaces
				// e.g. n.My Property instead of n.`My Property`
				re := regexp.MustCompile(`[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+\s+[a-zA-Z0-9_]+`)
				if matches := re.FindAllString(query, -1); len(matches) > 0 {
					return errors.NewSyntaxError(
						fmt.Sprintf("potential unquoted property with spaces: %s", matches[0]),
						"property validation",
						nil,
					)
				}
				return nil
			},
		},
		{
			ID:          "MISSING_RELATIONSHIP_DIRECTION",
			Description: "Missing or ambiguous relationship direction",
			Level:       ValidationLevelStrict,
			Check: func(query string) error {
				// Check for relationships without direction (MATCH (a)-[r]-(b))
				// This is not necessarily an error but worth flagging in strict mode
				if strings.Contains(query, "MATCH") &&
					regexp.MustCompile(`\)-\[[^\]]+\]-\(`).MatchString(query) {
					return errors.NewValidationError(
						"relationship without direction found, consider using -> or <-",
						"relationship validation",
						nil,
					)
				}
				return nil
			},
		},
	}
}

// Validator validates Cypher queries
type Validator struct {
	rules []ValidationRule
	level ValidationLevel
}

// NewValidator creates a new validator with the given rules and level
func NewValidator(rules []ValidationRule, level ValidationLevel) *Validator {
	return &Validator{
		rules: rules,
		level: level,
	}
}

// DefaultValidator creates a validator with default rules and basic validation
func DefaultValidator() *Validator {
	return NewValidator(DefaultRules(), ValidationLevelBasic)
}

// StrictValidator creates a validator with default rules and strict validation
func StrictValidator() *Validator {
	return NewValidator(DefaultRules(), ValidationLevelStrict)
}

// SetLevel sets the validation level
func (v *Validator) SetLevel(level ValidationLevel) {
	v.level = level
}

// AddRule adds a custom validation rule
func (v *Validator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// ValidateQuery validates a Cypher query against the rules
func (v *Validator) ValidateQuery(query string) []error {
	if v.level == ValidationLevelOff {
		return nil
	}

	var errs []error

	for _, rule := range v.rules {
		if rule.Level <= v.level {
			if err := rule.Check(query); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

// IsValid checks if a query is valid
func (v *Validator) IsValid(query string) bool {
	return len(v.ValidateQuery(query)) == 0
}
