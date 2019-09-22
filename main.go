package main

import (
	"github.com/lroman242/redirective/controllers"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	//start chrome
	// /usr/bin/google-chrome --addr=localhost --port=9222 --remote-debugging-port=9222 --remote-debugging-address=0.0.0.0 --disable-extensions --disable-gpu --headless --hide-scrollbars --no-first-run --no-sandbox
	cmd := exec.Command("/usr/bin/google-chrome", "--addr=localhost", "--port=9222", "--remote-debugging-port=9222", "--remote-debugging-address=0.0.0.0", "--disable-extensions", "--disable-gpu", "--headless", "--hide-scrollbars", "--no-first-run", "--no-sandbox")

	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("google-chrome headless runned with PID: %d\n", cmd.Process.Pid)
	log.Println("google-chrome headless runned on 9222 port")

	defer func() {
		log.Printf("killing google-chrom PID %d\n", cmd.Process.Pid)
		// Kill chrome:
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
	}()

	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// add routes
	mux.HandleFunc("/api/screenshot/chrome", controllers.ChromeScreenshot)
	mux.HandleFunc("/api/trace/chrome", controllers.ChromeTrace)

	fs := http.FileServer(http.Dir("assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.Default().Handler(mux)

	// Insert the middleware
	handler = c.Handler(handler)

	// start http server
	go func() {
		err = http.ListenAndServe(":8080", handler)
		if err != nil {
			log.Printf("ListenAndServe error: %s", err)
		}
		//TODO: dynamic host + certs
		err = http.ListenAndServeTLS(":8083", "/etc/letsencrypt/live/redirective.net/fullchain.pem", "/etc/letsencrypt/live/redirective.net/privkey.pem", handler)
		if err != nil {
			log.Printf("ListenAndServeTLS error: %s", err)
		}
	}()

	// awaiting to exit signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
