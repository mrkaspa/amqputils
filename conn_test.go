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
	ch, q, close, err := Connect("amqp://guest:guest@localhost", "demo")
	defer close()
	assert.NotNil(t, ch)
	assert.NotNil(t, q)
	assert.Nil(t, err)
}

func whenConnectIsInvalid(t *testing.T) {
	_, _, _, err := Connect("amqp://guest:xxxx@localhost", "demo")
	assert.NotNil(t, err)
}

func TestSubscribe(t *testing.T) {
	ch, q, close, _ := Connect("amqp://guest:guest@localhost", "demo")
	defer close()
	resp := make(chan []byte)
	go Subscribe(ch, q, func(d amqp.Delivery) []byte {
		resp <- d.Body
		return nil
	})
	err := Publish("amqp://guest:guest@localhost", "demo", []byte("xxx"))
	assert.Nil(t, err)
	assert.NotNil(t, <-resp)
}
