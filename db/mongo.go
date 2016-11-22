package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	MONGO_USERS_COLL = "users"
	MONGO_FILES_COLL = "files"
)

type MongoEngine struct {
	Host string
	Port string
	DBName string
	session *mgo.Session
}

func (repo MongoEngine) getDBConnection() *mgo.Database {
	var err error
	repo.session, err = mgo.Dial(repo.Host)
	if err != nil {
		log.Fatal(" [ERROR] Can't connect to mongo. Check settings file")
	}
	return repo.session.DB(repo.DBName)
}

func (repo MongoEngine) FindUser(userName string) (user User, err error) {
	db := repo.getDBConnection()
	collection := db.C(MONGO_USERS_COLL)
	err = collection.Find(bson.M{"username": userName}).One(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (repo MongoEngine) FindFile(fileName string) (file File, err error) {
	db := repo.getDBConnection()
	collection := db.C(MONGO_FILES_COLL)
	err = collection.Find(bson.M{"name": fileName}).One(&file)
	if err != nil {
		return File{}, err
	}
	return file, nil
}

func (repo MongoEngine) FindFiles(userName string) (files []File, err error) {
	db := repo.getDBConnection()
	collection := db.C(MONGO_FILES_COLL)
	err = collection.Find(bson.M{"user": userName}).All(&files)
	if err != nil {
		return []File{}, err
	}
	return files, nil
}

func (repo MongoEngine) SaveUser(user User) error {
	db := repo.getDBConnection()
	repo.session.Close()
	collection := db.C(MONGO_USERS_COLL)
	err := collection.Insert(user)
	return err
}

func (repo MongoEngine) SaveFile(file File) error {
	db := repo.getDBConnection()
	collection := db.C(MONGO_FILES_COLL)
	err := collection.Insert(file)
	return err
}

func (repo MongoEngine) SaveFiles(files []File) error {
	db := repo.getDBConnection()
	collection := db.C(MONGO_FILES_COLL)
	if len(files) == 0 {
		log.Print(" [INFO] No files to save")
		return nil
	}

	for _, file := range files {
		err := collection.Insert(file)
		if err != nil {
			return err
		}
	}
	return nil
}