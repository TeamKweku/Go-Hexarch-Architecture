package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("healthcheck failed: %s", err)
	}
}

func run() error {
	port, ok := os.LookupEnv("CODE_ODESSEY_PORT")
	if !ok {
		return errors.New("CODE_ODESSEY_PORT is not set")
	}

	client := http.Client{Timeout: 1 * time.Second}
	healthCheckURL := fmt.Sprintf("http://localhost:%s/healthcheck", port)
	res, err := client.Get(healthCheckURL) //nolint:noctx
	if err != nil {
		return fmt.Errorf("healthcheck client error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("app server returned status %d", res.StatusCode)
	}

	return nil
}
