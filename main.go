package main

import (
	"database/sql"
	"fmt"
	"log"
	"visa_assignment/config"
	"visa_assignment/controller"
	"visa_assignment/repository"
	"visa_assignment/router"
	"visa_assignment/service"
)

func main() {
	// 1. Load configuration from resources/application.yml
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Connect to MS SQL Server
	db := initDatabase()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// 3. Initialize Repositories
	memberRepo := repository.NewMemberRepository(db)
	rewardRepo := repository.NewRewardRepository(db)

	// 4. Initialize Services
	memberService := service.NewMemberService(memberRepo, rewardRepo)
	rewardService := service.NewRewardService(memberRepo, rewardRepo)

	// 5. Initialize Controllers
	memberController := controller.NewMemberController(memberService, rewardService)
	rewardController := controller.NewRewardController(rewardService)

	// 6. Setup Router
	startRouter(memberController, rewardController)
}

func startRouter(memberController *controller.MemberController, rewardController *controller.RewardController) {
	r := router.SetupRouter(memberController, rewardController)

	port := config.GetConfig().Port
	if port == 0 {
		port = 8080 // Default port if not specified in config
	}

	// 8. Start HTTP Server
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting Service on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initDatabase() *sql.DB {
	db, err := repository.InitDB(config.GetConfig().Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to SQL Server")
	return db
}
