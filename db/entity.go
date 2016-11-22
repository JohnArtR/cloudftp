package db

type User struct {
	Username	string
	Pssword		string
	HomeFolder	string
	StorageType	string
}

type File struct {
	Name string
	Size int
	User string
}