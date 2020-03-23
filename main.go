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
	"time"

	"github.com/lroman242/redirective/controllers"
	"github.com/rs/cors"
)

func main() {
	screenshotsStoragePath := flag.String("screenshotsPath", envString("SCREENSHOTS_PATH", "assets/screenshots"), "Path to directory where screenshots would be stored | set this flag or env SCREENSHOTS_PATH")
	certFile := flag.String("certPath", envString("CERT_PATH", ""), "Path to the certificate file | set this flag or env CERT_PATH")
	keyFile := flag.String("keyPath", envString("KEY_PATH", ""), "Path to the key file | set this flag or env KEY_PATH")
	logPath := flag.String("logPath", envString("LOG_PATH", "log/redirective.log"), "Path to the log file | set this flag or env LOG_PATH")

	logFile, err := os.Create(*logPath)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

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

	err = checkScreenshotsStorageDir(*screenshotsStoragePath)
	if err != nil {
		log.Fatalf("folder to store screenshots not found and couldn`t be created. error: %s", err)
	}

	if !strings.HasSuffix(*screenshotsStoragePath, "/") {
		*screenshotsStoragePath = *screenshotsStoragePath + "/"
	}

	handler := makeHandler(*screenshotsStoragePath, logger)

	// start http server
	go func(handler *http.Handler) {
		log.Println("Listening http on 8080")
		err := http.ListenAndServe(":8080", *handler)
		if err != nil {
			log.Printf("ListenAndServe error: %s", err)
		}
	}(handler)

	// start https server
	if *certFile != "" && *keyFile != "" {
		go func(certFile, keyFile string, handler *http.Handler) {
			log.Println("Listening http on 8083")
			err = http.ListenAndServeTLS(":8083", certFile, keyFile, *handler)
			if err != nil {
				log.Printf("ListenAndServeTLS error: %s", err)
			}
		}(*certFile, *keyFile, handler)
	}

	// awaiting to exit signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}

// Run google chrome headless
func runBrowser() *os.Process {
	// /usr/bin/google-chrome --addr=localhost --port=9222 --remote-debugging-port=9222 --remote-debugging-address=0.0.0.0 --disable-extensions --disable-gpu --headless --hide-scrollbars --no-first-run --no-sandbox --user-agent="Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/77.0.3854.3 Chrome/77.0.3854.3 Safari/537.36"
	cmd := exec.Command("/usr/bin/google-chrome",
		"--addr=localhost",
		"--port=9222",
		"--remote-debugging-port=9222",
		"--remote-debugging-address=0.0.0.0",
		"--disable-extensions",
		"--disable-gpu",
		"--headless",
		"--hide-scrollbars",
		"--no-first-run",
		"--no-sandbox",
		"--user-agent=Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

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
func makeHandler(screenshotsStoragePath string, logger *log.Logger) *http.Handler {
	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// add routes
	mux.HandleFunc("/api/screenshot/chrome", func(writer http.ResponseWriter, request *http.Request) {
		logger.Printf("[%s] Screenshot request: %s", time.Now().Format(time.RFC3339), request.URL.Query().Get("url"))
		controllers.ChromeScreenshot(writer, request, screenshotsStoragePath)
	})
	mux.HandleFunc("/api/trace/chrome", func(writer http.ResponseWriter, request *http.Request) {
		logger.Printf("[%s] Trace request: %s", time.Now().Format(time.RFC3339), request.URL.Query().Get("url"))
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
