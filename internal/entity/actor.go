package entity

type Actor struct {
	Id       int    `db:"id"`
	Name     string `json:"name" db:"name"`
	Gender   string `json:"gender" db:"gender"`
	Birthday string `json:"birthday" db:"birthday"`
	Films    []string
}

type ActorCreateInput struct {
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}
