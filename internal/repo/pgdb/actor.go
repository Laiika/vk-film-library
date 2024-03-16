package pgdb

import (
	"context"
	"fmt"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/repo/repoerrs"
	"vk-film-library/pkg/postgres"
)

type ActorRepo struct {
	client postgres.Client
}

func NewActorRepo(client postgres.Client) *ActorRepo {
	return &ActorRepo{
		client: client,
	}
}

func (r *ActorRepo) CreateActor(ctx context.Context, actor *entity.Actor) (int, error) {
	query := `INSERT INTO actors (name, gender, birthday) VALUES ($1, $2, $3) RETURNING id`
	var id int

	err := r.client.QueryRow(ctx, query, actor.Name, actor.Gender, actor.Birthday).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ActorRepo CreateActor: %v", err)
	}

	return id, nil
}

func (r *ActorRepo) GetAllActors(ctx context.Context) ([]*entity.Actor, error) {
	query := `SELECT id, name, gender, birthday FROM actors`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ActorRepo GetAllActors: %v", err)
	}

	actors := make([]*entity.Actor, 0)
	for rows.Next() {
		var ac entity.Actor

		err = rows.Scan(&ac.Id, &ac.Name, &ac.Gender, &ac.Birthday)
		if err != nil {
			return nil, fmt.Errorf("ActorRepo GetAllActors: %v", err)
		}

		actors = append(actors, &ac)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ActorRepo GetAllActors: %v", err)
	}

	for _, ac := range actors {
		query = `SELECT name FROM films f JOIN films_actors fa ON fa.film_id = f.id WHERE fa.actor_id=$1`
		rows, err = r.client.Query(ctx, query, ac.Id)
		if err != nil {
			return nil, fmt.Errorf("ActorRepo GetAllActors: %v", err)
		}

		films := make([]string, 0)
		for rows.Next() {
			var f string

			err = rows.Scan(&f)
			if err != nil {
				return nil, fmt.Errorf("ActorRepo GetAllActors: %v", err)
			}

			films = append(films, f)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("ActorRepo GetAllActors: %v", err)
		}

		ac.Films = films
	}

	return actors, nil
}

func (r *ActorRepo) EditActor(ctx context.Context, actor *entity.Actor) error {
	fields := ""
	if actor.Name != "" {
		fields += fmt.Sprintf("name = %s", actor.Name)
	}
	if actor.Gender != "" {
		fields += fmt.Sprintf("gender = %s", actor.Gender)
	}
	query := fmt.Sprintf(`UPDATE actors SET %s WHERE id = $1`, fields)

	commandTag, err := r.client.Exec(ctx, query, actor.Id)
	if err != nil {
		return fmt.Errorf("ActorRepo EditActor: %v", err)
	}
	if commandTag.RowsAffected() != 1 {
		return repoerrs.ErrNotFound
	}

	return nil
}

func (r *ActorRepo) DeleteActor(ctx context.Context, id int) error {
	query := `DELETE FROM actors WHERE id = $1`

	commandTag, err := r.client.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ActorRepo DeleteActor: %v", err)
	}
	if commandTag.RowsAffected() != 1 {
		return repoerrs.ErrNotFound
	}

	return nil
}
