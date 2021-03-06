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

func TestMsgNoSubscribers(t *testing.T) {
	m := broker.MsgNoSubscribers()
	assert.Equal(t, m.Op, broker.OpNoSubscribers)
}

func TestMsgHasSubscribers(t *testing.T) {
	m := broker.MsgHasSubscribers()
	assert.Equal(t, m.Op, broker.OpHasSubscribers)
}

func TestMsgSyncSubscribers(t *testing.T) {
	m := broker.MsgSyncSubscribers()
	assert.Equal(t, m.Op, broker.OpSyncSubscribers)
}

func TestMsgMessage(t *testing.T) {
	data := "random_data"
	m := broker.MsgMessage(data)
	assert.Equal(t, m.Op, broker.OpMessage)
	assert.Equal(t, m.Payload, data)
}
