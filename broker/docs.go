package broker

// This package implements broker that satisfy requirements:
//
// * Has subscriber and publisher.
// * Informs publisher if there is any subscribers left.
// * Broadcast publishers message for all subscribers.
// * Allow to have multiple subscribers.
// * Allow to have multiple publishers.
//
// Communication between publisher, broker and subscriber is made by custom
// broker protocol defined in message.go file. Each message has defined operand
// and payload. Possible message with operands:
//
// * OpNewSubscriber (0) - will be broadcasted by broker to all publishers when
// new subscriber connects to broker. This can be used in case you want to
// trigger additional logic for publishers.
//
// * OpNoSubscriber (1) - will be broadcasted by broker to all publishers if no
// subscribers left in broker. On every subscribers disconnection, broker will
// check if any subscibers are left, and if no this message will be broadcasted
// to publishers. Also this message will be broadcasted as response to
// OpSyncSubscribers.
//
// * OpHasSubscribers (2) - will be broadcasted to all publishers everytime when
// subscriber disconnects and there is still some subscribers left in broker.
// Also this message will be broadcasted as response to OpSyncSubscribers.
//
// * OpSyncSubscribers (3) - can be sent by publisher to check if there is any
// subscribers left in broker. This message usualy will be sent by publisher on
// startup.
//
// * OpMessage (4) - this is the only message that can be accepted by subscibers.
// This message will be broadcasted for all subscribers by publisher.
