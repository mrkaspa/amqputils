package amqputils

import "github.com/streadway/amqp"

// Connect to amqp server
func Connect(url string, queueName string) (*amqp.Channel, *amqp.Queue, func(), error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return nil, nil, nil, err
	}

	close := func() {
		conn.Close()
		ch.Close()
	}

	return ch, &q, close, nil
}

func Subscribe(ch *amqp.Channel, q *amqp.Queue, do func(*amqp.Channel, amqp.Delivery)) error {
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
		return nil
	}

	for d := range msgs {
		go do(ch, d)
	}
	return nil
}
