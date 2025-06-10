package shell

import (
	"os"
	"strings"
	"unicode"
	"fmt"
)

// Tokenize takes an input command string and returns a slice of parsed tokens.
func Tokenize(input string) ([]string, error) {
	var tokens []string			// Final list of parsed tokens
	var current strings.Builder // Current token being constructed
	var insideSingleQuote bool  // Track if inside single quotes
	var insideDoubleQuote bool  // Track if inside double quotes
	var escapeNext bool 		// Track if the next character should be escaped

	for _, r := range input {
		switch {
			// Handle escape character "\"
		case escapeNext:
			current.WriteRune(r)

		case r == '\\':
			escapeNext = true

		// Handle single quotes (literal strings)
		case r == '\'' && !insideDoubleQuote:
			if insideSingleQuote {
				// End single-quoted string
				insideSingleQuote = false
			} else {
				// Begin single-quoted string
				insideSingleQuote = true
			}
		
		// Handle double quotes (strings with expansion)
		case r == '"' && !insideSingleQuote:
			if insideDoubleQuote {
				// End double-quoted string
				insideDoubleQuote = false
			} else {
				insideDoubleQuote = true
			}

		// Handle spaces (token delimiters) if not in quote
		case unicode.IsSpace(r) && !insideSingleQuote && !insideDoubleQuote:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		// Anything else: normal characters
		default:
			current.WriteRune(r)
		}
	}

	// Add last token if any
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	// Handle unbalanced quotes
	if insideSingleQuote || insideDoubleQuote {
		return nil, fmt.Errorf("unbalanced quotes in input")
	}

	// Exapnd variables in tokens
	tokens = expandVariables(tokens)

	return tokens, nil
}

// expandVariables replaces any $VAR with its value in the environment.
func expandVariables(tokens []string) []string {
	var expanded []string
	for _, token := range tokens {
		if strings.Contains(token, "$") {
			expanded = append(expanded, os.Expand(token, os.Getenv))
		} else {
			expanded = append(expanded, token)
		}
	}
	return expanded
}