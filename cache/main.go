package main

import (
	api "distributed-cache/cache/api"
	cache "distributed-cache/cache/internal"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	randPort := seed.Intn(9999-9000) + 9000

	c := cache.NewCache(60)

	r := gin.Default()

	api.RegisterRoutes(r, c)

	r.Run(fmt.Sprintf(":%d", randPort))
}
