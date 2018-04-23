package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	health struct {
		ServiceName string   `json:"serviceName"`
		Alive       bool     `json:"alive"`
		Version     string   `json:"version"`
		Services    []health `json:"services"`
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
		services := []health{}
		svcTilsHealth, svcTilsErr := testService("svc-tils")
		if svcTilsErr != nil {
			log.Error(svcTilsErr)
		} else {
			services = append(services, svcTilsHealth)
		}
		u := health{
			Alive:       true,
			ServiceName: "api-gateway",
			Version:     version,
			Services:    services,
		}
		return c.JSON(http.StatusOK, u)
	})

	// Start server
	e.Logger.Fatal(e.Start(addr))
}

func testService(serviceName string) (health, error) {
	url := "http://" + serviceName + "/health"

	httpClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return health{}, err
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return health{}, getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return health{}, readErr
	}

	healthResp := health{}
	jsonErr := json.Unmarshal(body, &healthResp)
	if jsonErr != nil {
		return health{}, jsonErr
	}
	return healthResp, nil
}
