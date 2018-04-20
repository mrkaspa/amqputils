package amqputils

import (
	"github.com/streadway/amqp"
)

// Server for receiving amqp messages
type Server struct {
	Version   string
	Event     string
	Do        SubscribeFunc
	AMQPChan  *amqp.Channel
	AMQPQueue *amqp.Queue
}

// NewServer creates one
func NewServer(version string, ch *amqp.Channel, event string, do SubscribeFunc) (*Server, error) {
	q, err := CreateQueue(ch, event)
	if err != nil {
		return nil, err
	}
	return &Server{
		Version:   version,
		Event:     event,
		Do:        do,
		AMQPChan:  ch,
		AMQPQueue: q,
	}, nil
}

// Start the server
func (s *Server) Start() {
	f := func(delivery amqp.Delivery) []byte {
		if s.Version == delivery.AppId {
			return s.Do(delivery)
		} else {
			return []byte("invalid message version, expecting version " + s.Version)
		}
	}
	Subscribe(s.AMQPChan, s.AMQPQueue, f)
}
