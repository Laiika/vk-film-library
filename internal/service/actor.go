package service

import (
	"context"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/repo"
	"vk-film-library/internal/repo/repoerrs"
)

type ActorService struct {
	repo repo.ActorRepo
}

func NewActorService(repo repo.ActorRepo) *ActorService {
	return &ActorService{
		repo: repo,
	}
}

func (a *ActorService) CreateActor(ctx context.Context, input *entity.ActorCreateInput) (int, error) {
	actor := &entity.Actor{
		Name:     input.Name,
		Gender:   input.Gender,
		Birthday: input.Birthday,
	}

	return a.repo.CreateActor(ctx, actor)
}

func (a *ActorService) GetAllActors(ctx context.Context) ([]*entity.Actor, error) {
	return a.repo.GetAllActors(ctx)
}

func (a *ActorService) EditActor(ctx context.Context, input *entity.Actor) error {
	err := a.repo.EditActor(ctx, input)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return ErrActorNotFound
		}
		return err
	}

	return nil
}

func (a *ActorService) DeleteActor(ctx context.Context, id int) error {
	err := a.repo.DeleteActor(ctx, id)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return ErrActorNotFound
		}
		return err
	}

	return nil
}
