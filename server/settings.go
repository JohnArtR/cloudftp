package server

import "github.com/naoina/toml"
import "os"
import "io/ioutil"
import (
	"crypto/tls"
	"strings"
	"strconv"
	"log"
	"github.com/JohnArtR/cloudftp/db"
	"github.com/JohnArtR/cloudftp/storage"
)

const (
	ENV_DB_ENGINE	= "DB_ENGINE"
	ENV_DB_HOST	= "DB_HOST"
	ENV_DB_PORT	= "DB_PORT"
	ENV_DB_NAME	= "DB_NAME"

	DB_ENGINE_MONGO		= "mongo"
	DB_ENGINE_POSTGRES	= "postgres"
)

type ParadiseSettings struct {
	Host           string
	Port           int
	MaxConnections int
	MaxPassive     int
	Exec           string
	Pem            string
	Key            string
	StorageDirectory string
}

func Load509Config() *tls.Config {
	// use https://letsencrypt.org to get the pem and key files
	cert, cerr := tls.LoadX509KeyPair(Settings.Pem, Settings.Key)
	if cerr != nil {
		return nil
	}

	config := &tls.Config{}
	if config.NextProtos == nil {
		config.NextProtos = []string{"ftp"}
	}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = cert

	return config
}

func ReadSettings() ParadiseSettings {
	f, err := os.Open("conf/settings.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config ParadiseSettings
	if err := toml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}
	FileService = storage.FileManager(storage.NewGCloudFileService())
	return config
}

func DBSetup()  {
	var dbEngine string
	if dbEngine = os.Getenv(ENV_DB_ENGINE); strings.EqualFold(dbEngine, "") {
		log.Fatalf(" [ERROR] Env variable $%s not set", ENV_DB_ENGINE)
	}
	FileService = storage.FileManager(storage.NewGCloudFileService())
	dbHost, dbPort, dbName := getDBParams()
	switch dbEngine {
	case DB_ENGINE_MONGO:
		db.UserRepo = db.UserRepository(db.MongoEngine{Host: dbHost, Port:dbPort, DBName:dbName})
		db.FileRepo = db.FileRepository(db.MongoEngine{Host: dbHost, Port:dbPort, DBName:dbName})
	case DB_ENGINE_POSTGRES:
		log.Println(" [INFO] Postgres engine not supported now. CloudFTP will be exit.")
	default:
		log.Fatal(" [ERROR] Wrong DB_ENGINE setting")
	}
}

func getDBParams() (dbHost, dbPort, dbName string) {
	if dbHost = os.Getenv(ENV_DB_HOST); strings.EqualFold(dbHost, "") {
		log.Fatalf(" [ERROR] Env variable $%s not set", ENV_DB_HOST)
	}

	if dbPort = os.Getenv(ENV_DB_PORT); strings.EqualFold(dbPort, "") {
		log.Fatalf(" [ERROR] Env variable $%s not set", ENV_DB_PORT)
	}

	_, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatal("[ ERROR] DB_PORT must be a number")
	}

	if dbName = os.Getenv(ENV_DB_NAME); strings.EqualFold(dbName, "") {
		log.Fatalf(" [ERROR] Env variable $%s not set", ENV_DB_NAME)
	}

	return dbHost, dbPort, dbName
}