package main

import (
	"testing"
	"log"
	"eventview.online/bigftp/server"
	"eventview.online/bigftp/db"
)

func TestUserRepositorySave(t *testing.T)  {
	server.DBSetup()
	user := db.User{Username:"JohnAR", Pssword:"Gyhyfgz1", HomeFolder:"/home/johnar", StorageType:"LOCAL"}
	err := server.UserRepository.SaveUser(user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepositoryFind(t *testing.T)  {
	server.DBSetup()
	user, err := server.UserRepository.FindUser("JohnAR")
	if err != nil {
		t.Fatal(err)
	}
	log.Print(user)
}

func TestFileRepositorySave(t *testing.T)  {
	server.DBSetup()
	file1 := db.File{Name:"me.jpg", Size:10240, User:"JohnAR"}
	file2 := db.File{Name:"alen.jpg", Size:14336, User:"JohnAR"}
	err := server.FileRepository.SaveFiles([]db.File{file1, file2})
	if err != nil {
		t.Fatal(err)
	}
}


func TestFileRepositoryFind(t *testing.T)  {
	server.DBSetup()
	files, err := server.FileRepository.FindFiles("JohnAR")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		log.Print(file)
	}
}
