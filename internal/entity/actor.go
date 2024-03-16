package entity

import "time"

type Actor struct {
	Id       int       `db:"id"`
	Name     string    `json:"name" db:"name"`
	Gender   string    `json:"gender" db:"gender"`
	Birthday time.Time `json:"birthday" db:"birthday"`
	Films    []string
}

type ActorCreateInput struct {
	Name     string    `json:"name"`
	Gender   string    `json:"gender"`
	Birthday time.Time `json:"birthday"`
}
