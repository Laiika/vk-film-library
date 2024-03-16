package service

import (
	"context"
	"time"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/repo"
)

type Auth interface {
	CreateUser(ctx context.Context, input *entity.CreateInput) (int, error)
	GenerateToken(ctx context.Context, input *entity.AuthInput) (string, error)
	ParseToken(token string) (string, error)
}

type Actor interface {
	CreateActor(ctx context.Context, input *entity.ActorCreateInput) (int, error)
	GetAllActors(ctx context.Context) ([]*entity.Actor, error)
	EditActor(ctx context.Context, actor *entity.Actor) error
	DeleteActor(ctx context.Context, id int) error
}

type Film interface {
	CreateFilm(ctx context.Context, input *entity.FilmCreateInput) (int, error)
	GetFilmsByName(ctx context.Context, namePart string) ([]*entity.Film, error)
	GetFilmsByActor(ctx context.Context, namePart string) ([]*entity.Film, error)
	DeleteFilm(ctx context.Context, id int) error
}

type Services struct {
	Auth  Auth
	Actor Actor
	Film  Film
}

type ServicesDependencies struct {
	Repos *repo.Repositories

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth:  NewAuthService(deps.Repos.UserRepo, deps.SignKey, deps.TokenTTL),
		Actor: NewActorService(deps.Repos.ActorRepo),
		Film:  NewFilmService(deps.Repos.FilmRepo),
	}
}
