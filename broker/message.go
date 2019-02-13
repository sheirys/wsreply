package broker

// Operand type defines operand.
type Operand int

// here possible opernads are defined. See more documentation about opernads in
// README.md file or broker/docs.go.
const (
	OpNoSubscribers Operand = iota
	OpHasSubscribers
	OpSyncSubscribers
	OpMessage
)

// Message is used in broker as intercomm protocol and will be used to accept
// messages from publishers and broadcast to subscribers.
type Message struct {
	Op      Operand `json:"op"`
	Payload string  `json:"payload"`
}

// TranslateOp will translate operand to string.
func (m Message) TranslateOp() string {
	switch m.Op {
	case OpNoSubscribers:
		return "OpNoSubscribers"
	case OpHasSubscribers:
		return "OpHasSubscribers"
	case OpSyncSubscribers:
		return "OpSyncSubscribers"
	case OpMessage:
		return "OpMessage"
	default:
		return "unknown"
	}
}

// MsgNoSubscribers produces message with OpNoSubscribers operand.
func MsgNoSubscribers() Message {
	return Message{
		Op: OpNoSubscribers,
	}
}

// MsgHasSubscribers produces message with OpHasSubscribers operand.
func MsgHasSubscribers() Message {
	return Message{
		Op: OpHasSubscribers,
	}
}

// MsgSyncSubscribers produces message with OpSyncSubscribers operand.
func MsgSyncSubscribers() Message {
	return Message{
		Op: OpSyncSubscribers,
	}
}

// MsgMessage produces message with defined payload that will be broadcasted to
// subscribers.
func MsgMessage(data string) Message {
	return Message{
		Op:      OpMessage,
		Payload: data,
	}
}
