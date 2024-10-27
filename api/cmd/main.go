package main

import (
	"fmt"
	"github.com/kacperhemperek/twitter-v2/router"
	"log"
	"net/http"
	"time"
)

func main() {
	r := router.New()

	server := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:1337",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server listening on :1337")
	log.Fatalln(server.ListenAndServe())
}
