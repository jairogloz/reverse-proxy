package server

import (
	"encoding/json"
	"net/http"

	"github.com/AndresKenji/reverse-proxy/internal/config"
	"go.mongodb.org/mongo-driver/bson"
)


func (s *Server) GetConfigsHandler(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("query")
	if param == "latest" {
		cfgFile, err := s.GetLatestConfig()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		out, err := json.Marshal(cfgFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(out)
		return
	}
	cfgFiles, err := s.GetAllConfigs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out, err := json.Marshal(cfgFiles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(out)

}

func (s *Server) SaveConfigHandler(w http.ResponseWriter, r *http.Request) {
	var cfgFile config.ConfigFile
	err := json.NewDecoder(r.Body).Decode(&cfgFile)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	err = s.SaveConfig(&cfgFile)
	if err != nil {
		http.Error(w, "Failed to save config file", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Config saved"))
}

func (s *Server) DeleteConfigHandler(w http.ResponseWriter, r *http.Request) {
    
    // Extraer el ID de la URL
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }
    
    // Eliminar el documento
    filter := bson.D{{Key: "id", Value: id}}
	result, err := s.DeleteConfig(filter)
	if err != nil {
        http.Error(w, "Failed to delete document", http.StatusInternalServerError)
        return
    }
	if result.DeletedCount == 0 {
        http.Error(w, "No document found with the given ID", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Document deleted successfully"))
}

func (s *Server) UpdateConfigHandler(w http.ResponseWriter, r *http.Request) {
	 // Extraer el ID de la URL
	 id := r.URL.Query().Get("id")
	 if id == "" {
		 http.Error(w, "ID is required", http.StatusBadRequest)
		 return
	}
	var updatedConfig config.Config
    err := json.NewDecoder(r.Body).Decode(&updatedConfig)
    if err != nil {
        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
        return
    }
	update := bson.D{
        {Key: "$set", Value: bson.D{
            {Key: "prefix", Value: updatedConfig.Prefix},
            {Key: "header_identifier", Value: updatedConfig.HeaderIdentifier},
            {Key: "backend_urls", Value: updatedConfig.BackendUrls},
            {Key: "secure", Value: updatedConfig.Secure},
        }},
    }
	 
	// Eliminar el documento
	filter := bson.D{{Key: "id", Value: id}}
	err = s.UpdateConfig(filter, update)
	if err != nil {
        http.Error(w, "Failed to update document", http.StatusInternalServerError)
        return
    }

	w.WriteHeader(http.StatusOK)
    w.Write([]byte("Document updated successfully"))
}