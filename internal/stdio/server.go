package stdio

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
)

type Server struct {
	transport thrift.TTransport
	protocol  thrift.TProtocol
	processor thrift.TProcessor
}

func NewServer(processor thrift.TProcessor, opts ...serverOption) *Server {
	transport := NewTransport()
	defaultServer := &Server{
		transport: transport,
		protocol:  thrift.NewTJSONProtocol(transport),
		processor: processor,
	}
	for _, opt := range opts {
		opt(defaultServer)
	}
	return defaultServer
}

func (s *Server) Serve(ctx context.Context) error {
	errChan := make(chan error)

	go func() {
		for {
			ok, err := s.processor.Process(ctx, s.protocol, s.protocol)
			if err != nil {
				errChan <- err
				return
			}
			if !ok {
				continue
			}
			if err := s.transport.Flush(ctx); err != nil {
				errChan <- err
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

type serverOption func(*Server)

func WithProtocol(factory func(thrift.TTransport) thrift.TProtocol) serverOption {
	return func(s *Server) {
		s.protocol = factory(s.transport)
	}
}
