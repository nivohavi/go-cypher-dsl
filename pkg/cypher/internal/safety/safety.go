package safety

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	// Safety level for query generation
	safetyLevel = getLevelFromEnv()

	// Default logger for safety warnings
	safetyLogger = log.New(os.Stderr, "[cypher-safety] ", log.LstdFlags)
)

// SafetyLevel represents the safety enforcement level for the DSL
type SafetyLevel int

const (
	// SafetyLevelOff disables all safety warnings
	SafetyLevelOff SafetyLevel = iota

	// SafetyLevelWarn issues warnings but allows unsafe operations
	SafetyLevelWarn

	// SafetyLevelStrict prevents unsafe operations by panicking
	SafetyLevelStrict
)

// SetSafetyLevel sets the global safety level
func SetSafetyLevel(level SafetyLevel) {
	safetyLevel = level
}

// SetSafetyLogger sets a custom logger for safety warnings
func SetSafetyLogger(logger *log.Logger) {
	safetyLogger = logger
}

// WarnOnStringLiteral issues a warning when a string literal is used in a query
// Returns true if the operation should proceed
func WarnOnStringLiteral(value string, context string) bool {
	if safetyLevel == SafetyLevelOff {
		return true
	}

	// Check if value potentially contains special characters that could be used for injection
	isSuspicious := strings.ContainsAny(value, `'"();{}[]`)

	if isSuspicious {
		msg := fmt.Sprintf("String literal '%s' in %s could lead to injection vulnerabilities. Use parameters instead.", value, context)

		switch safetyLevel {
		case SafetyLevelWarn:
			safetyLogger.Println(msg)
			return true
		case SafetyLevelStrict:
			safetyLogger.Println(msg)
			panic(msg)
		}
	}

	return true
}

// getLevelFromEnv determines the safety level from environment variables
func getLevelFromEnv() SafetyLevel {
	level := os.Getenv("GO_CYPHER_DSL_SAFETY_LEVEL")
	switch strings.ToUpper(level) {
	case "OFF":
		return SafetyLevelOff
	case "WARN":
		return SafetyLevelWarn
	case "STRICT":
		return SafetyLevelStrict
	default:
		// Default to warn level
		return SafetyLevelWarn
	}
}
