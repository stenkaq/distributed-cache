package main

import (
	"log"

	serviceRegistry "cache/service-registry/internal"
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
}
