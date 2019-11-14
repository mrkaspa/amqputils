package amqputils

import (
	"github.com/streadway/amqp"
)

// Server for receiving amqp messages
type Server struct {
	Event     string
	Do        SubscribeFunc
	AMQPChan  *amqp.Channel
	AMQPQueue *amqp.Queue
	PoolSize  int
}

// NewServer creates one
func NewServer(ch *amqp.Channel, event string, do SubscribeFunc, poolSize int) (*Server, error) {
	q, err := CreateQueue(ch, event)
	if err != nil {
		return nil, err
	}
	return &Server{
		Event:     event,
		Do:        do,
		AMQPChan:  ch,
		AMQPQueue: q,
		PoolSize:  poolSize,
	}, nil
}

// HealtCheck send a response each time that receive a message
func HealtCheck(ch *amqp.Channel, queueService string) (*Server, error) {
	q, err := CreateQueue(ch, queueService)
	if err != nil {
		return nil, err
	}
	return &Server{
		Event: queueService,
		Do: func(m amqp.Delivery) ([]byte, error) {
			return []byte(queueService + " active "), nil
		},
		AMQPChan:  ch,
		AMQPQueue: q,
	}, nil
}

// Start the server
func (s *Server) Start() {
	f := func(delivery amqp.Delivery) ([]byte, error) {
		return s.Do(delivery)
	}
	Subscribe(s.AMQPChan, s.AMQPQueue, f, s.PoolSize)
}
