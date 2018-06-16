package interfaces

import (
	"github.com/valinurovam/garagemq/amqp"
	"github.com/valinurovam/garagemq/qos"
)

type Queue interface {
	Push(item interface{})
	Pop() (res interface{})
	Length() int64
}

type AmqpQueue interface {
	Start()
	Push(message *amqp.Message)
	Pop() *amqp.Message
	PopQos(qosList []*qos.AmqpQos) *amqp.Message
	RemoveConsumer(cTag string)
	GetName() string
	IsExclusive() bool
	IsAutoDelete() bool
	IsDurable() bool
	ConnId() uint64
	Length() uint64
	ConsumersCount() int
	Purge() uint64
	AddConsumer(consumer Consumer)
	EqualWithErr(qu AmqpQueue) error
}

type Channel interface {
	SendContent(method amqp.Method, message *amqp.Message)
	NextDeliveryTag() uint64
	AddUnackedMessage(dTag uint64, cTag string, message *amqp.Message)
}

type Consumer interface {
	Consume()
	Tag() string
	Stop()
	Start()
}

type Binding interface {
	MatchDirect(exchange string, routingKey string) bool
	MatchFanout(exchange string) bool
	MatchTopic(exchange string, routingKey string) bool
	GetExchange() string
	GetRoutingKey() string
	GetQueue() string
	Equal(biding Binding) bool
}
