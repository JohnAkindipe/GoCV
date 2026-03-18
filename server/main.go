package main

import (
	"log"
	"net/http"
)

func main() {
	srvrPtr := &http.Server{
		Addr:    ":4000",
		Handler: routes(),
	}

	log.Printf("Starting server on %s", srvrPtr.Addr)
	err := srvrPtr.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
