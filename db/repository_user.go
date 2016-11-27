package db

type UserFinder interface {
	FindUser(userName string) (User, error)
}

type UserSaver interface {
	SaveUser(user User) error
}

type UserRepository interface {
	FindUser(userName string) (*User, error)
	SaveUser(user User) error
}
