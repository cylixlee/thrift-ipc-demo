package graceful

import (
	"context"
	"os"
	"os/signal"
)

func Run(f func(context.Context) error) error {
	interruptChan := make(chan os.Signal, 1)
	errChan := make(chan error)
	signal.Notify(interruptChan, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		errChan <- f(ctx)
	}()
	<-interruptChan

	select {
	case err := <-errChan:
		return err
	default:
		return ctx.Err()
	}
}
