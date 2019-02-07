package broker

type Operand int

const (
	OpNewSubscriber Operand = iota
	OpNoSubscribers
	OpMessage
)

type Message struct {
	Op      Operand `json:"op"`
	Payload []byte  `json:"payload"`
}
