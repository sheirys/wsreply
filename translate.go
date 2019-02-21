package wsreply

import "strings"

// Translate will change "!" symbols to "?" in provided message.
func Translate(m string) string {
	return strings.Replace(m, "!", "?", -1)
}
