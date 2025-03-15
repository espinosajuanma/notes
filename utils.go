package notes

import (
	"regexp"

	"github.com/rwxrob/term"
)

func paint(s string, c string) string {
	rx := regexp.MustCompile("`([^`]+)`")
	text := rx.ReplaceAllFunc([]byte(s), func(match []byte) []byte {
		l := len(match)
		return []byte(c + string(match[1:l-1]) + term.Reset)
	})
	return string(text)
}
