package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleHello)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
func handleHello(w http.ResponseWriter, _ *http.Request) {
	wc, err := w.Write([]byte("Hello world"))
	if err != nil {
		log.Fatal("error writing response", err)
		return
	}
	fmt.Println(wc)

}
