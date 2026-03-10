package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/spaolacci/murmur3"
)

const serviceRegistryURL = "http://service-registry:8080"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomValue(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func SetValue(value string, registryURL string) (string, error) {
	hash := murmur3.Sum32([]byte(value))

	registryURL = fmt.Sprintf("%s/services/instances/?hash_key=%d", registryURL, hash)
	resp, err := http.Get(registryURL)
	if err != nil {
		return "", fmt.Errorf("service registry request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("service registry returned status %d", resp.StatusCode)
	}

	var instance struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return "", fmt.Errorf("failed to decode service registry response: %w", err)
	}

	payload, _ := json.Marshal(map[string]string{"value": value})
	cacheURL := fmt.Sprintf("http://%s:%d/cache/", instance.Host, instance.Port)
	cacheResp, err := http.Post(cacheURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("cache request failed: %w", err)
	}
	defer cacheResp.Body.Close()

	if cacheResp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("cache returned status %d", cacheResp.StatusCode)
	}

	var result struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(cacheResp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode cache response: %w", err)
	}

	return result.Key, nil
}

func main() {
	for {
		value := randomValue(16)
		key, err := SetValue(value, serviceRegistryURL)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		} else {
			fmt.Printf("stored value=%q at key=%s\n", value, key)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
