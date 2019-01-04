package main

import (
	"github.com/lroman242/redirective/controllers"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	//start chrome
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

	// add routes
	http.HandleFunc("/api/screenshot/chrome", controllers.ChromeScreenshot)
	http.HandleFunc("/api/trace/chrome", controllers.ChromeTrace)

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// start http server
	go func() {
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Printf("ListenAndServe error: %s", err)
		}
	}()

	// awaiting to exit signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
