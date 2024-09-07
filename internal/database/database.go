package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	URL string
	Port string
	Mongo *mongo.Client
}

func NewDatabase() *Database {
	url := os.Getenv("mongo_url")
	port := os.Getenv("mongo_port")
	db := &Database{}
	db.Port = port
	db.URL = url
	err := db.StartMongo()
	if err != nil {
		log.Panic("Error setting mongo db")
	}
	return db
}

func (d *Database) StartMongo() error {
	usr := os.Getenv("mongo_user")
	pwd := os.Getenv("mongo_pwd")
	// crear opciones de conexión
	clientOptions := options.Client().ApplyURI(d.URL)
	clientOptions.SetAuth(options.Credential{
		Username: usr,
		Password: pwd,
	})

	// conexión
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Connection Error:", err)
		return err
	}

	err = c.Ping(context.Background(), nil)
	if err != nil {
		log.Println("Test Connection failed wit error:",err.Error())
		return err
	}

	log.Println("Connected to mongo!")
	d.Mongo = c
	return nil	
}

