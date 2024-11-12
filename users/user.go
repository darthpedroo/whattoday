package users

type User struct {
	Id   int
	Name string
}

func NewUser(id int, name string) User {
	return User{
		Id:   id,
		Name: name,
	}
}
