package broker

type Operand int

const (
	OpNoSubscribers Operand = iota
	OpHasSubscribers
	OpSyncSubscribers
	OpMessage
)

type Message struct {
	Op      Operand `json:"op"`
	Payload string  `json:"payload"`
}

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

func MsgNoSubscribers() Message {
	return Message{
		Op: OpNoSubscribers,
	}
}

func MsgHasSubscribers() Message {
	return Message{
		Op: OpHasSubscribers,
	}
}

func MsgSyncSubscribers() Message {
	return Message{
		Op: OpSyncSubscribers,
	}
}

func MsgMessage(data string) Message {
	return Message{
		Op:      OpMessage,
		Payload: data,
	}
}
