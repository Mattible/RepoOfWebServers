package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {

	port := os.Getenv("WEBSERVER_PORT")
	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/", handler)
	fmt.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}