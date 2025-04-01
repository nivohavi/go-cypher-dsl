package renderer

import (
	"fmt"
	"regexp"
	"strings"
)

// FormattingOptions contains configuration options for Cypher formatting
type FormattingOptions struct {
	IndentString     string // String used for indentation
	KeywordCase      KeywordCase
	ClauseNewline    bool // Whether to put each clause on a new line
	IndentSubClauses bool // Whether to indent subclauses
	MaxLineLength    int  // Maximum line length before wrapping (0 = no limit)
}

// KeywordCase defines how to format Cypher keywords
type KeywordCase int

const (
	// KeywordCaseAsIs leaves keywords as they are
	KeywordCaseAsIs KeywordCase = iota
	// KeywordCaseUpper converts all keywords to uppercase
	KeywordCaseUpper
	// KeywordCaseLower converts all keywords to lowercase
	KeywordCaseLower
)

// DefaultFormattingOptions returns a default set of formatting options
func DefaultFormattingOptions() FormattingOptions {
	return FormattingOptions{
		IndentString:     "  ",
		KeywordCase:      KeywordCaseUpper,
		ClauseNewline:    true,
		IndentSubClauses: true,
		MaxLineLength:    80,
	}
}

// CypherFormatter formats Cypher queries with customizable options
type CypherFormatter struct {
	options FormattingOptions
}

// NewCypherFormatter creates a new CypherFormatter with the given options
func NewCypherFormatter(options FormattingOptions) *CypherFormatter {
	return &CypherFormatter{
		options: options,
	}
}

// NewDefaultFormatter creates a new CypherFormatter with default options
func NewDefaultFormatter() *CypherFormatter {
	return NewCypherFormatter(DefaultFormattingOptions())
}

// Format formats a Cypher query string
func (f *CypherFormatter) Format(query string) string {
	// List of Cypher keywords that should trigger formatter rules
	keywords := []string{
		"MATCH", "OPTIONAL MATCH", "WHERE", "WITH", "RETURN", "ORDER BY",
		"SKIP", "LIMIT", "CREATE", "MERGE", "DELETE", "DETACH DELETE",
		"SET", "REMOVE", "UNWIND", "CALL", "YIELD", "UNION", "UNION ALL",
	}

	// First pass: normalize whitespace
	query = normalizeWhitespace(query)

	// Format keywords according to case preference
	for _, keyword := range keywords {
		var replacement string
		switch f.options.KeywordCase {
		case KeywordCaseUpper:
			replacement = strings.ToUpper(keyword)
		case KeywordCaseLower:
			replacement = strings.ToLower(keyword)
		default:
			replacement = keyword
		}

		// Replace only stand-alone keywords (not substrings of other words)
		// This uses a case-insensitive search to find keywords regardless of their current case
		regex := fmt.Sprintf(`(?i)(^|\s)%s(\s|$)`, regexp.QuoteMeta(keyword))
		query = regexp.MustCompile(regex).ReplaceAllStringFunc(query, func(match string) string {
			// Preserve the spacing around the keyword
			if strings.HasPrefix(match, " ") && strings.HasSuffix(match, " ") {
				return " " + replacement + " "
			} else if strings.HasPrefix(match, " ") {
				return " " + replacement
			} else if strings.HasSuffix(match, " ") {
				return replacement + " "
			}
			return replacement
		})
	}

	// Apply formatting rules if requested
	if f.options.ClauseNewline {
		// Replace keywords with newline + keyword
		for _, keyword := range keywords {
			if f.options.KeywordCase == KeywordCaseUpper {
				keyword = strings.ToUpper(keyword)
			} else if f.options.KeywordCase == KeywordCaseLower {
				keyword = strings.ToLower(keyword)
			}

			pattern := fmt.Sprintf(` %s `, keyword)
			replacement := fmt.Sprintf("\n%s ", keyword)
			query = strings.ReplaceAll(query, pattern, replacement)
		}

		// Split into lines and indent
		lines := strings.Split(query, "\n")
		for i := 1; i < len(lines); i++ {
			if f.options.IndentSubClauses {
				lines[i] = f.options.IndentString + lines[i]
			}
		}

		query = strings.Join(lines, "\n")
	}

	return query
}

// normalizeWhitespace reduces multiple spaces to a single space
func normalizeWhitespace(s string) string {
	// Replace tabs and newlines with spaces
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\n", " ")

	// Collapse multiple spaces into a single space
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}

	return strings.TrimSpace(s)
}
