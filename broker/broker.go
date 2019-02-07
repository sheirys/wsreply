package broker

// Broker defines publish/subscribe system. In general this is simple fan
// out system but with some rules. When new sub
type Broker interface {

	// Start should initialize all variables required by broker. This
	// should be called before usage of other broker functions.
	Start() error

	// Publish should send message to all subscribers in broker.
	Publish(Message) error

	// NewPublisher should return new publishers stream, where messages
	// produced by broker should appear e.g. in this stream notifications
	// about new subcscriber connection or if there is no other
	// subscribers left should be pushed.
	NewPublisherStream() (*Stream, error)

	// NewSubscriberStream should return new subscribers stream where
	// messages pushed by publishers should appear.
	NewSubscriberStream() (*Stream, error)

	// Unsubscribe should close passed stream.
	Unsubscribe(*Stream) error

	// Stop should unsubscribe all streams and stop broker.
	Stop() error
}
