package v1

import (
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"vk-film-library/internal/service"
	"vk-film-library/pkg/logger"
)

func NewRouter(mux *http.ServeMux, services *service.Services, log *logger.Logger) {
	mux.HandleFunc("/swagger/", httpSwagger.Handler())

	newAuthRoutes(mux, services.Auth, log)
	authMiddleware := &AuthMiddleware{services.Auth, log}

	newActorRoutes(mux, services.Actor, authMiddleware, log)
	newFilmRoutes(mux, services.Film, authMiddleware, log)
}
