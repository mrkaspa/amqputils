package amqputils

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	t.Run("When the connection is valid", whenConnectIsValid)
	t.Run("When the connection is invalid", whenConnectIsInvalid)
}

func whenConnectIsValid(t *testing.T) {
	ch, close, err := CreateConnection("amqp://guest:guest@localhost")
	defer close()
	assert.NoError(t, err)
	assert.NotNil(t, ch)
}

func whenConnectIsInvalid(t *testing.T) {
	_, _, err := CreateConnection("amqp://guest:xxxx@localhost")
	assert.Error(t, err)
}

func TestSubscribe(t *testing.T) {
	ch, close, _ := CreateConnection("amqp://guest:guest@localhost")
	defer close()
	q, _ := CreateQueue(ch, "demo")
	resp := make(chan []byte)
	go Subscribe(ch, q, func(d amqp.Delivery) []byte {
		resp <- d.Body
		return nil
	})
	err := Publish("amqp://guest:guest@localhost", "demo", "v1.0", []byte("xxx"))
	assert.Nil(t, err)
	assert.NotNil(t, <-resp)
}
