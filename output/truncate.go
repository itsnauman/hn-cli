package output

import (
	"fmt"
	"regexp"
	"strings"
)

const DefaultTruncateLen = 300

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

// StripHTML removes HTML tags and decodes common entities.
func StripHTML(s string) string {
	s = htmlTagRe.ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#x27;", "'")
	s = strings.ReplaceAll(s, "&#x2F;", "/")
	s = strings.ReplaceAll(s, "&#39;", "'")
	// Collapse whitespace
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimSpace(s)
}

// MakeTruncator returns a truncation function configured with the given max length.
// If full is true, the returned function only strips HTML but does not truncate.
func MakeTruncator(maxLen int, full bool) func(string) string {
	return func(s string) string {
		s = StripHTML(s)
		runes := []rune(s)
		if full || len(runes) <= maxLen {
			return s
		}
		return fmt.Sprintf("%s (truncated, %d chars total)", string(runes[:maxLen]), len(runes))
	}
}
