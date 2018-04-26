package amqputils

import (
	"github.com/streadway/amqp"
)

type cmpFunc = func(string, string) (bool, []byte)

// Server for receiving amqp messages
type Server struct {
	Version        string
	Event          string
	Do             SubscribeFunc
	AMQPChan       *amqp.Channel
	AMQPQueue      *amqp.Queue
	CompareVersion cmpFunc
}

// NewServer creates one
func NewServer(version string, ch *amqp.Channel, event string, do SubscribeFunc, cFun cmpFunc) (*Server, error) {
	q, err := CreateQueue(ch, event)
	if err != nil {
		return nil, err
	}
	return &Server{
		Version:        version,
		Event:          event,
		Do:             do,
		AMQPChan:       ch,
		AMQPQueue:      q,
		CompareVersion: cFun,
	}, nil
}

// Start the server
func (s *Server) Start() {
	f := func(delivery amqp.Delivery) []byte {
		// fmt.Printf("msg \n%#v\n", delivery)
		if ok, data := s.CompareVersion(s.Version, delivery.AppId); !ok {
			return data
		}
		return s.Do(delivery)
	}
	Subscribe(s.AMQPChan, s.AMQPQueue, f)
}
