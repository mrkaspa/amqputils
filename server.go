package amqputils

import "github.com/streadway/amqp"

// Server for receiving amqp messages
type Server struct {
	URL       string
	Event     string
	Do        SubscribeFunc
	Close     func()
	AMQPChan  *amqp.Channel
	AMQPQueue *amqp.Queue
}

// NewServer creates one
func NewServer(url, event string, do SubscribeFunc) (*Server, error) {
	ch, q, close, err := Connect(url, event)
	if err != nil {
		return nil, err
	}
	return &Server{
		URL:       url,
		Event:     event,
		Do:        do,
		Close:     close,
		AMQPChan:  ch,
		AMQPQueue: q,
	}, nil
}

// Start the server
func (s *Server) Start() {
	Subscribe(s.AMQPChan, s.AMQPQueue, s.Do)
}

// Stop the server
func (s *Server) Stop() {
	s.Close()
}
