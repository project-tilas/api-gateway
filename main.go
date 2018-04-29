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
		Services    []health `json:"services,omitempty"`
		PodName     string   `json:"podName,omitempty"`
		NodeName    string   `json:"nodeName,omitempty"`
		Hits        int      `json:"hits,omitempty"`
	}
)

var version string
var addr string
var podName string
var nodeName string

func init() {
	fmt.Println("Running API_GATEWAY version: " + version)
	addr = getEnvVar("API_GATEWAY_ADDR", ":8080")
	nodeName = getEnvVar("API_GATEWAY_NODE_NAME", "N/A")
	podName = getEnvVar("API_GATEWAY_POD_NAME", "N/A")
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

		svcAuthHealth, svcAuthErr := testService("svc-auth")
		if svcAuthErr != nil {
			log.Error(svcAuthErr)
		} else {
			services = append(services, svcAuthHealth)
		}

		u := health{
			Alive:       true,
			ServiceName: "api-gateway",
			Version:     version,
			Services:    services,
			PodName:     podName,
			NodeName:    nodeName,
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

func getEnvVar(env string, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
