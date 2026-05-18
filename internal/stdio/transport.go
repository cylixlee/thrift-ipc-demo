package stdio

import (
	"bufio"
	"context"
	"io"
	"math"
	"os"
)

type Transport struct {
	reader *bufio.Reader
	writer *bufio.Writer
	closed bool
}

func NewTransport() *Transport {
	return &Transport{
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}

func (t *Transport) Close() error {
	t.closed = true
	return nil
}

func (t *Transport) Read(p []byte) (n int, err error) {
	if t.closed {
		return 0, io.EOF
	}
	return t.reader.Read(p)
}

func (t *Transport) Write(p []byte) (n int, err error) {
	if t.closed {
		return 0, io.EOF
	}
	return t.writer.Write(p)
}

func (t *Transport) Flush(ctx context.Context) (err error) {
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

func (Transport) RemainingBytes() uint64 { return math.MaxUint64 }

func (t *Transport) Open() error  { t.closed = false; return nil }
func (t *Transport) IsOpen() bool { return !t.closed }
