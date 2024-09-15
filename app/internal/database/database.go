package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	URL     string
	Port    string
	Mongo   *mongo.Client
	MongoDb string
}

func NewDatabase() *Database {
	url := os.Getenv("mongo_url")
	if url == "" {
		url = "mongodb://localhost:27017"
	}
	port := os.Getenv("mongo_port")
	if port == "" {
		port = "27017"
	}
	mdb := os.Getenv("mongo_db")
	if mdb == "" {
		mdb = "reverse-proxy"
	}
	db := &Database{}
	db.Port = port
	db.URL = url
	db.MongoDb = mdb

	log.Println("Mongo URL:", db.URL)
	log.Println("Mongo Port:", db.Port)

	err := db.StartMongo()
	if err != nil {
		log.Fatal("Error setting MongoDB:", err)
	}

	// Verificar si la base de datos y la colección existen
	err = db.CheckAndCreateDatabase()
	if err != nil {
		log.Fatal("Error checking or creating database:", err)
	}

	return db
}

func (d *Database) StartMongo() error {
	usr := os.Getenv("mongo_user")
	pwd := os.Getenv("mongo_pwd")

	// Crear opciones de conexión
	clientOptions := options.Client().ApplyURI(d.URL)
	clientOptions.SetAuth(options.Credential{
		Username: usr,
		Password: pwd,
	})

	// Conectar
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Connection Error:", err)
		return err
	}

	// Verificar conexión
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Test Connection failed with error:", err)
		return err
	}

	// Listar bases de datos para verificar si la conexión funciona
	databases, err := client.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error listing databases:", err)
		return err
	}
	log.Println("Databases available:", databases)

	// Asignar el cliente de Mongo a la estructura
	d.Mongo = client

	log.Println("Connected to MongoDB!")
	return nil
}

func (d *Database) CheckAndCreateDatabase() error {
	// Listar todas las bases de datos
	dbs, err := d.Mongo.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error listing databases:", err)
		return err
	}

	// Verificar si la base de datos existe
	dbExists := false
	for _, dbName := range dbs {
		if dbName == d.MongoDb {
			dbExists = true
			break
		}
	}

	if !dbExists {
		log.Printf("Database '%s' does not exist, creating...\n", d.MongoDb)

		// Crear una colección dentro de la base de datos para que MongoDB la cree
		collection := d.Mongo.Database(d.MongoDb).Collection("init")
		_, err := collection.InsertOne(context.TODO(), bson.M{"status": "initialized"})
		if err != nil {
			log.Println("Error creating database or collection:", err)
			return err
		}

		log.Printf("Database '%s' and collection 'init' created successfully.\n", d.MongoDb)
	} else {
		log.Printf("Database '%s' already exists.\n", d.MongoDb)

		// Verificar si la colección existe
		collections, err := d.Mongo.Database(d.MongoDb).ListCollectionNames(context.TODO(), bson.D{})
		if err != nil {
			log.Println("Error listing collections:", err)
			return err
		}

		collectionExists := false
		for _, collection := range collections {
			if collection == "init" {
				collectionExists = true
				break
			}
		}

		// Si la colección no existe, crearla
		if !collectionExists {
			log.Println("Collection 'init' does not exist, creating...")
			collection := d.Mongo.Database(d.MongoDb).Collection("init")
			_, err := collection.InsertOne(context.TODO(), bson.M{"status": "initialized"})
			if err != nil {
				log.Println("Error creating collection:", err)
				return err
			}
			log.Println("Collection 'init' created successfully.")
		} else {
			log.Println("Collection 'init' already exists.")
		}
	}

	return nil
}
