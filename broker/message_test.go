package broker_test

import (
	"testing"

	"github.com/sheirys/wsreply/broker"
	"github.com/stretchr/testify/assert"
)

func TestTranslateOp(t *testing.T) {
	testTable := []struct {
		Op       broker.Operand
		Expected string
	}{
		{
			Op:       broker.OpMessage,
			Expected: "OpMessage",
		},
		{
			Op:       broker.OpNewSubscriber,
			Expected: "OpNewSubscriber",
		},
		{
			Op:       broker.OpNoSubscribers,
			Expected: "OpNoSubscribers",
		},
		{
			Op:       broker.OpHasSubscribers,
			Expected: "OpHasSubscribers",
		},
		{
			Op:       broker.OpSyncSubscribers,
			Expected: "OpSyncSubscribers",
		},
		{
			Op:       -7,
			Expected: "unknown",
		},
	}

	for _, v := range testTable {
		got := broker.Message{
			Op: v.Op,
		}.TranslateOp()

		assert.Equal(t, got, v.Expected)
	}
}
