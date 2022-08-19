package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

type ReceiveState struct {
	Count      int   `json:"count"`
	LastUpdate int64 `json:"lastUpdate"`
}

// https://docs.dapr.io/developing-applications/building-blocks/state-management/howto-get-save-state/

const STATE_STORE_NAME = "statestore2"
const KEY = "key"

func main() {

	ctx := context.Background()
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// ステート削除
	err = client.DeleteState(ctx, STATE_STORE_NAME, KEY, nil)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(5)
	go count(ctx, client, &wg)
	go count(ctx, client, &wg)
	go count(ctx, client, &wg)
	go count(ctx, client, &wg)
	go count(ctx, client, &wg)

	log.Println("waiting...")
	wg.Wait()
	log.Println("finish")
}

func deseralize(bin []byte) *ReceiveState {
	state := &ReceiveState{}
	_ = json.Unmarshal(bin, state)
	return state
}

func serialize(state *ReceiveState) []byte {
	bin, _ := json.Marshal(state)
	return bin
}

func count(ctx context.Context, client dapr.Client, wg *sync.WaitGroup) {

	for i := 1; i <= 10; i++ {

		var retry bool = true

		for retry {

			result, err := client.GetState(ctx, STATE_STORE_NAME, KEY, nil)
			if err != nil {
				panic(err)
			}

			var state *ReceiveState

			item := &dapr.SetStateItem{
				Key: KEY,
				Options: &dapr.StateOptions{
					Concurrency: dapr.StateConcurrencyFirstWrite,
					Consistency: dapr.StateConsistencyEventual,
				},
			}

			if len(result.Value) == 0 {
				state = &ReceiveState{
					Count:      1,
					LastUpdate: time.Now().Unix(),
				}
			} else {
				state = deseralize(result.Value)
				state.Count++
				state.LastUpdate = time.Now().Unix()
				item.Etag = &dapr.ETag{Value: result.Etag}
			}

			item.Value = serialize(state)

			err = client.SaveBulkState(ctx, STATE_STORE_NAME, item)
			if err == nil {
				retry = false
			} else if strings.Contains(err.Error(), "possible etag mismatch") {
				log.Println("mismatch etag. reload state.")
			} else {
				panic(err)
			}
			log.Printf("Count=%d, %s\n", state.Count, time.Unix(state.LastUpdate, 0))
		}
	}
	wg.Done()

}
