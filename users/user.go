package users

type User struct {
	Id       int
	Name     string
	Password string
}

func NewUser(id int, name string, password string) User {
	return User{
		Id:       id,
		Name:     name,
		Password: password,
	}
}
