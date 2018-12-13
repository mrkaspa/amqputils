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
	go Subscribe(ch, q, func(d amqp.Delivery) ([]byte, error) {
		return nil, nil
	})
	err = Publish("amqp://guest:guest@localhost", "demo", []byte("xxx"))
	assert.NoError(t, err)
}

func TestCall(t *testing.T) {
	ch, close, err := CreateChannelConnection("amqp://guest:guest@localhost")
	defer close()
	assert.NoError(t, err)
	q, err := CreateQueue(ch, "echo")
	assert.NoError(t, err)
	msg := []byte("xxx")
	go Subscribe(ch, q, func(d amqp.Delivery) ([]byte, error) {
		return d.Body, nil
	})
	resp, err := Call("amqp://guest:guest@localhost", "echo", msg)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Equal(t, resp, msg)
}
