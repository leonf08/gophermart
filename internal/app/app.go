package app

import (
	"github.com/leonf08/gophermart.git/internal/config"
	"github.com/leonf08/gophermart.git/internal/controller/http"
	"github.com/leonf08/gophermart.git/internal/controller/http/handlers"
	"github.com/leonf08/gophermart.git/internal/database/postgres"
	"github.com/leonf08/gophermart.git/internal/logger"
	"github.com/leonf08/gophermart.git/internal/services"
	"github.com/leonf08/gophermart.git/internal/services/repo"
	"os"
	"os/signal"
	"syscall"
)

// Run runs the application.
func Run(cfg *config.Config) {
	log := logger.NewLogger()

	db, err := postgres.NewConnection(cfg.DatabaseAddress)
	if err != nil {
		log.Error("app - Run - postgres.NewConnection", "error", err)
		return
	}
	defer db.Close()

	repository := repo.NewRepository(db)
	auth := services.NewAuthenticator(cfg.JWTSecret)
	userService := services.NewUserManager(repository, auth)
	orderService := services.NewOrderManager(repository, services.NewAccrual(cfg.AccrualAddress, repository, log))

	r := handlers.NewRouter(userService, orderService, auth, log)

	server := http.NewServer(r, cfg.ServerAddress)
	log.Info("app - Run - server.ListenAndServe", "address", cfg.ServerAddress)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-server.Err():
		log.Error("app - Run - server.Err", "error", err)
	case sig := <-interrupt:
		log.Info("app - Run - interrupt", "signal", sig.String())
	}

	log.Info("app - Run - shutdown")
	err = server.Shutdown()
	if err != nil {
		log.Error("app - Run - server.Shutdown", "error", err)
	}
}
