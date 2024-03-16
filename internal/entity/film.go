package entity

import (
	"fmt"
	"time"
)

type Film struct {
	Id          int       `db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Rating      int       `json:"rating" db:"rating"`
	Actors      []string
}

type FilmCreateInput struct {
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Rating      int       `json:"rating" db:"rating"`
	Actors      []string  `json:"actors"`
}

func (form *FilmCreateInput) Validate() error {
	if len(form.Name) < 1 || len(form.Name) > 150 {
		return fmt.Errorf("film name is invalid")
	}
	if len(form.Description) > 1000 {
		return fmt.Errorf("film description is invalid")
	}
	if form.Rating < 0 || form.Rating > 10 {
		return fmt.Errorf("film rating is invalid")
	}

	return nil
}
