package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/lroman242/redirective/interface/api/controllers"
	"github.com/rs/cors"
	"net/http"
)

func NewRouter(controller controllers.TraceController) http.Handler {
	router := httprouter.New()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// add routes
	router.GET("/api/find/:id", controller.FindTraceResults)
	router.GET("/api/screenshot", controller.Screenshot)
	router.GET("/api/trace", controller.TraceUrl)

	// Serve static files from the ./assets directory
	// http(s)://api.redirective.net/screenshots/{filename.png}
	router.NotFound = http.FileServer(http.Dir("assets/"))

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.Default().Handler(router)

	// Insert the middleware
	handler = c.Handler(handler)

	return handler
}
