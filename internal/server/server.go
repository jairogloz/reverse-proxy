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
	Port    string
	Mux     *http.ServeMux
	Context context.Context
	Database *database.Database
}

func NewServer(ctx context.Context) *Server {
	server := &Server{}
	server.Context = ctx
	port := os.Getenv("port")
	if port == "" {
		log.Println("There is no environment variable for server port, server will use port 8080 ")
		server.Port = "8080"
	} else {
		server.Port = port
	}
	server.Database = database.NewDatabase()

	return server
}

func (s *Server) SetServerMux(cfgFile *config.ConfigFile) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server listening"))
	})
	for _, cfg := range cfgFile.Endpoints {
		log.Println(cfg)
		mux.HandleFunc(cfg.Prefix, middleware.RequestLoggerMiddleware(cfg.GenerateHandler()))
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
	collection := s.Database.Mongo.Database(os.Getenv("mongo_db")).Collection("configurations")

	_, err := collection.InsertOne(context.TODO(), configFile)
	if err != nil {
		log.Println("Error saving config:",err)
	}

	log.Println("Config saved successfully")
	return nil
}

// GetLatestConfig fetches the latest configuration from MongoDB
func (s *Server) GetLatestConfig() (*config.ConfigFile, error) {
	collection := s.Database.Mongo.Database(os.Getenv("mongo_db")).Collection("configurations")

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

func (s *Server) UpdateConfig(filter bson.M, update bson.M) error {
	collection := s.Database.Mongo.Database(os.Getenv("mongo_db")).Collection("configurations")

	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set":update})
	if err != nil {
		log.Println("Error updating config:",err)
	}

	log.Println("Config updated successfully!")
	return nil
}