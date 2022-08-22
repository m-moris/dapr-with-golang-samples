package main

import (
	"context"
	"log"

	"net/http"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

func main() {

	// Create Service
	s := daprd.NewService(":5005")

	// Bind cron service
	if err := s.AddBindingInvocationHandler("run", runHandler); err != nil {
		log.Fatalf("error adding binding handler: %v", err)
	}

	// start the service
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting service: %v", err)
	}
}

func runHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {

	log.Printf("Handler was invoked\n")
	log.Printf("Binding - Metadata:%v, Data:%v", in.Metadata, in.Data)

	return nil, nil
}
