package main

import (
	"crypto/sha256"
	"testing"
	"log"
	"io"
	"github.com/JohnArtR/cloudftp/server"
	"github.com/JohnArtR/cloudftp/db"
	"fmt"
)

func TestUserRepositorySave(t *testing.T)  {
	server.DBSetup()
	enc := sha256.New()
	io.WriteString(enc, "Gyhyfgz1")
	pass := fmt.Sprintf("%x", enc.Sum(nil))
	user := db.User{Username:"JohnAR", Password:pass, HomeFolder:"/home/johnar", StorageType:"LOCAL"}
	err := db.UserRepo.SaveUser(user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepositoryFind(t *testing.T)  {
	server.DBSetup()
	user, err := db.UserRepo.FindUser("JohnAR")
	if err != nil {
		t.Fatal(err)
	}
	log.Print(user)
}

func TestFileRepositorySave(t *testing.T)  {
	server.DBSetup()
	file1 := db.File{Name:"me.jpg", Size:10240, User:"JohnAR"}
	file2 := db.File{Name:"alen.jpg", Size:14336, User:"JohnAR"}
	err := db.FileRepo.SaveFiles([]db.File{file1, file2})
	if err != nil {
		t.Fatal(err)
	}
}


func TestFileRepositoryFind(t *testing.T)  {
	server.DBSetup()
	files, err := db.FileRepo.FindFiles("JohnAR")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		log.Print(file)
	}
}

func TestEncryption(t *testing.T)  {
	sha := sha256.New()
	io.WriteString(sha, "Na5da41$$;")
	log.Printf("%x", sha.Sum(nil))
}
