package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lroman242/redirective/controllers"
	"github.com/rs/cors"
)

func main() {
	screenshotsStoragePath := flag.String("screenshotsPath", envString("SCREENSHOTS_PATH", "assets"), "Path to directory where screenshots would be stored | set this flag or env SCREENSHOTS_PATH")
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

	err := checkScreenshotsStorageDir(*screenshotsStoragePath)
	if err != nil {
		log.Fatalln("folder to store screenshots not found and couldn`t be created")
	}

	if !strings.HasSuffix(*screenshotsStoragePath, "/") {
		*screenshotsStoragePath = *screenshotsStoragePath + "/"
	}

	// start http server
	go func(certFile, keyFile, screenshotsStoragePath string) {
		handler := makeHandler(screenshotsStoragePath)

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
	}(*certFile, *keyFile, *screenshotsStoragePath)

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
func makeHandler(screenshotsStoragePath string) *http.Handler {
	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// add routes
	mux.HandleFunc("/api/screenshot/chrome", func(writer http.ResponseWriter, request *http.Request) {
		controllers.ChromeScreenshot(writer, request, screenshotsStoragePath)
	})
	mux.HandleFunc("/api/trace/chrome", func(writer http.ResponseWriter, request *http.Request) {
		controllers.ChromeTrace(writer, request, screenshotsStoragePath)
	})

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

func checkScreenshotsStorageDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	_, err := isWritable(path)
	if err != nil {
		return err
	}

	return nil
}

func isWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, errors.New("path isn't a directory")
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, errors.New("write permission bit is not set on this file for user")
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		return false, errors.New("unable to get stat. error " + err.Error())
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, errors.New("user doesn't have permission to write to this directory")
	}

	return true, nil
}
