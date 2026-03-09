package main

import (
	"log"

	"distributed-cache/service-registry/api"
	serviceRegistry "distributed-cache/service-registry/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := serviceRegistry.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := serviceRegistry.NewServiceRepository(db)
	svc := serviceRegistry.NewServiceRegistryService(repo)
	_ = svc

	r := gin.Default()

	api.RegisterRoutes(r, svc)

	r.Run(":8080")
}
