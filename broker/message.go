package broker

type Operand int

const (
	OpNewSubscriber Operand = iota
	OpNoSubscribers
	OpHasSubscribers
	OpSyncSubscribers
	OpMessage
)

type Message struct {
	Op      Operand `json:"op"`
	Payload []byte  `json:"payload"`
}

func (m Message) TranslateOp() string {
	switch m.Op {
	case OpNewSubscriber:
		return "OpNewSubscriber"
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

func MsgNewSubscriber() Message {
	return Message{
		Op: OpNewSubscriber,
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

func MsgMessage(data []byte) Message {
	return Message{
		Op:      OpMessage,
		Payload: data,
	}
}
