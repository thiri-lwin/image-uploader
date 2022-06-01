package utilities

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rogpeppe/fastuuid"
	"go_skill_test/models"
	"gopkg.in/mgo.v2"
	"os"
)

var DBSession *mgo.Session
var AuthKey string
var UUIDGenerator *fastuuid.Generator
var DBConfig models.DBConfig

func GetConfig() (models.Config, error) {
	var config models.Config

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file", err)
		return config, err
	}

	AuthKey = os.Getenv("AUTHORIZATION_KEY")

	config.Port = os.Getenv("PORT")
	config.DBConfig = models.DBConfig{
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		UserName:     os.Getenv("USER_NAME"),
		Password:     os.Getenv("PASSWORD"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
	}

	return config, nil
}

func GetConnection(dbConfig models.DBConfig) (*mgo.Session, error) {
	var url string
	if dbConfig.UserName == "" {
		url = dbConfig.Host + ":" + dbConfig.Port
	} else {
		url = "mongodb://" + dbConfig.UserName + ":" + dbConfig.Password + "@" + dbConfig.Host + ":" + dbConfig.Port
	}
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println("Could not connect to database")
	} else {
		fmt.Println("Database connected!")
	}
	return session, err
}

func CreateGenerator() error {
	generator, err := fastuuid.NewGenerator()
	if err != nil {
		fmt.Println("could not create *generator")
	}
	UUIDGenerator = generator
	return err
}

func CreateNewUUID() string {
	return fastuuid.Hex128(UUIDGenerator.Next())
}
