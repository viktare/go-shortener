package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr: ":3000",
		Handler: http.HandlerFunc(BasicHandler),
	}

	err := server.ListenAndServe();

	if err != nil {
		fmt.Println("failet to listen to server", err)
	}
}

func BasicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
