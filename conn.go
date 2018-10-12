package amqputils

import (
	"log"

	"github.com/streadway/amqp"
)

// SubscribeFunc function to handle an incoming message
type SubscribeFunc func(amqp.Delivery) ([]byte, error)

// CreateConnection connection - channel and its respective close function
func CreateConnection(url string) (*amqp.Connection, *amqp.Channel, func(), error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}

	close := func() {
		conn.Close()
		ch.Close()
	}
	return conn, ch, close, nil
}

// CreateChannelConnection channel and its respective close function
func CreateChannelConnection(url string) (*amqp.Channel, func(), error) {
	_, ch, close, err := CreateConnection(url)
	return ch, close, err
}

// CreateQueue in the amqp server
func CreateQueue(ch *amqp.Channel, queueName string) (*amqp.Queue, error) {
	return declareQueue(ch, queueName, true)
}

// CreateQueue in the amqp server not durable
func CreateQueueNotDurable(ch *amqp.Channel, queueName string) (*amqp.Queue, error) {
	return declareQueue(ch, queueName, false)
}

func declareQueue(ch *amqp.Channel, queueName string, durable bool) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		durable,   // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return nil, err
	}

	return &q, nil
}

// Subscribe to a queue and handle the messages
func Subscribe(ch *amqp.Channel, q *amqp.Queue, do SubscribeFunc) error {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		msg, err := do(d)

		if err != nil {
			log.Printf("AMQPUTILS, an error has occurred: %v", err.Error())
			continue
		}

		d.Ack(false)

		if msg != nil && d.ReplyTo != "" && d.CorrelationId != "" {
			ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          msg,
					AppId:         d.AppId,
				})
		}
	}
	return nil
}
