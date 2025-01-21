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

	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Property endpoints
	mux.HandleFunc("GET /properties/{mls}", s.GetProperty)
	mux.HandleFunc("POST /properties/coordinates", s.GetPropertyByCoordinates)
	mux.HandleFunc("POST /properties", s.GetAllProperties)
	mux.HandleFunc("POST /properties/broker", s.GetAllBrokerProperties)
	mux.HandleFunc("POST /properties/agency", s.GetAllAgencyProperties)
	mux.HandleFunc("POST /properties/city", s.GetAllCityProperties)
	mux.HandleFunc("POST /properties/radius", s.GetAllRadiusProperties)

	// Broker enpoints
	mux.HandleFunc("GET /brokers/{brokerId}", s.GetBroker)
	mux.HandleFunc("POST /brokers", s.GetAllBrokers)

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

// GetProperty godoc
//
//	@Summary		Get property by MLS number
//	@Description	Retrieves property details using the provided MLS number
//	@Tags			Properties
//	@Accept			json
//	@Produce		json
//	@Param			mls	path		int	true	"MLS number"
//	@Success		200	{object}	repository.Property
//	@Failure		400	{string}	string	"Invalid MLS number"
//	@Failure		404	{string}	string	"Property not found"
//	@Failure		500	{string}	string	"Failed to get property"
//	@Router			/properties/{mls} [get]
func (s *Server) GetProperty(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	mls, err := strconv.ParseInt(strings.TrimSpace(r.PathValue("mls")), 10, 32)
	if err != nil {
		http.Error(w, "Invalid MLS number", http.StatusBadRequest)
		return
	}

	property, err := s.queries.GetProperty(ctx, int32(mls))
	if err != nil {
		http.Error(w, "Property not found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(property)
	if err != nil {
		http.Error(w, "Failed to marshal property response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// GetAllProperties godoc
//
//	@Summary		Get all properties
//	@Description	Retrieves a list of properties with pagination
//	@Tags			Properties
//	@Accept			json
//	@Produce		json
//	@Param			request	body		repository.GetAllPropertiesParams	true	"Pagination parameters"
//	@Success		200		{array}		repository.Property
//	@Failure		400		{object}	string	"Invalid request body"
//	@Failure		500		{object}	string	"Internal server error"
//	@Router			/properties [post]
func (s *Server) GetAllProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req repository.GetAllPropertiesParams
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllProperties(ctx, req)
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

// GetPropertyByCoordinates godoc
//
//	@Summary		Get a property by coordinates
//	@Description	Retrieves a property based on latitude and longitude
//	@Tags			Properties
//	@Accept			json
//	@Produce		json
//	@Param			request	body		repository.GetPropertyByCoordinatesParams	true	"Coordinates parameters"
//	@Success		200		{object}	repository.Property
//	@Failure		400		{object}	string	"Invalid request body"
//	@Failure		404		{object}	string	"Property not found"
//	@Failure		500		{object}	string	"Internal server error"
//	@Router			/properties/coordinates [post]
func (s *Server) GetPropertyByCoordinates(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req repository.GetPropertyByCoordinatesParams
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetPropertyByCoordinates(ctx, req)
	if err != nil {
		log.Printf("Failed to get properties: %v", err)
		http.Error(w, "Failed to get properties", http.StatusInternalServerError)
	}
	if len(properties) == 0 {
		http.Error(w, "Property not found", http.StatusNotFound)
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal property response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// GetAllBrokerProperties godoc
//
//	@Summary		Get all properties by broker
//	@Description	Retrieves a list of properties for a specific broker with pagination
//	@Tags			Properties
//	@Accept			json
//	@Produce		json
//	@Param			request	body		repository.GetAllBrokerPropertiesParams	true	"Pagination parameters"
//	@Success		200		{array}		repository.Property
//	@Failure		400		{object}	string	"Invalid broker ID or request body"
//	@Failure		404		{object}	string	"Broker properties not found"
//	@Failure		500		{object}	string	"Internal server error"
//	@Router			/properties/broker [post]
func (s *Server) GetAllBrokerProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req repository.GetAllBrokerPropertiesParams
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllBrokerProperties(ctx, req)
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

// GetAllAgencyProperties godoc
//
//	@Summary		Get all properties by agency
//	@Description	Retrieves a list of properties for a specific agency with pagination
//	@Tags			Properties
//	@Accept			json
//	@Produce		json
//	@Param			request	body		repository.GetAllAgencyPropertiesParams	true	"Pagination parameters"
//	@Success		200		{array}		repository.Property
//	@Failure		400		{object}	string	"Invalid agency name or request body"
//	@Failure		404		{object}	string	"Agency properties not found"
//	@Failure		500		{object}	string	"Internal server error"
//	@Router			/properties/agency [post]
func (s *Server) GetAllAgencyProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req repository.GetAllAgencyPropertiesParams
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllAgencyProperties(ctx, req)
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

// GetAllCategoryProperties godoc
//
//	@Summary		Get all properties by category
//	@Description	Retrieves a list of properties for a specific category with pagination
//	@Tags			Properties
//	@Accept			json
//	@Produce		json
//	@Param			request	body		repository.GetAllCategoryPropertiesParams	true	"Pagination parameters"
//	@Success		200		{array}		repository.Property
//	@Failure		400		{object}	string	"Invalid category name or request body"
//	@Failure		404		{object}	string	"Category properties not found"
//	@Failure		500		{object}	string	"Internal server error"
//	@Router			/properties/category [post]
func (s *Server) GetAllCategoryProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req repository.GetAllCategoryPropertiesParams
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllCategoryProperties(ctx, req)
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

// GetAllCityProperties godoc
//
//	@Summary		Get all properties by city
//	@Description	Retrieves a list of properties for a specific city with pagination
//	@Tags			Properties
//	@Accept			json
//	@Produce		json
//	@Param			request	body		repository.GetAllCityPropertiesParams	true	"Pagination parameters"
//	@Success		200		{array}		repository.Property
//	@Failure		400		{object}	string	"Invalid city name or request body"
//	@Failure		404		{object}	string	"City properties not found"
//	@Failure		500		{object}	string	"Internal server error"
//	@Router			/properties/city [post]
func (s *Server) GetAllCityProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req repository.GetAllCityPropertiesParams
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllCityProperties(ctx, req)
	if err != nil {
		http.Error(w, "Failed to get city properties", http.StatusInternalServerError)
		return
	}

	if properties == nil {
		http.Error(w, "City properties not found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(properties)
	if err != nil {
		http.Error(w, "Failed to marshal city properties response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// GetAllRadiusProperties godoc
//
//	@Summary		Get all properties by radius
//	@Description	Retrieves a list of properties for a specific radius with pagination
//	@Tags			Properties
//	@Accept			json
//	@Produce		json
//	@Param			request	body		repository.GetAllRadiusPropertiesParams	true	"Pagination parameters"
//	@Success		200		{array}		repository.Property
//	@Failure		400		{object}	string	"Invalid radius or request body"
//	@Failure		404		{object}	string	"Properties not found"
//	@Failure		500		{object}	string	"Internal server error"
//	@Router			/properties/radius [post]
func (s *Server) GetAllRadiusProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req repository.GetAllRadiusPropertiesParams
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	properties, err := s.queries.GetAllRadiusProperties(ctx, req)
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

// GetAllBrokers godoc
//
//	@Summary		Get all brokers
//	@Description	Retrieves a list of brokers with pagination
//	@Tags			Brokers
//	@Accept			json
//	@Produce		json
//	@Param			request	body		repository.GetAllBrokersParams	true	"Pagination parameters"
//	@Success		200		{array}		repository.Broker
//	@Failure		400		{object}	string	"Invalid request body"
//	@Failure		500		{object}	string	"Internal server error"
//	@Router			/brokers [post]
func (s *Server) GetAllBrokers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req repository.GetAllBrokersParams
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.NumberOfItems < 1 {
		http.Error(w, "Number of items must be above 0", http.StatusBadRequest)
		return
	}

	if req.NumberOfItems > 20 {
		http.Error(w, "Number of items must be below 21", http.StatusBadRequest)
		return
	}

	brokers, err := s.queries.GetAllBrokers(ctx, req)
	if err != nil {
		http.Error(w, "Failed to get brokers", http.StatusInternalServerError)
		return
	}

	var completeBrokers []CompleteBroker
	for _, broker := range brokers {
		completeBroker, err := GetCompleteBroker(s, ctx, broker)
		if err != nil {
			http.Error(w, "Complete broker not found", http.StatusNotFound)
			return
		}
		completeBrokers = append(completeBrokers, completeBroker)
	}

	resp, err := json.Marshal(completeBrokers)
	if err != nil {
		http.Error(w, "Failed to marshal brokers response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// GetBroker godoc
//
//	@Summary		Get broker by broker number
//	@Description	Retrieves broker details using the provided broker number
//	@Tags			Brokers
//	@Accept			json
//	@Produce		json
//	@Param			brokerId	path		int	true	"Broker number"
//	@Success		200			{object}	repository.Broker
//	@Failure		400			{string}	string	"Invalid broker number"
//	@Failure		404			{string}	string	"Broker not found"
//	@Failure		500			{string}	string	"Failed to get broker"
//	@Router			/brokers/{brokerId} [get]
func (s *Server) GetBroker(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	brokerId, err := strconv.ParseInt(strings.TrimSpace(r.PathValue("brokerId")), 10, 32)
	if err != nil {
		http.Error(w, "Invalid MLS number", http.StatusBadRequest)
		return
	}

	var completeBroker CompleteBroker

	broker, err := s.queries.GetBroker(ctx, int32(brokerId))
	if err != nil {
		http.Error(w, "Broker not found", http.StatusNotFound)
		return
	}

	completeBroker, err = GetCompleteBroker(s, ctx, broker)
	if err != nil {
		http.Error(w, "Complete broker not found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(completeBroker)
	if err != nil {
		http.Error(w, "Failed to marshal broker response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
