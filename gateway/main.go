package main

import (
	"log"

	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

type CustomHandlerFunc func(http.ResponseWriter, *http.Request) error

func (h CustomHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		log.Println("Error:", err)
	}
}
func main() {
	e := echo.New()
	// Define target URLs for multiple proxies
	targetURLs := map[string]string{
		"proxy1*": "http://http-service1:8082",
		"proxy2":  "http://http-service2:8083",
	}
	// Create reverse proxies for each target URL

	proxies := make(map[string]*httputil.ReverseProxy)
	for name, targetURL := range targetURLs {
		target, err := url.Parse(targetURL)
		if err != nil {
			log.Fatalf("Invalid target URL for %s: %s", name, err)
		}
		proxies[name] = httputil.NewSingleHostReverseProxy(target)
	}
	// Handler function to handle proxying requests
	proxyHandler := func(proxy *httputil.ReverseProxy) CustomHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			log.Println("Proxying request:", r.URL.Path)
			// Remove the "/proxy1" prefix from the request path
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/proxy1")
			// Set the proper target URL for the reverse proxy
			target := proxy.Director
			target(r)
			// Serve the request using the reverse proxy
			proxy.ServeHTTP(w, r)
			return nil
		}
	}
	// Register reverse proxy handlers for each target URL
	for name, proxy := range proxies {
		handler := proxyHandler(proxy)
		path := "/" + name + "/*"
		e.Any(path, echo.WrapHandler(handler))
	}
	// Log all requests
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Println("Received request:", c.Request().URL.Path)
			return next(c)
		}
	})
	e.Use(AddHeaders)
	// Start the server
	e.Start(":8089")

}
func AddHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Your logic here...

		// Set headers
		c.Response().Header().Set("Custom-Header-1", "Header-Value-1")
		c.Response().Header().Set("Custom-Header-2", "Header-Value-2")

		// Call the next middleware or handler in the chain
		return next(c)
	}
}
