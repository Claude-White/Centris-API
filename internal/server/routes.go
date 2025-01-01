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

type RequestBody struct {
	StartPosition int32 `json:"start_position"`
	NumberOfItems int32 `json:"number_of_items"`
}
type Coordinates struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

type GeoFilterCoordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Radius    string  `json:"radius"`
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Property endpoints
	mux.HandleFunc("GET /properties/{mls}", s.GetProperty)
	mux.HandleFunc("POST /properties/coordinates", s.GetPropertyByCoordinates)
	mux.HandleFunc("POST /properties", s.GetAllProperties)
	mux.HandleFunc("POST /properties/broker/{brokerId}", s.GetAllBrokerProperties)
	mux.HandleFunc("POST /properties/agency/{agencyName}", s.GetAllAgencyProperties)
	mux.HandleFunc("POST /properties/city/{cityName}", s.GetAllCityProperties)
	mux.HandleFunc("POST /properties/neighbourhood/{neighbourhoodName}", s.GetAllNeighbourhoodProperties)
	mux.HandleFunc("POST /properties/geo-filter/{radius}", s.GetAllRadiusProperties)
	mux.HandleFunc("POST /properties/create", s.CreateProperty)

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

// GetProperty godoc
// @Summary      Get property by MLS number
// @Description  Retrieves property details using the provided MLS number
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        mls   path      int     true  "MLS number"
// @Success      200   {object}  repository.Property
// @Failure      400   {string}  string  "Invalid MLS number"
// @Failure      404   {string}  string  "Property not found"
// @Failure      500   {string}  string  "Failed to get property"
// @Router       /properties/{mls} [get]
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

// GetAllProperties godoc
// @Summary      Get all properties
// @Description  Retrieves a list of properties with pagination
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        request body RequestBody true "Pagination parameters"
// @Success      200 {array} repository.Property
// @Failure      400 {object} string "Invalid request body"
// @Failure      500 {object} string "Internal server error"
// @Router       /properties [post]
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

// GetPropertyByCoordinates godoc
// @Summary      Get a property by coordinates
// @Description  Retrieves a property based on latitude and longitude
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        request body Coordinates true "Coordinates parameters"
// @Success      200 {object} repository.Property
// @Failure      400 {object} string "Invalid request body"
// @Failure      404 {object} string "Property not found"
// @Failure      500 {object} string "Internal server error"
// @Router       /properties/coordinates [post]
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

// GetAllBrokerProperties godoc
// @Summary      Get all properties by broker
// @Description  Retrieves a list of properties for a specific broker with pagination
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        brokerId path int true "Broker ID"
// @Param        request body RequestBody true "Pagination parameters"
// @Success      200 {array} repository.Property
// @Failure      400 {object} string "Invalid broker ID or request body"
// @Failure      404 {object} string "Broker properties not found"
// @Failure      500 {object} string "Internal server error"
// @Router       /properties/broker/{brokerId} [post]
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
		BrokerID:      brokerId,
		NumberOfItems: req.NumberOfItems,
		StartPosition: req.StartPosition,
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

// GetAllAgencyProperties godoc
// @Summary      Get all properties by agency
// @Description  Retrieves a list of properties for a specific agency with pagination
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        agencyName path string true "Agency Name"
// @Param        request body RequestBody true "Pagination parameters"
// @Success      200 {array} repository.Property
// @Failure      400 {object} string "Invalid agency name or request body"
// @Failure      404 {object} string "Agency properties not found"
// @Failure      500 {object} string "Internal server error"
// @Router       /properties/agency/{agencyName} [post]
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
		AgencyName:    agency,
		NumberOfItems: req.NumberOfItems,
		StartPosition: req.StartPosition,
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

