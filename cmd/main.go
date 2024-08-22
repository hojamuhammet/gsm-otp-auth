package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gsm-otp-auth/internal/config"
	"gsm-otp-auth/internal/delivery/routers"
	"gsm-otp-auth/internal/repository"
	"gsm-otp-auth/internal/service"
	db "gsm-otp-auth/pkg/database"
	"gsm-otp-auth/pkg/lib/logger"
	"gsm-otp-auth/pkg/lib/utils"
)

func main() {
	cfg := config.LoadConfig()

	logger, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("failed to set up logger: %v", err)
	}

	database, err := db.InitDB(cfg)
	if err != nil {
		logger.ErrorLogger.Error("failed to initialize database: %v", utils.Err(err))
		os.Exit(1)
	}

	repo := repository.NewOTPRepository(database.GetClient(), logger)

	otpService := service.NewOTPService(repo, logger, cfg)

	r := routers.SetupOTPRoutes(otpService, logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(cfg.HTTPServer.Address, r); err != nil && err != http.ErrServerClosed {
			logger.ErrorLogger.Error("Server failed to start:", utils.Err(err))
			os.Exit(1)
		}
	}()

	logger.InfoLogger.Info("Server is up and running")

	<-stop
	logger.InfoLogger.Info("Shutting down the server gracefully...")
	if err := database.Close(); err != nil {
		logger.ErrorLogger.Error("Error closing database:", utils.Err(err))
	}
	os.Exit(0)
}
