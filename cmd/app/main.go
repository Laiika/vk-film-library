package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vk-film-library/config"
	v1 "vk-film-library/internal/controller/http/v1"
	"vk-film-library/internal/httpserver"
	"vk-film-library/internal/repo"
	"vk-film-library/internal/service"
	"vk-film-library/pkg/logger"
	"vk-film-library/pkg/postgres"
)

// @title           Film Service
// @version         1.0
// @description     This is a service for viewing information about films and actors.

// @contact.name   Zavoiskih Evgenia
// @contact.email  evgeniazavojskih@gmail.com

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  JWT
// @in                          header
// @name                        Authorization
// @description					JWT token
func main() {
	log := logger.GetLogger()
	cfg := config.GetConfig(log)

	log.Info("initializing postgres")
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)
	client, err := postgres.NewClient(context.Background(), 3, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	log.Info("initializing repositories")
	repos := repo.NewRepositories(client)

	log.Info("initializing services")
	deps := service.ServicesDependencies{
		Repos:    repos,
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.NewServices(deps)

	mux := http.NewServeMux()
	v1.NewRouter(mux, services, log)
	log.Info("starting http server")
	log.Debug("server port: %s", cfg.HTTPServer.Port)
	address := fmt.Sprintf("%s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	httpServer := httpserver.New(mux, address)
	go func() {
		err = httpServer.Start()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = httpServer.Stop()
	if err != nil {
		panic(errors.Wrap(err, "Httpserver shutdown"))
	}
	log.Debug("Httpserver exited")
}
