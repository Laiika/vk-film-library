package entity

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

type AuthInput struct {
	Username string
	Password string
}

type CreateInput struct {
	Username string
	Password string
	Role     string
}
