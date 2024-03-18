package service

import (
	"context"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/repo"
	"vk-film-library/internal/repo/repoerrs"
)

type FilmService struct {
	repo repo.FilmRepo
}

func NewFilmService(repo repo.FilmRepo) *FilmService {
	return &FilmService{
		repo: repo,
	}
}

func (f *FilmService) CreateFilm(ctx context.Context, input *entity.FilmCreateInput) (int, error) {
	err := input.Validate()
	if err != nil {
		return 0, err
	}

	film := &entity.Film{
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   input.CreatedAt,
		Rating:      input.Rating,
		Actors:      input.Actors,
	}

	return f.repo.CreateFilm(ctx, film)
}

func (f *FilmService) GetSortFilms(ctx context.Context, sort string) ([]*entity.Film, error) {
	return f.repo.GetSortFilms(ctx, sort)
}

func (f *FilmService) GetFilmsByName(ctx context.Context, namePart string) ([]*entity.Film, error) {
	return f.repo.GetFilmsByName(ctx, namePart)
}

func (f *FilmService) GetFilmsByActor(ctx context.Context, namePart string) ([]*entity.Film, error) {
	return f.repo.GetFilmsByActor(ctx, namePart)
}

func (f *FilmService) DeleteFilm(ctx context.Context, id int) error {
	err := f.repo.DeleteFilm(ctx, id)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return ErrFilmNotFound
		}
		return err
	}

	return nil
}
