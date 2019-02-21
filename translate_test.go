package wsreply_test

import (
	"testing"

	"github.com/sheirys/wsreply"
	"github.com/stretchr/testify/assert"
)

func TestTranslate(t *testing.T) {
	testTable := []struct {
		Message  string
		Expected string
	}{
		{
			Message:  "labas!",
			Expected: "labas?",
		},
		{
			Message:  "!!!",
			Expected: "???",
		},
		{
			Message:  "123",
			Expected: "123",
		},
		{
			Message:  "",
			Expected: "",
		},
	}

	for i, v := range testTable {
		got := wsreply.Translate(v.Message)
		assert.Equal(t, got, v.Expected, "case %d failed", i)
	}
}
