package util

import (
	"regexp"
	"strings"
)

func GlobToRegexp(glob string) (*regexp.Regexp, error) {
	var out strings.Builder
	out.WriteString("^")

	escaping := false
	inClass := false

	for _, c := range glob {
		switch c {
		case '\\':
			if escaping {
				out.WriteString("\\\\")
				escaping = false
			} else {
				escaping = true
			}
		case '*':
			if escaping {
				out.WriteString("\\*")
				escaping = false
			} else {
				out.WriteString(".*")
			}
		case '?':
			if escaping {
				out.WriteString("\\?")
				escaping = false
			} else {
				out.WriteString(".")
			}
		case '[':
			if escaping {
				out.WriteString("\\[")
				escaping = false
			} else {
				inClass = true
				out.WriteRune(c)
			}
		case ']':
			if escaping || !inClass {
				out.WriteString("\\]")
				escaping = false
			} else {
				inClass = false
				out.WriteRune(c)
			}
		case '{':
			if escaping {
				out.WriteString("\\{")
				escaping = false
			} else {
				out.WriteString("(")
			}
		case '}':
			if escaping {
				out.WriteString("\\}")
				escaping = false
			} else {
				out.WriteString(")")
			}
		case ',':
			if escaping {
				out.WriteString("\\,")
				escaping = false
			} else if inClass {
				out.WriteRune(c)
			} else {
				out.WriteString("|")
			}
		default:
			if escaping {
				out.WriteString("\\")
				escaping = false
			}

			// escape regex meta chars
			if strings.ContainsRune(".^$+()|", c) {
				out.WriteRune('\\')
			}
			out.WriteRune(c)
		}
	}

	out.WriteString("$")

	return regexp.Compile(out.String())
}
