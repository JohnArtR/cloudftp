package db

import "time"

type User struct {
	Username	string
	Password	string
	HomeFolder	string
	StorageType	string
}

type File struct {
	Name string
	Size int
	UploadTime time.Time
	User string
}