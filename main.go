package main

import (
	"github.com/lroman242/redirective/config"
	"github.com/lroman242/redirective/registry"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.ParseConsole()

	r := registry.NewRegistry(conf)
	handler := r.NewHandler()

	// start http server
	go func(handler http.Handler) {
		log.Println("Listening http on 8080")

		err := http.ListenAndServe(":8080", handler)
		if err != nil {
			log.Printf("ListenAndServe error: %s", err)
		}
	}(handler)

	// awaiting to exit signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
