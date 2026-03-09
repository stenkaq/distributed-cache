package main

import (
	"bytes"
	api "distributed-cache/cache/api"
	cache "distributed-cache/cache/internal"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		panic("Empty host env")
	}

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	randPort := seed.Intn(9999-9000) + 9000

	c := cache.NewCache(60)

	r := gin.Default()

	api.RegisterRoutes(r, c)

	serviceID := registerService()
	registerInstance(serviceID, host, randPort)

	r.Run(fmt.Sprintf(":%d", randPort))
}

func registerService() string {
	url := "http://service-registry:8080/services"

	data := map[string]string{
		"name": "cache",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	type ServiceResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var svc ServiceResponse
	if err := json.Unmarshal(body, &svc); err != nil {
		panic(err)
	}

	return svc.ID
}

func registerInstance(serviceID string, host string, port int) {
	url := "http://service-registry:8080/instances"

	data := map[string]interface{}{
		"ID":     serviceID,
		"host":   host,
		"port":   port,
		"status": "UP",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
