package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"visa_assignment/config"
	"visa_assignment/controller"
	"visa_assignment/repository"
	"visa_assignment/router"
	"visa_assignment/service"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	initLogger()

	// 1. Load configuration from resources/application.yml
	err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// 2. Connect to MS SQL Server
	db := initDatabase()
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
		}
	}()

	// 3. Initialize Repositories
	memberRepo := repository.NewMemberRepository(db)
	rewardRepo := repository.NewRewardRepository(db)

	// 4. Initialize Services
	memberService := service.NewMemberService(memberRepo, rewardRepo)
	rewardService := service.NewRewardService(memberRepo, rewardRepo, config.GetConfig().PointTypes)

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
	log.Info().Str("address", addr).Msg("Starting service")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

func initDatabase() *sql.DB {
	db, err := repository.InitDB(config.GetConfig().Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	log.Info().Msg("Successfully connected to SQL Server")
	return db
}

func initLogger() {
	// Standardize on RFC3339 / UTC timestamps for ledger auditing
	zerolog.TimeFieldFormat = time.RFC3339

	// Set global log level (InfoLevel for prod, DebugLevel for local debugging)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Configure global logger with timestamp and file:line caller info
	log.Logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}
