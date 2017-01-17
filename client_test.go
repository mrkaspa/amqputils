package amqputils

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	ch, q, close, _ := Connect("amqp://guest:guest@localhost", "demo")
	defer close()
	go Subscribe(ch, q, func(ch *amqp.Channel, d amqp.Delivery) {
	})
	err := Publish("amqp://guest:guest@localhost", "demo", []byte("xxx"))
	assert.Nil(t, err)
}

func TestCall(t *testing.T) {
	ch, q, close, _ := Connect("amqp://guest:guest@localhost", "echo")
	defer close()
	msg := []byte("xxx")
	go Subscribe(ch, q, func(ch *amqp.Channel, d amqp.Delivery) {
		ch.Publish(
			"",        // exchange
			d.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          d.Body,
			})
	})
	resp, err := Call("amqp://guest:guest@localhost", "echo", msg)
	assert.Nil(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, resp, msg)
}
