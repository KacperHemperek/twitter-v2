package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/router"
)

func main() {
	r := router.New()
	handler := api.ApplyCors(r)

	server := &http.Server{
		Handler:      handler,
		Addr:         "0.0.0.0:1337",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server listening on :1337")
	log.Fatalln(server.ListenAndServe())
}
