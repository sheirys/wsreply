package broker_test

import (
	"testing"

	"github.com/sheirys/wsreply/broker"
	"github.com/stretchr/testify/assert"
)

// TestBrokerStop will test if broker can be stopped. This will cause a
// deadlock if anything is wrong.
func TestBrokerStop(t *testing.T) {
	b := broker.NewInMemBroker(false)
	assert.NoError(t, b.Start())
	assert.NoError(t, b.Stop())
}
