package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/lroman242/redirective/config"
	"github.com/lroman242/redirective/infrastructure/logger"
	"github.com/lroman242/redirective/registry"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

func main() {
	conf := config.ParseConsole()

	r := registry.NewRegistry(conf)
	handler := r.NewHandler()

	// start http server
	go func(handler *http.Handler) {
		log.Println("Listening http on 8080")
		err := http.ListenAndServe(":8080", *handler)
		if err != nil {
			log.Printf("ListenAndServe error: %s", err)
		}
	}(handler)

}

// Create web server handler
//  - define routes
//  - add CORS middleware
func makeHandler() *http.Handler {
	router := httprouter.New()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// add routes
	router.GET("/api/find/:id", func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		id := ps.ByName("id")
		logger.Printf("[%s] Find: %s", time.Now().Format(time.RFC3339), id)
		controllers.LoadTraceResults(writer, request, col, id)
	})
	router.GET("/api/screenshot/chrome", func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		logger.Printf("[%s] Screenshot request: %s", time.Now().Format(time.RFC3339), request.URL.Query().Get("url"))
		controllers.ChromeScreenshot(writer, request, screenshotsStoragePath)
	})
	router.GET("/api/trace/chrome", func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		logger.Printf("[%s] Trace request: %s", time.Now().Format(time.RFC3339), request.URL.Query().Get("url"))
		controllers.ChromeTrace(writer, request, screenshotsStoragePath, col)
	})

	// Serve static files from the ./assets directory
	// http(s)://api.redirective.net/screenshots/{filename.png}
	router.NotFound = http.FileServer(http.Dir("assets/"))

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.Default().Handler(router)

	// Insert the middleware
	handler = c.Handler(handler)

	return &handler
}
