package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("started")

	http.HandleFunc("/", hello)

	srv := &http.Server{
		Addr: "localhost:8000",
	}

	err := srv.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		log.Fatalln(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, 世界")
}
