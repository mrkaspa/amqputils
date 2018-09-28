package amqputils

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	ch, close, err := CreateChannelConnection("amqp://guest:guest@localhost")
	defer close()
	assert.NoError(t, err)
	q, err := CreateQueue(ch, "demo")
	assert.NoError(t, err)
	go Subscribe(ch, q, func(d amqp.Delivery) []byte {
		return nil
	})
	err = Publish("amqp://guest:guest@localhost", "demo", "v1.0", []byte("xxx"))
	assert.NoError(t, err)
}

func TestCall(t *testing.T) {
	ch, close, err := CreateChannelConnection("amqp://guest:guest@localhost")
	defer close()
	assert.NoError(t, err)
	q, err := CreateQueue(ch, "echo")
	assert.NoError(t, err)
	msg := []byte("xxx")
	go Subscribe(ch, q, func(d amqp.Delivery) []byte {
		return d.Body
	})
	resp, err := Call("amqp://guest:guest@localhost", "echo", "v1.0", msg)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, resp, msg)
}
