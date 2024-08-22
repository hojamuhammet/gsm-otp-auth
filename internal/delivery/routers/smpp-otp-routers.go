package routers

import (
	"gsm-otp-auth/internal/delivery/handlers"
	"gsm-otp-auth/internal/service"
	"gsm-otp-auth/pkg/lib/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupOTPRoutes(otpService service.OTPService, logger *logger.Loggers) http.Handler {
	otpRouter := chi.NewRouter()
	otpHandler := handlers.NewOTPHandler(otpService)

	otpRouter.Post("/sendOTP", otpHandler.GenerateAndSaveOTPHandler)
	otpRouter.Post("/validateOTP", otpHandler.ValidateOTPHandler)

	return otpRouter
}
