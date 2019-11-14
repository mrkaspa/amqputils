package amqputils

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

const connString = "amqp://guest:guest@localhost"

func createServerTest(queue string, response []byte) (*Server, func(), error) {
	ch, close, err := CreateChannelConnection(connString)
	if err != nil {
		return nil, nil, err
	}
	server, err := NewServer(ch, queue,
		func(d amqp.Delivery) ([]byte, error) {
			return response, nil
		}, 1)

	if err != nil {
		return nil, nil, err
	}

	return server, close, nil
}

func TestNewServer(t *testing.T) {
	server, close, err := createServerTest("demo", []byte("xxx"))
	defer close()
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestServer_Start(t *testing.T) {
	msg := []byte("xxx")
	server, close, _ := createServerTest("demo", msg)
	defer close()
	go server.Start()
	resp, err := Call(connString, server.Event, msg, server.PoolSize)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestServer_DoesntRespondWhenReturnNil(t *testing.T) {
	server, close, _ := createServerTest("demo", nil)
	defer close()
	go server.Start()
	_, err := Call(connString, server.Event, []byte("xxx"), server.PoolSize)
	assert.NotNil(t, err)
}
