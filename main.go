package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("bbbbbb")
	fmt.Fprintf(w, "Hello, from Docker container!")
	log.Println("aaaaaa")
}
