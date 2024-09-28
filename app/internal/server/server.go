package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/AndresKenji/reverse-proxy/internal/config"
	"github.com/AndresKenji/reverse-proxy/internal/database"
	"github.com/AndresKenji/reverse-proxy/internal/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	Port     string
	Mux      *http.ServeMux
	Context  context.Context
	Database *database.Database
	RestartChan chan bool
}

func NewServer(ctx context.Context, restart chan bool) *Server {
	server := &Server{}
	server.Context = ctx
	port := os.Getenv("port")
	if port == "" {
		log.Println("There is no environment variable for server port, server will use port 80 ")
		server.Port = "80"
	} else {
		server.Port = port
	}
	server.Database = database.NewDatabase()
	server.RestartChan = restart

	return server
}

func (s *Server) SetServerMux(cfgFile *config.ConfigFile) {
	mux := http.NewServeMux()

	// Middleware Chain
	chain := middleware.MiddlewareChain(
		middleware.CORSMiddleware,
		middleware.RequestLoggerMiddleware,
	)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server listening"))
	})
	mux.HandleFunc("/restart",func(w http.ResponseWriter, r *http.Request) {
		s.RestartChan <- true
		w.Write([]byte("Restart signal sent"))
	})
	
	fs := http.FileServer(http.Dir("/app/web"))
    mux.Handle("/admin/", http.StripPrefix("/admin", fs))

	mux.Handle("GET /admin/config", chain(http.HandlerFunc(s.GetConfigsHandler)))
	mux.Handle("POST /admin/config", middleware.CORSMiddleware(http.HandlerFunc(s.SaveConfigHandler)))

	for _, cfg := range cfgFile.Endpoints {
		mux.Handle(cfg.Prefix, chain(cfg.GenerateProxyHandler()))
	}
	s.Mux = mux
}

func (s *Server) StartServer() error {
	srv := &http.Server{
		Addr:    ":" + s.Port,
		Handler: s.Mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Printf("Server started on port %s\n", s.Port)

	<-s.Context.Done()

	log.Println("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		return err
	}

	log.Println("Server stopped gracefully.")
	return nil
}

func (s *Server) SaveConfig(configFile *config.ConfigFile) error {
	collection := s.Database.Mongo.Database(s.Database.MongoDb).Collection("configurations")

	_, err := collection.InsertOne(context.TODO(), configFile)
	if err != nil {
		log.Println("Error saving config:", err)
	}

	log.Println("Config saved successfully")
	return nil
}

// GetLatestConfig fetches the latest configuration from MongoDB
func (s *Server) GetLatestConfig() (*config.ConfigFile, error) {
	collection := s.Database.Mongo.Database(s.Database.MongoDb).Collection("configurations")

	// Find the latest document based on the created_at field
	var latestConfig config.ConfigFile
	filter := bson.D{} // You can modify this filter if needed
	options := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	err := collection.FindOne(context.TODO(), filter, options).Decode(&latestConfig)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No documents found, return nil
		}
		return nil, err
	}

	return &latestConfig, nil
}

// GetAllConfigs fetches all configuration documents from MongoDB
func (s *Server) GetAllConfigs() ([]config.ConfigFile, error) {
	collection := s.Database.Mongo.Database(s.Database.MongoDb).Collection("configurations")

	var configs []config.ConfigFile
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var config config.ConfigFile
		if err := cursor.Decode(&config); err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

// UpdateConfig updates a configuration document in MongoDB based on the provided filter
func (s *Server) UpdateConfig(filter bson.D, update bson.D) error {
	collection := s.Database.Mongo.Database(s.Database.MongoDb).Collection("configurations")

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeleteConfig deletes a configuration document from MongoDB based on the provided filter
func (s *Server) DeleteConfig(filter bson.D) ( *mongo.DeleteResult, error) {
	collection := s.Database.Mongo.Database(s.Database.MongoDb).Collection("configurations")

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) SetDefaultConfig() *config.ConfigFile {
	return config.DefaultConfig()
}