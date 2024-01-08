package mdutil

import (
	"regexp"
	"strings"
)

var (
	whitespaceOnly    = regexp.MustCompile("(?m)^[ \t]+$")
	leadingWhitespace = regexp.MustCompile("(?m)(^[ \t]*)(?:[^ \t\n])")
	// goDocListBullet   = regexp.MustCompile("^\t(\\*|\\+|-|•)").
)

// Dedent removes any common leading whitespace from every line in text.
//
// This can be used to make multiline strings to line up with the left edge of
// the display, while still presenting them in the source code in indented
// form.
func Dedent(text string) string {
	var margin string

	if text[0] == '\n' {
		text = whitespaceOnly.ReplaceAllString(text[1:], "")
	} else {
		text = whitespaceOnly.ReplaceAllString(text, "")
	}
	indents := leadingWhitespace.FindAllStringSubmatch(text, -1)

	// Look for the longest leading string of spaces and tabs common to all
	// lines.
	for i, indent := range indents {
		if i == 0 { //nolint: gocritic
			margin = indent[1]
		} else if strings.HasPrefix(indent[1], margin) {
			// Current line more deeply indented than previous winner:
			// no change (previous winner is still on top).
			continue
		} else if strings.HasPrefix(margin, indent[1]) {
			// Current line consistent with and no deeper than previous winner:
			// it's the new winner.
			margin = indent[1]
		} else {
			// Current line and previous winner have no common whitespace:
			// there is no margin.
			margin = ""
			break
		}
	}

	if margin != "" {
		text = regexp.MustCompile("(?m)^"+margin).ReplaceAllString(text, "")
	}
	return text
}

func ToMarkdown(s string) string {
	s = strings.ReplaceAll(s, "__CODE_SPAN__", "`")
	s = strings.ReplaceAll(s, "__CODE_BLOCK__", "```")
	s = Dedent(s)
	s = strings.ReplaceAll(s, "\t", "  ")
	return s
}
