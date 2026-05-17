package stdio

import (
	"bufio"
	"context"
	"io"
	"math"
	"os"
)

type StdioTransport struct {
	reader *bufio.Reader
	writer *bufio.Writer
	closed bool
}

func NewTransport() *StdioTransport {
	return &StdioTransport{
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}

func (t *StdioTransport) Close() error {
	t.closed = true
	return nil
}

func (t *StdioTransport) Read(p []byte) (n int, err error) {
	if t.closed {
		return 0, io.EOF
	}
	return t.reader.Read(p)
}

func (t *StdioTransport) Write(p []byte) (n int, err error) {
	if t.closed {
		return 0, io.EOF
	}
	return t.writer.Write(p)
}

func (t *StdioTransport) Flush(ctx context.Context) (err error) {
	if t.closed {
		return io.ErrClosedPipe
	}

	errChan := make(chan error)
	go func() {
		errChan <- t.writer.Flush()
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		return
	case err := <-errChan:
		return err
	}
}

func (StdioTransport) RemainingBytes() uint64 { return math.MaxUint64 }

func (t *StdioTransport) Open() error  { t.closed = false; return nil }
func (t *StdioTransport) IsOpen() bool { return !t.closed }
