package amqputils

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	ch, q, close, _ := Connect("amqp://guest:guest@localhost", "demo")
	defer close()
	go Subscribe(ch, q, func(d amqp.Delivery) []byte {
		return nil
	})
	err := Publish("amqp://guest:guest@localhost", "demo", []byte("xxx"))
	assert.Nil(t, err)
}

func TestCall(t *testing.T) {
	ch, q, close, _ := Connect("amqp://guest:guest@localhost", "echo")
	defer close()
	msg := []byte("xxx")
	go Subscribe(ch, q, func(d amqp.Delivery) []byte {
		return d.Body
	})
	resp, err := Call("amqp://guest:guest@localhost", "echo", msg)
	assert.Nil(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, resp, msg)
}
