package amqputils

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func createServerTest() (*Server, error) {
	return NewServer("amqp://guest:guest@localhost", "demo", func(d amqp.Delivery) []byte {
		return []byte("xxx")
	})
}

func TestNewServer(t *testing.T) {
	server, err := createServerTest()
	assert.Nil(t, err)
	assert.NotNil(t, server)
}

func TestServer_Start(t *testing.T) {
	server, _ := createServerTest()
	go server.Start()
	resp, err := Call(server.URL, server.Event, []byte("xxx"))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestServer_Stop(t *testing.T) {
	server, _ := NewServer("amqp://guest:guest@localhost", "stop", func(d amqp.Delivery) []byte {
		return []byte("xxx")
	})
	go server.Start()
	server.Stop()
	_, err := Call(server.URL, server.Event, []byte("xxx"))
	assert.NotNil(t, err)
}
