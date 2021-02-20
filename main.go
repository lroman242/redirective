package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/lroman242/redirective/config"
	"github.com/lroman242/redirective/registry"
)

func main() {
	conf := config.Parse()

	r := registry.NewRegistry(conf)
	handler := r.NewHandler()

	// start http server
	go func(handler http.Handler) {
		log.Println("Listening http on " + strconv.Itoa(conf.HTTPServer.Port))

		err := http.ListenAndServe(conf.HTTPServer.Host+":"+strconv.Itoa(conf.HTTPServer.Port), handler)
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
