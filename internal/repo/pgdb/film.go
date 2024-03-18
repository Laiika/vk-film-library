package pgdb

import (
	"context"
	"fmt"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/repo/repoerrs"
	"vk-film-library/pkg/postgres"
)

type FilmRepo struct {
	client postgres.Client
}

func NewFilmRepo(client postgres.Client) *FilmRepo {
	return &FilmRepo{
		client: client,
	}
}

func (r *FilmRepo) CreateFilm(ctx context.Context, film *entity.Film) (int, error) {
	query := `INSERT INTO films (name, description, created_at, rating) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int

	err := r.client.QueryRow(ctx, query, film.Name, film.Description, film.CreatedAt, film.Rating).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("FilmRepo CreateFilm: %v", err)
	}

	for _, actor := range film.Actors {
		id2, err := r.getActorIdByName(ctx, actor)
		if err != nil {
			return 0, err
		}

		query = `INSERT INTO films_actors (film_id, actor_id) VALUES ($1, $2)`
		_, err = r.client.Exec(ctx, query, id, id2)
		if err != nil {
			return 0, fmt.Errorf("FilmRepo CreateFilm: %v", err)
		}
	}

	return id, nil
}

func (r *FilmRepo) getActorIdByName(ctx context.Context, name string) (int, error) {
	query := `SELECT id FROM actors WHERE name = $1`
	var id int

	err := r.client.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("FilmRepo CreateFilm: %v", err)
	}

	return id, nil
}

func (r *FilmRepo) GetSortFilms(ctx context.Context, sort string) ([]*entity.Film, error) {
	var query string
	switch sort {
	case "rating":
		query = `SELECT id, name, description, created_at, rating FROM films ORDER BY rating`
	case "name":
		query = `SELECT id, name, description, created_at, rating FROM films ORDER BY name`
	case "created_at":
		query = `SELECT id, name, description, created_at, rating FROM films ORDER BY created_at`
	default:
		return nil, fmt.Errorf("FilmRepo GetFilmsByName: invalid sort field")
	}

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
	}

	films := make([]*entity.Film, 0)
	for rows.Next() {
		var f entity.Film

		err = rows.Scan(&f.Id, &f.Name, &f.Description, &f.CreatedAt, &f.Rating)
		if err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
		}

		films = append(films, &f)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
	}

	for _, f := range films {
		query = `SELECT name FROM actors ac JOIN films_actors fa ON fa.actor_id = ac.id WHERE fa.film_id=$1`
		rows, err = r.client.Query(ctx, query, f.Id)
		if err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
		}

		actors := make([]string, 0)
		for rows.Next() {
			var ac string

			err = rows.Scan(&ac)
			if err != nil {
				return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
			}

			actors = append(actors, ac)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
		}

		f.Actors = actors
	}

	return films, nil
}

func (r *FilmRepo) GetFilmsByName(ctx context.Context, namePart string) ([]*entity.Film, error) {
	query := `SELECT id, name, description, created_at, rating FROM films WHERE name LIKE '%` + namePart + `%'`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
	}

	films := make([]*entity.Film, 0)
	for rows.Next() {
		var f entity.Film

		err = rows.Scan(&f.Id, &f.Name, &f.Description, &f.CreatedAt, &f.Rating)
		if err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
		}

		films = append(films, &f)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
	}

	for _, f := range films {
		query = `SELECT name FROM actors ac JOIN films_actors fa ON fa.actor_id = ac.id WHERE fa.film_id=$1`
		rows, err = r.client.Query(ctx, query, f.Id)
		if err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
		}

		actors := make([]string, 0)
		for rows.Next() {
			var ac string

			err = rows.Scan(&ac)
			if err != nil {
				return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
			}

			actors = append(actors, ac)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByName: %v", err)
		}

		f.Actors = actors
	}

	return films, nil
}

func (r *FilmRepo) GetFilmsByActor(ctx context.Context, namePart string) ([]*entity.Film, error) {
	query := `SELECT fa.film_id FROM films_actors fa JOIN actors ac ON ac.id = fa.actor_id WHERE ac.name LIKE '%` +
		namePart + `%'`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FilmRepo GetFilmsByActor: %v", err)
	}

	films := make([]*entity.Film, 0)
	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByActor: %v", err)
		}

		f, err := r.getFilmById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByActor: %v", err)
		}
		f.Id = id
		films = append(films, f)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("FilmRepo GetFilmsByActor: %v", err)
	}

	for _, f := range films {
		query = `SELECT name FROM actors ac JOIN films_actors fa ON fa.actor_id = ac.id WHERE fa.film_id=$1`
		rows, err = r.client.Query(ctx, query, f.Id)
		if err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByActors: %v", err)
		}

		actors := make([]string, 0)
		for rows.Next() {
			var ac string

			err = rows.Scan(&ac)
			if err != nil {
				return nil, fmt.Errorf("FilmRepo GetFilmsByActors: %v", err)
			}

			actors = append(actors, ac)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("FilmRepo GetFilmsByActors: %v", err)
		}

		f.Actors = actors
	}

	return films, nil
}

func (r *FilmRepo) getFilmById(ctx context.Context, id int) (*entity.Film, error) {
	query := `SELECT name, description, created_at, rating FROM films WHERE id = $1`
	var film entity.Film

	err := r.client.QueryRow(ctx, query, id).Scan(&film.Name, &film.Description, &film.CreatedAt, &film.Rating)
	if err != nil {
		return nil, fmt.Errorf("FilmRepo GetFilmsByActor: %v", err)
	}

	return &film, nil
}

func (r *FilmRepo) DeleteFilm(ctx context.Context, id int) error {
	query := `DELETE FROM films WHERE id = $1`

	commandTag, err := r.client.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("FilmRepo DeleteFilm: %v", err)
	}
	if commandTag.RowsAffected() != 1 {
		return repoerrs.ErrNotFound
	}

	return nil
}