// GetAllCategoryProperties godoc
// @Summary      Get all properties by category
// @Description  Retrieves a list of properties for a specific category with pagination
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        categoryName path string true "Category Name"
// @Param        request body RequestBody true "Pagination parameters"
// @Success      200 {array} repository.Property
// @Failure      400 {object} string "Invalid category name or request body"
// @Failure      404 {object} string "Category properties not found"
// @Failure      500 {object} string "Internal server error"
// @Router       /properties/category/{categoryName} [post]
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
		Category:      category,
		NumberOfItems: req.NumberOfItems,
		StartPosition: req.StartPosition,
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

// GetAllCityProperties godoc
// @Summary      Get all properties by city
// @Description  Retrieves a list of properties for a specific city with pagination
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        cityName path string true "City Name"
// @Param        request body RequestBody true "Pagination parameters"
// @Success      200 {array} repository.Property
// @Failure      400 {object} string "Invalid city name or request body"
// @Failure      404 {object} string "City properties not found"
// @Failure      500 {object} string "Internal server error"
// @Router       /properties/city/{cityName} [post]
func (s *Server) GetAllCityProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var city *string
	cityName := strings.TrimSpace(r.PathValue("cityName"))
	if cityName == "" {
		http.Error(w, "City name is required", http.StatusBadRequest)
		return
	} else {
		lowerCaseCityName := strings.ToLower(cityName)
		city = &lowerCaseCityName
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

	properties, err := s.queries.GetAllCityProperties(ctx, repository.GetAllCityPropertiesParams{
		CityName:      city,
		NumberOfItems: req.NumberOfItems,
		StartPosition: req.StartPosition,
	})
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

// GetAllNeighbourhoodProperties godoc
// @Summary      Get all properties by neighbourhood
// @Description  Retrieves a list of properties for a specific neighbourhood with pagination
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        neighbourhoodName path string true "Neighbourhood Name"
// @Param        request body RequestBody true "Pagination parameters"
// @Success      200 {array} repository.Property
// @Failure      400 {object} string "Invalid neighbourhood name or request body"
// @Failure      404 {object} string "Neighbourhood properties not found"
// @Failure      500 {object} string "Internal server error"
// @Router       /properties/neighbourhood/{neighbourhoodName} [post]
func (s *Server) GetAllNeighbourhoodProperties(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var neighbourhood *string
	neighbourhoodName := strings.TrimSpace(r.PathValue("neighbourhoodName"))
	if neighbourhoodName == "" {
		http.Error(w, "Neighbourhood name is required", http.StatusBadRequest)
		return
	} else {
		lowerCaseCityName := strings.ToLower(neighbourhoodName)
		neighbourhood = &lowerCaseCityName
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

	properties, err := s.queries.GetAllNeighbourhoodProperties(ctx, repository.GetAllNeighbourhoodPropertiesParams{
		NeighbourhoodName: neighbourhood,
		NumberOfItems:     req.NumberOfItems,
		StartPosition:     req.StartPosition,
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

	floatValue, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		http.Error(w, "Invalid radius value", http.StatusBadRequest)
	}

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
		Radius:    floatValue,
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

// CreateProperty godoc
// @Summary      Create a new property
// @Description  Creates a new property given a json request body
// @Tags         dev
// @Accept       json
// @Produce      json
// @Param        params body repository.CreatePropertyParams true "Property request body"
// @Success      200 {integer} string "New property id"
// @Failure      400 {object} string "Invalid request body"
// @Failure      500 {object} string "Internal server error"
// @Router       /properties/create [post]
func (s *Server) CreateProperty(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req repository.CreatePropertyParams
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	propertyId, err := s.queries.CreateProperty(ctx, req)
	if err != nil {
		http.Error(w, "Failed to create property", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(propertyId)
	if err != nil {
		http.Error(w, "Failed to marshal create property response", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// ********************** BROKER ENDPOINT FUNCTIONS **********************

// GetBrokers godoc
// @Summary      Get all brokers
// @Description  Retrieves a list of all brokers
// @Tags         Brokers
// @Produce      json
// @Success      200 {array} Broker
// @Failure      500 {object} string "Internal server error"
// @Router       /brokers [get]
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
