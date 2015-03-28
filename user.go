package townsita

type User struct {
	ID string
}

func NewUser() *User {
	return &User{}
}
