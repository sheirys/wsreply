package broker

type Operand int

const (
	OpNewSubscriber Operand = iota
	OpNoSubscribers
)

type Message struct {
	Op      Operand `json:"op"`
	Payload string  `json:"string"`
}
