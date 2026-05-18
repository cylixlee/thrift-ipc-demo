package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/cylixlee/thrift-ipc-demo/internal/graceful"
	"github.com/cylixlee/thrift-ipc-demo/internal/stdio"
	"github.com/cylixlee/thrift-ipc-demo/internal/thrift/hello"
)

type helloHandler struct{}

func (helloHandler) Hello(ctx context.Context, request *hello.HelloRequest) (*hello.HelloResponse, error) {
	if request == nil {
		return &hello.HelloResponse{Msg: "I dont see a req"}, nil
	}
	return &hello.HelloResponse{Msg: fmt.Sprintf("(go) Hello, %s", request.Name)}, nil
}

func main() {
	server := stdio.NewServer(hello.NewHelloProcessor(helloHandler{}))

	fmt.Fprintln(os.Stderr, "(go) Thrift over stdio serving...")
	defer fmt.Fprintln(os.Stderr, "(go) Thrift over stdio shutting down...")
	if err := graceful.Run(server.Serve); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}
		log.Fatalln(err)
	}
}

// test data:
//
// [1,"hello",1,1,{"1":{"rec":{"1":{"str":"CYLIX"}}}}]
