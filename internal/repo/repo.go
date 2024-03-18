package repo

import (
	"context"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/repo/pgdb"
	"vk-film-library/pkg/postgres"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) (int, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*entity.User, error)
}

type ActorRepo interface {
	CreateActor(ctx context.Context, actor *entity.Actor) (int, error)
	GetAllActors(ctx context.Context) ([]*entity.Actor, error)
	EditActor(ctx context.Context, actor *entity.Actor) error
	DeleteActor(ctx context.Context, id int) error
}

type FilmRepo interface {
	CreateFilm(ctx context.Context, film *entity.Film) (int, error)
	GetSortFilms(ctx context.Context, sort string) ([]*entity.Film, error)
	GetFilmsByName(ctx context.Context, namePart string) ([]*entity.Film, error)
	GetFilmsByActor(ctx context.Context, namePart string) ([]*entity.Film, error)
	//EditFilmByName(ctx context.Context, film *entity.Film) error
	DeleteFilm(ctx context.Context, id int) error
}

type Repositories struct {
	UserRepo
	ActorRepo
	FilmRepo
}

func NewRepositories(client postgres.Client) *Repositories {
	return &Repositories{
		UserRepo:  pgdb.NewUserRepo(client),
		ActorRepo: pgdb.NewActorRepo(client),
		FilmRepo:  pgdb.NewFilmRepo(client),
	}
}
