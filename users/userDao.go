package users

type UserDao interface {
	GetUsers() ([]User, error)
	AddUser(user User) error
	GetUser(userId int) (User, error)
	GetUserFromName(userName string) (User, error)
	//Login()
}
