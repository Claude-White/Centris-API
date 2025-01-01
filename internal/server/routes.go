package server

import (
	"centris-api/internal/repository"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type RequestBody struct {
	StartPosition int32
	NumberOfItems int32
}

type Coordinates struct {
	Longitude pgtype.Numeric
	Latitude  pgtype.Numeric
}

type GeoFilterCoordinates struct {
	Longitude float64
	Latitude  float64
	Radius    pgtype.Numeric
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	// Property endpoints
	mux.HandleFunc("GET /properties/{mls}", s.GetProperty)
	mux.HandleFunc("POST /properties/coordinates", s.GetPropertyByCoordinates)
	mux.HandleFunc("POST /properties", s.GetAllProperties)
	mux.HandleFunc("POST /properties/broker/{brokerId}", s.GetAllBrokerProperties)
	mux.HandleFunc("POST /properties/agency/{agencyName}", s.GetAllAgencyProperties)
	mux.HandleFunc("POST /properties/city/{cityName}", s.GetAllCityProperties)
	mux.HandleFunc("POST /properties/neighbourhood/{NeighbourhoodName}", s.GetAllNeighbourhoodProperties)
	mux.HandleFunc("POST /properties/geo-filter/{radius}", s.GetAllRadiusProperties)

	// Broker enpoints
	mux.HandleFunc("/brokers", s.GetBrokers)

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

// ********************** PROPERTY ENDPOINT FUNCTIONS **********************
func (s *Server) GetProperty(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	mls, err := strconv.ParseInt(strings.TrimSpace(r.PathValue("mls")), 10, 64)
	if err != nil {
		http.Error(w, "Invalid MLS number", http.StatusBadRequest)
	}

	property, err := s.queries.GetProperty(ctx, mls)
	if err != nil {
		log.Printf("Failed to get property: %v", err)
		http.Error(w, "Failed to get property", http.StatusInternalServerError)
	}

	if property == (repository.Property{}) {
		http.Error(w, "Property not found", http.StatusNotFound)
	}

	resp, err := json.Marshal(property)
	if err != nil {
		http.Error(w, "Failed to marshal property response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) GetAllProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req RequestBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllProperties(ctx, repository.GetAllPropertiesParams{
		Limit:  req.NumberOfItems,
		Offset: req.StartPosition,
	})
	if err != nil {
		log.Printf("Failed to get properties: %v", err)
		http.Error(w, "Failed to get properties", http.StatusInternalServerError)
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal properties response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) GetPropertyByCoordinates(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req Coordinates
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	property, err := s.queries.GetPropertyByCoordinates(ctx, repository.GetPropertyByCoordinatesParams{
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
	})
	if err != nil {
		log.Printf("Failed to get properties: %v", err)
		http.Error(w, "Failed to get properties", http.StatusInternalServerError)
	}
	if property == (repository.Property{}) {
		http.Error(w, "Property not found", http.StatusNotFound)
	}

	resp, err := json.Marshal(property)
	if err != nil {
		http.Error(w, "Failed to marshal property response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) GetAllBrokerProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	brokerId, err := strconv.ParseInt(strings.TrimSpace(r.PathValue("brokerId")), 10, 64)
	if err != nil {
		http.Error(w, "Invalid broker id", http.StatusBadRequest)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req RequestBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllBrokerProperties(ctx, repository.GetAllBrokerPropertiesParams{
		BrokerID: brokerId,
		Limit:    req.NumberOfItems,
		Offset:   req.StartPosition,
	})
	if err != nil {
		log.Printf("Failed to get broker properties: %v", err)
		http.Error(w, "Failed to get broker properties", http.StatusInternalServerError)
	}

	if properties == nil {
		http.Error(w, "Broker properties not found", http.StatusNotFound)
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal broker properties response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) GetAllAgencyProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	agency := strings.ToLower(strings.TrimSpace(r.PathValue("agencyName")))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req RequestBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllAgencyProperties(ctx, repository.GetAllAgencyPropertiesParams{
		AgencyName: agency,
		Limit:      req.NumberOfItems,
		Offset:     req.StartPosition,
	})
	if err != nil {
		log.Printf("Failed to get agency properties: %v", err)
		http.Error(w, "Failed to get agency properties", http.StatusInternalServerError)
	}

	if properties == nil {
		http.Error(w, "Agency properties not found", http.StatusNotFound)
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal agency properties response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) GetAllCategoryProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	category := strings.ToLower(strings.TrimSpace(r.PathValue("categoryName")))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req RequestBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllCategoryProperties(ctx, repository.GetAllCategoryPropertiesParams{
		Category: category,
		Limit:    req.NumberOfItems,
		Offset:   req.StartPosition,
	})
	if err != nil {
		log.Printf("Failed to get category properties: %v", err)
		http.Error(w, "Failed to get category properties", http.StatusInternalServerError)
	}

	if properties == nil {
		http.Error(w, "Category properties not found", http.StatusNotFound)
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal category properties response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) GetAllCityProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	city := pgtype.Text{String: strings.ToLower(strings.TrimSpace(r.PathValue("cityName"))), Valid: true}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req RequestBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllCityProperties(ctx, repository.GetAllCityPropertiesParams{
		CityName: city,
		Limit:    req.NumberOfItems,
		Offset:   req.StartPosition,
	})
	if err != nil {
		log.Printf("Failed to get city properties: %v", err)
		http.Error(w, "Failed to get city properties", http.StatusInternalServerError)
	}

	if properties == nil {
		http.Error(w, "City properties not found", http.StatusNotFound)
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal city properties response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) GetAllNeighbourhoodProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	neighbourhood := pgtype.Text{String: strings.ToLower(strings.TrimSpace(r.PathValue("NeighbourhoodName"))), Valid: true}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req RequestBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllNeighbourhoodProperties(ctx, repository.GetAllNeighbourhoodPropertiesParams{
		NeighbourhoodName: neighbourhood,
		Limit:             req.NumberOfItems,
		Offset:            req.StartPosition,
	})
	if err != nil {
		log.Printf("Failed to get neighbourhood properties: %v", err)
		http.Error(w, "Failed to get neighbourhood properties", http.StatusInternalServerError)
	}

	if properties == nil {
		http.Error(w, "Neighbourhood properties not found", http.StatusNotFound)
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal neighbourhood properties response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) GetAllRadiusProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	radius := strings.TrimSpace(r.PathValue("radius"))
	var radiusNumeric pgtype.Numeric

	floatValue, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		http.Error(w, "Invalid radius value", http.StatusBadRequest)
	}

	radiusNumeric.Scan(strconv.FormatFloat(floatValue, 'f', -1, 64))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req GeoFilterCoordinates
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllRadiusProperties(ctx, repository.GetAllRadiusPropertiesParams{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Radius:    radiusNumeric,
	})
	if err != nil {
		log.Printf("Failed to get radius properties: %v", err)
		http.Error(w, "Failed to get radius properties", http.StatusInternalServerError)
	}

	if properties == nil {
		http.Error(w, "Radius properties not found", http.StatusNotFound)
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal radius properties response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// ********************** BROKER ENDPOINT FUNCTIONS **********************
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
