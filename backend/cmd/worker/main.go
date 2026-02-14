package main

import (
	"log"
	"time"

	"chronovault/internal/database"
	"chronovault/internal/repository"
	"chronovault/internal/services"
	"chronovault/internal/websocket"
)

func main() {
	log.Println("Starting ChronoVault Background Worker...")

	db, err := database.Connect("./chronovault.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.New(db)
	wsHub := websocket.NewHub()
	go wsHub.Run()

	obligationService := services.NewObligationService(repo, wsHub)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Worker started. Evaluating obligations every minute...")

	for {
		select {
		case <-ticker.C:
			orgs, err := repo.ListOrganizations()
			if err != nil {
				log.Printf("Error fetching organizations: %v", err)
				continue
			}

			for _, org := range orgs {
				if err := obligationService.EvaluateObligations(org.ID); err != nil {
					log.Printf("Error evaluating obligations for org %s: %v", org.ID, err)
				}
			}
		}
	}
}
