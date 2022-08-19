package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

// https://docs.dapr.io/developing-applications/building-blocks/state-management/howto-get-save-state/

func main() {
	const STATE_STORE_NAME = "statestore2"
	client, err := dapr.NewClient()
	defer client.Close()
	if err != nil {
		panic(err)
	}

	err = client.DeleteState(context.Background(), STATE_STORE_NAME, "order_1", nil)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixMicro())

	for i := 0; i < 10; i++ {
		orderId := rand.Intn(1000-1) + 1

		ctx := context.Background()
		err = client.SaveState(ctx, STATE_STORE_NAME, "order_1", []byte(strconv.Itoa(orderId)), nil)
		if err != nil {
			panic(err)
		}

		item, err := client.GetState(ctx, STATE_STORE_NAME, "order_1", nil)
		if err != nil {
			panic(err)
		}

		log.Printf("Value=%s, Etag=%s\n", string(item.Value), item.Etag)
		time.Sleep(1 * time.Second)
	}
}
