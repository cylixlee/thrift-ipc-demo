package stdio

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
)

type StdioServer struct {
	transport thrift.TTransport
	protocol  thrift.TProtocol
	processor thrift.TProcessor
}

func NewServer(processor thrift.TProcessor, opts ...serverOption) *StdioServer {
	transport := NewTransport()
	defaultServer := &StdioServer{
		transport: transport,
		protocol:  thrift.NewTJSONProtocol(transport),
		processor: processor,
	}
	for _, opt := range opts {
		opt(defaultServer)
	}
	return defaultServer
}

func (s *StdioServer) Serve(ctx context.Context) error {
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

type serverOption func(*StdioServer)

func WithProtocol(factory func(thrift.TTransport) thrift.TProtocol) serverOption {
	return func(s *StdioServer) {
		s.protocol = factory(s.transport)
	}
}
