package main

import (
	"bytes"
	api "distributed-cache/cache/api"
	cache "distributed-cache/cache/internal"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		panic("Empty host env")
	}

	var randPort int

	for {
		randPort = 9000 + rand.Intn(1000)
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", randPort))
		if err == nil {
			ln.Close()
			break
		}
	}

	c := cache.NewCache(60)

	r := gin.Default()

	api.RegisterRoutes(r, c)

	go func() {
		if err := r.Run(fmt.Sprintf(":%d", randPort)); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
	serviceID := registerService(host, randPort)

	fmt.Printf("serviceID: %s", serviceID)
	select {}
}

func registerService(host string, port int) string {
	url := "http://service-registry:8080/services"

	data := map[string]interface{}{
		"name": "cache",
		"host": host,
		"port": port,
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
