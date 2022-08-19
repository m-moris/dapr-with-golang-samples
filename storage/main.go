package main

import (
	"context"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"
)

// https://docs.dapr.io/reference/components-reference/supported-bindings/blobstorage/

func main() {

	fmt.Printf("hello world")

	ctx := context.Background()
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	str := "this is test"

	data := []byte(str)

	blobName := "foo/bar/boo.txt"

	// Uload blob
	fmt.Printf("Upload blob. %s\n", blobName)
	in := &dapr.InvokeBindingRequest{
		Name:      "myblob",
		Operation: "create",
		Data:      data,
		Metadata:  map[string]string{"blobName": blobName}}
	err = client.InvokeOutputBinding(ctx, in)

	if err != nil {
		panic(err)
	}

	// Download
	fmt.Printf("Download blob. %s\n", blobName)
	in = &dapr.InvokeBindingRequest{
		Name:      "myblob",
		Operation: "get",
		Metadata:  map[string]string{"blobName": blobName, "includeMetadat": "true"}}
	event, err := client.InvokeBinding(ctx, in)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Get Data = '%v'\n", string(event.Data))

	// List
	fmt.Print("List Blobs\n")
	in = &dapr.InvokeBindingRequest{
		Name:      "myblob",
		Operation: "list"}
	event, err = client.InvokeBinding(ctx, in)

	if err != nil {
		panic(err)
	}
	fmt.Printf("List Data = '%v'\n", string(event.Data))

	// Delete
	in = &dapr.InvokeBindingRequest{
		Name:      "myblob",
		Operation: "delete",
		Metadata:  map[string]string{"blobName": "foo/bar/boo.txt", "deleteSnapshots": "include"}}
	err = client.InvokeOutputBinding(ctx, in)

	if err != nil {
		panic(err)
	}

	fmt.Printf("deleted blob\n")

}
