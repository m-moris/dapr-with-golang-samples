package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

func main() {

	ctx := context.Background()
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Enqueue
	for i := 0; i < 10; i++ {

		msg := fmt.Sprintf("Hello dapr queue = %d", i)
		req := &dapr.InvokeBindingRequest{
			Name:      "myqueue",
			Operation: "create",
			Data:      []byte(msg),
		}
		err = client.InvokeOutputBinding(ctx, req)

		if err != nil {
			panic(err)
		}
		log.Println("send message")
	}

	// Dequeue
	s, err := daprd.NewService(":5005")

	if err != nil {
		panic(err)
	}

	if err := s.AddBindingInvocationHandler("myqueue", runHandler); err != nil {
		log.Fatalf("error adding binding handler: %v", err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting service: %v", err)
	}
}

func runHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {

	log.Printf("Handler was invoked\n")
	//log.Printf("Binding - Metadata:%v, Data:%v", in.Metadata, in.Data)
	log.Printf("%s", string(in.Data))

	return nil, nil
}
