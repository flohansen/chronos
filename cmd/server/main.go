package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "http_requests_total %d", rand.Int())
	})
	log.Fatal(http.ListenAndServe(":3000", mux))
}
