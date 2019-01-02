package backend

type User struct {
	ID       string
	Email    string `sql:",unique,notnull"`
	Password string `sql:",notnull"`
}

type Store interface {
	GetUserByID(token string) (*User, error)
}
