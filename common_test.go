package wsreply_test

import (
	"testing"

	"github.com/sheirys/wsreply"
	"github.com/stretchr/testify/assert"
)

func TestTranslate(t *testing.T) {
	testTable := []struct {
		Message  []byte
		Expected []byte
	}{
		{
			Message:  []byte("labas?"),
			Expected: []byte("labas!"),
		},
		{
			Message:  []byte("???"),
			Expected: []byte("!!!"),
		},
		{
			Message:  []byte("123"),
			Expected: []byte("123"),
		},
		{
			Message:  []byte(""),
			Expected: []byte(nil),
		},
		{
			Message:  []byte(nil),
			Expected: []byte(nil),
		},
	}

	for i, v := range testTable {
		got := wsreply.Translate(v.Message)
		assert.Equal(t, got, v.Expected, "case %d failed", i)
	}
}
