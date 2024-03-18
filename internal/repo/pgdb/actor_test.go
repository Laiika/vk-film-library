package pgdb

import (
	"context"
	"errors"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"vk-film-library/internal/entity"
)

func TestActorRepo_CreateActor(t *testing.T) {
	type args struct {
		ctx   context.Context
		actor *entity.Actor
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				actor: &entity.Actor{
					Name:     "asjdsk",
					Gender:   "men",
					Birthday: time.UnixMilli(123456),
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(1)

				m.ExpectQuery("INSERT INTO actors").
					WithArgs(args.actor.Name, args.actor.Gender, args.actor.Birthday).
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "unexpected error",
			args: args{
				ctx: context.Background(),
				actor: &entity.Actor{
					Name:     "asjdsk",
					Gender:   "men",
					Birthday: time.UnixMilli(123456),
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.actor.Name, args.actor.Gender, args.actor.Birthday).
					WillReturnError(errors.New("some error"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := poolMock
			userRepoMock := NewActorRepo(postgresMock)

			got, err := userRepoMock.CreateActor(tc.args.ctx, tc.args.actor)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
