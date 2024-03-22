package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNum string = ":9000"

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "received request for method %s", r.Method)
	fmt.Println("request method", r.Method)
}

func main() {
	http.HandleFunc("/", root)

	log.Println("Started on port", portNum)
	fmt.Println("To close connection CTRL+C :-)")

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
