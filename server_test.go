package amqputils

import (
	"log"
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
	server, err := NewServer("v1.0", ch, queue,
		func(d amqp.Delivery) []byte {
			return response
		},
		func(v1, v2 string) (bool, []byte) {
			if v1 != v2 {
				log.Printf("Version error. Expecting %s, received %s", v1, v2)
				return false, []byte("{\"error\":\"invalid message version, expecting " + v1 + "\"}")
			}
			return true, nil
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
	msg := []byte("xxx")
	server, close, _ := createServerTest("demo", msg)
	defer close()
	go server.Start()
	resp, err := Call(connString, server.Event, "v1.0", msg)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestServer_DoesntRespondWhenReturnNil(t *testing.T) {
	server, close, _ := createServerTest("demo", nil)
	defer close()
	go server.Start()
	_, err := Call(connString, server.Event, "v1.0", []byte("xxx"))
	assert.NotNil(t, err)
}
