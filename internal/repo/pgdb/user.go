package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/repo/repoerrs"
	"vk-film-library/pkg/postgres"
)

type UserRepo struct {
	client postgres.Client
}

func NewUserRepo(client postgres.Client) *UserRepo {
	return &UserRepo{
		client: client,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	query := `INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id`
	var id int

	err := r.client.QueryRow(ctx, query, user.Username, user.Password, user.Role).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrs.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("UserRepo CreateUser: %v", err)
	}

	return id, nil
}

func (r *UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, username, password, role FROM users WHERE username=$1 AND password=$2`

	err := r.client.QueryRow(ctx, query, username, password).Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &entity.User{}, repoerrs.ErrNotFound
		}
		return &entity.User{}, fmt.Errorf("UserRepo GetUserByUsernameAndPassword: %v", err)
	}

	return &user, nil
}
