package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Cron handler was invoked.")
	log.Printf("%v", r)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	w.Write([]byte("success"))
}

func main() {
	http.HandleFunc("/run", handler)
	err := http.ListenAndServe(":5005", nil)
	if err != nil {
		panic(err)
	}
}
