package amqputils

import (
	"math/rand"

	"github.com/streadway/amqp"
)

func Call(url string, queueName string, info []byte) ([]byte, error) {
	ch, q, close, err := Connect(url, queueName)
	defer close()
	if err != nil {
		return nil, err
	}

	qRec, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err != nil {
		return nil, err
	}

	corrId := randomString(32)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       qRec.Name,
			Body:          info,
		})
	if err != nil {
		return nil, err
	}

	resp := make(chan []byte)
	go Subscribe(ch, &qRec, func(ch *amqp.Channel, d amqp.Delivery) {
		if corrId == d.CorrelationId {
			resp <- d.Body
		}
	})
	return <-resp, nil
}

func Publish(url string, queueName string, info []byte) error {
	return PublishWithCorrelation(url, queueName, "", info)
}

func PublishWithCorrelation(url string, queueName string, corrId string, info []byte) error {
	ch, q, close, err := Connect(url, queueName)
	defer close()
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "application/json",
			Body:          info,
			CorrelationId: corrId,
		})
	if err != nil {
		return err
	}
	return nil
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
