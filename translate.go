package wsreply

import "bytes"

// Translate will change "!" symbols to "?" in provided message.
func Translate(m []byte) []byte {
	return bytes.Replace(m, []byte("!"), []byte("?"), -1)
}
