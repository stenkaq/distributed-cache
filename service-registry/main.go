package main

import (
	"log"
	"time"

	"distributed-cache/service-registry/api"
	serviceRegistry "distributed-cache/service-registry/internal"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/otter/v2"
	"github.com/maypok86/otter/v2/stats"
)

func main() {
	db, err := serviceRegistry.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	counter := stats.NewCounter()

	cache := otter.Must(&otter.Options[int, *serviceRegistry.ServiceInstance]{
		MaximumSize:       10_000,
		RefreshCalculator: otter.RefreshWriting[int, *serviceRegistry.ServiceInstance](60 * time.Second),
		StatsRecorder:     counter,
	})

	repo := serviceRegistry.NewServiceRepository(db)
	ringRepo := serviceRegistry.NewRingRepository(*cache)
	svc := serviceRegistry.NewServiceRegistryService(repo, ringRepo)
	_ = svc

	r := gin.Default()

	api.RegisterRoutes(r, svc)

	r.Run(":8080")
}
