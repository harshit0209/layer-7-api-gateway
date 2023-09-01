package main

import (
	"log"
	"os"

	"net/http"

	"github.com/labstack/echo/v4"
)

var Host string

func ForwardRequ(c echo.Context) error {
	Header := c.Request().Header
	Request := c.Request().Body
	URL := c.Request().URL
	Method := c.Request().Method
	log.Print(URL)
	log.Print(Method)
	log.Print(Request)
	log.Println(Header)
	log.Println("Service 1 logged")
	return c.JSON(http.StatusOK, map[string]interface{}{"header": Header, "body": Request, "url": URL, "service": Host})

}

func main() {
	Host = os.Getenv("HTTP_HOST")
	if Host == "" {
		Host = "localhost"
	}
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = ":8080"
	}
	e := echo.New()

	// Define a route that will act as the gateway
	e.Any("*", ForwardRequ)
	// Start the Echo server
	err := e.Start(port) // Replace with your desired port
	log.Print(err)
}
