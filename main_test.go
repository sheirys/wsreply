package main

import (
	"testing"

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
			Expected: []byte(""),
		},
	}

	for i, v := range testTable {
		got := translate(v.Message)
		assert.Equal(t, got, v.Expected, "case %d failed", i)
	}
}
