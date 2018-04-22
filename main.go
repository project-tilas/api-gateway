package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	health struct {
		ServiceName string `json:"serviceName"`
		Alive       bool   `json:"alive"`
		Version     string `json:"version"`
	}
)

var version string
var addr string

func init() {
	fmt.Println("Running API_GATEWAY version: " + version)
	addr = os.Getenv("API_GATEWAY_ADDR")
	if addr == "" {
		addr = ":8080"
	}

}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/health", func(c echo.Context) error {
		u := health{Alive: true, ServiceName: "api-gateway", Version: version}
		return c.JSON(http.StatusOK, u)
	})

	// Start server
	e.Logger.Fatal(e.Start(addr))
}
