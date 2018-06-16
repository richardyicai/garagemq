package server

import (
	"github.com/valinurovam/garagemq/amqp"
	"github.com/valinurovam/garagemq/consumer"
	"github.com/valinurovam/garagemq/qos"
	"github.com/valinurovam/garagemq/interfaces"
)

func (channel *Channel) basicRoute(method amqp.Method) *amqp.Error {
	switch method := method.(type) {
	case *amqp.BasicQos:
		return channel.basicQos(method)
	case *amqp.BasicPublish:
		return channel.basicPublish(method)
	case *amqp.BasicConsume:
		return channel.basicConsume(method)
	case *amqp.BasicAck:
		return channel.basicAck(method)
	case *amqp.BasicCancel:
		return channel.basicCancel(method)

	}

	return amqp.NewConnectionError(amqp.NotImplemented, "unable to route basic method "+method.Name(), method.ClassIdentifier(), method.MethodIdentifier())
}

func (channel *Channel) basicQos(method *amqp.BasicQos) (err *amqp.Error) {
	channel.updateQos(method.PrefetchCount, method.PrefetchSize, method.Global)
	channel.sendMethod(&amqp.BasicQosOk{})

	return nil
}

func (channel *Channel) basicAck(method *amqp.BasicAck) (err *amqp.Error) {
	return channel.handleAck(method)
}

func (channel *Channel) basicPublish(method *amqp.BasicPublish) (err *amqp.Error) {
	if method.Immediate {
		return amqp.NewChannelError(amqp.NotImplemented, "Immediate = true", method.ClassIdentifier(), method.MethodIdentifier())
	}

	if _, err = channel.getExchangeWithError(method.Exchange, method); err != nil {
		return err
	}

	channel.currentMessage = amqp.NewMessage(method)
	return nil
}

func (channel *Channel) basicConsume(method *amqp.BasicConsume) (err *amqp.Error) {
	var qu interfaces.AmqpQueue

	if qu, err = channel.getQueueWithError(method.Queue, method); err != nil {
		return err
	}

	cmr := consumer.New(method.Queue, method.ConsumerTag, method.NoAck, channel, qu, []*qos.AmqpQos{channel.qos, channel.conn.qos})
	channel.addConsumer(cmr)

	if !method.NoWait {
		channel.sendMethod(&amqp.BasicConsumeOk{ConsumerTag: cmr.ConsumerTag})
	}

	cmr.Start()
	qu.AddConsumer(cmr)

	return nil
}

func (channel *Channel) basicCancel(method *amqp.BasicCancel) (err *amqp.Error) {
	if _, ok := channel.consumers[method.ConsumerTag]; !ok {
		return amqp.NewChannelError(amqp.NotFound, "Consumer not found", method.ClassIdentifier(), method.MethodIdentifier())
	}
	channel.removeConsumer(method.ConsumerTag)
	channel.sendMethod(&amqp.BasicConsumeOk{ConsumerTag: method.ConsumerTag})
	return nil
}
