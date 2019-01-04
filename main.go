package main

import (
	"github.com/lroman242/redirective/controllers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/screenshot/chrome", controllers.ChromeScreenshot)
	http.HandleFunc("/api/trace/chrome", controllers.ChromeTrace)

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("ListenAndServe error: %s", err)
	}
}
