package amqputils

import "github.com/streadway/amqp"

// Server for receiving amqp messages
type Server struct {
	Event     string
	Do        SubscribeFunc
	AMQPChan  *amqp.Channel
	AMQPQueue *amqp.Queue
}

// NewServer creates one
func NewServer(ch *amqp.Channel, event string, do SubscribeFunc) (*Server, error) {
	q, err := CreateQueue(ch, event)
	if err != nil {
		return nil, err
	}
	return &Server{
		Event:     event,
		Do:        do,
		AMQPChan:  ch,
		AMQPQueue: q,
	}, nil
}

// Start the server
func (s *Server) Start() {
	Subscribe(s.AMQPChan, s.AMQPQueue, s.Do)
}
