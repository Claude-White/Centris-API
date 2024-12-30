package server

import (
	"centris-api/internal/repository"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.HelloWorldHandler)

	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/brokers", s.GetBrokers)
	mux.HandleFunc("/properties", s.GetAllProperties)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT 1"
	var result int
	err := s.db.QueryRow(context.Background(), query).Scan(&result)
	if err != nil {
		log.Printf("Database health check failed: %v", err)
		return
	}
}

func (s *Server) GetBrokers(w http.ResponseWriter, r *http.Request) {
	domain := "www.centris.ca"
	path := "/fr/courtiers-immobiliers"
	numberOfBrokers := getTotalNumberOfBrokers(domain, path)
	startPositions := make([]int, numberOfBrokers)
	for i := 0; i < numberOfBrokers; i++ {
		startPositions[i] = i
	}
	brokers, err := GetBrokersConcurrently(startPositions, 500)
	if err != nil {
		log.Fatalf("Error getting broker info: %v", err)
	}
	resp, err := json.Marshal(brokers)
	if err != nil {
		http.Error(w, "Failed to marshal brokers response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) GetAllProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	properties, err := s.queries.GetAllProperties(ctx, repository.GetAllPropertiesParams{
		Limit:  1,
		Offset: 0,
	})

	if err != nil {
		log.Printf("Failed to get all properties: %v", err)
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal properties response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
