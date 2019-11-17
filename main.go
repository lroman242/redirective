package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/lroman242/redirective/controllers"
	"github.com/rs/cors"
)

func main() {
	certFile := flag.String("certPath", envString("CERT_PATH", ""), "Path to the certificate file | set this flag or env CERT_PATH")
	keyFile := flag.String("keyPath", envString("KEY_PATH", ""), "Path to the key file | set this flag or env KEY_PATH")

	//parse arguments
	flag.Parse()

	//run browser (google chrome headless)
	browserProcess := runBrowser()

	//stop browser (google chrome headless)
	defer func() {
		log.Printf("killing google-chrom PID %d\n", browserProcess.Pid)
		// Kill chrome:
		if err := browserProcess.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
	}()

	// start http server
	go func(certFile, keyFile string) {
		handler := makeHandler()

		err := http.ListenAndServe(":8080", *handler)
		if err != nil {
			log.Printf("ListenAndServe error: %s", err)
		}

		if certFile != "" && keyFile != "" {
			err = http.ListenAndServeTLS(":8083", certFile, keyFile, *handler)
			if err != nil {
				log.Printf("ListenAndServeTLS error: %s", err)
			}
		}
	}(*certFile, *keyFile)

	// awaiting to exit signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}

// Run google chrome headless
func runBrowser() *os.Process {
	// /usr/bin/google-chrome --addr=localhost --port=9222 --remote-debugging-port=9222 --remote-debugging-address=0.0.0.0 --disable-extensions --disable-gpu --headless --hide-scrollbars --no-first-run --no-sandbox
	cmd := exec.Command("/usr/bin/google-chrome", "--addr=localhost", "--port=9222", "--remote-debugging-port=9222", "--remote-debugging-address=0.0.0.0", "--disable-extensions", "--disable-gpu", "--headless", "--hide-scrollbars", "--no-first-run", "--no-sandbox")

	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("google-chrome headless runned with PID: %d\n", cmd.Process)
	log.Println("google-chrome headless runned on 9222 port")

	return cmd.Process
}

// Create web server handler
//  - define routes
//  - add CORS middleware
func makeHandler() *http.Handler {
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

	return &handler
}

// parse string value from os environment
// return default value if not found
func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}

	return def
}
