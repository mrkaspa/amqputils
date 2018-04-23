package amqputils

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

const connString = "amqp://guest:guest@localhost"

func createServerTest(queue string, response []byte) (*Server, func(), error) {
	ch, close, err := CreateConnection(connString)
	if err != nil {
		return nil, nil, err
	}
	server, err := NewServer("v.1.0", ch, queue, func(d amqp.Delivery) []byte {
		return response
	})
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
	server, close, _ := createServerTest("demo", []byte("xxx"))
	defer close()
	go server.Start()
	resp, err := Call(connString, server.Event, []byte("xxx"))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestServer_DoesntRespondWhenReturnNil(t *testing.T) {
	server, close, _ := createServerTest("demo", nil)
	defer close()
	go server.Start()
	msg, err := Call(connString, server.Event, []byte("xxx"))
	assert.Nil(t, err)

	m := "{\"error\":\"invalid message version, expecting version v.1.0\"}"
	assert.True(t, string(msg) == m)
}
