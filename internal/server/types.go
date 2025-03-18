package server

import (
	"centris-api/internal/repository"

	"github.com/google/uuid"
)

// Coordinate represents latitude and longitude
type Coordinate struct {
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
}

// Result wraps a value and potential error
type Result[T any] struct {
	Value T
	Err   error
}

// MapBounds represents the geographic bounds
type MapBounds struct {
	SouthWest Coordinate `json:"SouthWest"`
	NorthEast Coordinate `json:"NorthEast"`
}

// InputData represents the structure of each item in test.json
type InputData struct {
	ZoomLevel int       `json:"zoomLevel"`
	MapBounds MapBounds `json:"mapBounds"`
}

// MarkerPosition represents the position of a marker
type MarkerPosition struct {
	Position Coordinate `json:"Position"`
}

// MarkerResponse represents the response from the GetMarkers endpoint
type MarkerResponse struct {
	D struct {
		Result struct {
			Markers []MarkerPosition `json:"Markers"`
		} `json:"Result"`
	} `json:"d"`
}

// MarkerInfoResponse represents the response from the GetMarkerInfo endpoint
type MarkerInfoResponse struct {
	D struct {
		Result struct {
			Html string `json:"Html"`
		} `json:"Result"`
	} `json:"d"`
}

type BrokerResponse struct {
	D struct {
		Result struct {
			Html  string `json:"Html"`
			Title string `json:"Title"`
		} `json:"Result"`
		Succeeded bool `json:"Succeeded"`
	} `json:"d"`
}

// ResponseData contains the message, result, and succeeded fields
type ResponseData struct {
	Message   string     `json:"Message"`
	Result    ResultData `json:"Result"`
	Succeeded bool       `json:"Succeeded"`
}

// ResultData contains the HTML content and title
type ResultData struct {
	Html  string `json:"Html"`
	Title string `json:"Title"`
}

type OutputData struct {
	Link string `json:"link"`
}

// BrokerInfo holds all the information we want to collect

type Response struct {
	D struct {
		Result struct {
			Markers []Marker `json:"Markers"`
		} `json:"Result"`
	} `json:"d"`
}

type Marker struct {
	ContainsSubject     bool     `json:"ContainsSubject"`
	GeoHash             string   `json:"GeoHash"`
	HasStrictQueryMatch bool     `json:"HasStrictQueryMatch"`
	Key                 Key      `json:"Key"`
	NoMls               *string  `json:"NoMls"`
	PointsCount         int      `json:"PointsCount"`
	Position            Position `json:"Position"`
	SubjectIndex        int      `json:"SubjectIndex"`
	Title               *string  `json:"Title"`
}

type Key struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

type Position struct {
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
}

type MarkerInfoInputData struct {
	PageIndex int     `json:"PageIndex"`
	ZoomLevel int     `json:"ZoomLevel"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	GeoHash   string  `json:"GeoHash"`
}

type RequestBodyPhoto struct {
	Lang                   string `json:"lang"`
	CentrisNo              int64  `json:"centrisNo"`
	Track                  bool   `json:"track"`
	AuthorizationMediaCode string `json:"authorizationMediaCode"`
}

type Resolution struct {
	Domain      string `json:"Domain"`
	Width       int    `json:"Width"`
	Height      int    `json:"Height"`
	QueryParams string `json:"QueryParams"`
}

type PropertyResponsePhoto struct {
	ID        string `json:"Id"`
	UrlThumb  string `json:"UrlThumb"`
	Desc      string `json:"Desc"`
	MaxWidth  int    `json:"MaxWidth"`
	MaxHeight int    `json:"MaxHeight"`
	DescSupp  string `json:"DescSupp,omitempty"`
}

type PhotoResponse struct {
	PhotoList                        []PropertyResponsePhoto `json:"PhotoList"`
	UseIntegralForIdenticalDimension bool                    `json:"UseIntegralForIdenticalDimension"`
	Resolutions                      []Resolution            `json:"Resolutions"`
	Track                            bool                    `json:"Track"`
	CentrisNo                        string                  `json:"CentrisNo"`
	VirtualTourUrl                   *string                 `json:"VirtualTourUrl"`
}

type CompleteBroker struct {
	Broker        repository.Broker
	Broker_Phones []repository.BrokerPhone
	Broker_Links  []repository.BrokerExternalLink
}

type Broker struct {
	ID                int64                `bson:"id,omitempty"`
	FirstName         string               `bson:"first_name,omitempty"`
	MiddleName        string               `bson:"middle_name,omitempty"`
	LastName          string               `bson:"last_name,omitempty"`
	Title             string               `bson:"title,omitempty"`
	ProfilePhoto      string               `bson:"profile_photo,omitempty"`
	ComplementaryInfo string               `bson:"complementary_info,omitempty"`
	ServedAreas       string               `bson:"served_areas,omitempty"`
	Presentation      string               `bson:"presentation,omitempty"`
	CorporationName   string               `bson:"corporation_name,omitempty"`
	AgencyName        string               `bson:"agency_name,omitempty"`
	AgencyAddress     string               `bson:"agency_address,omitempty"`
	AgencyLogo        *string              `bson:"agency_logo,omitempty"`
	ExternalLinks     []BrokerExternalLink `bson:"external_links,omitempty"`
	PhoneNumbers      []BrokerPhone        `bson:"phone_numbers,omitempty"`
	Properties        []Property           `bson:"properties,omitempty"`
}

type BrokerExternalLink struct {
	ID   uuid.UUID `bson:"id,omitempty"`
	Type string    `bson:"type,omitempty"`
	Link string    `bson:"link,omitempty"`
}

type BrokerPhone struct {
	ID     uuid.UUID `bson:"id,omitempty"`
	Type   string    `bson:"type,omitempty"`
	Number string    `bson:"number,omitempty"`
}

type Property struct {
	// MLS number
	ID             int64             `bson:"mls,omitempty"`
	Title          string            `bson:"title,omitempty"`
	Category       string            `bson:"category,omitempty"`
	Address        string            `bson:"address,omitempty"`
	CityName       string            `bson:"city_name,omitempty"`
	Price          float32           `bson:"price,omitempty"`
	Description    *string           `bson:"description,omitempty"`
	BedroomNumber  *int32            `bson:"bedroom_number,omitempty"`
	RoomNumber     *int32            `bson:"room_number,omitempty"`
	BathroomNumber *int32            `bson:"bathroom_number,omitempty"`
	Latitude       float32           `bson:"latitude,omitempty"`
	Longitude      float32           `bson:"longitude,omitempty"`
	Expenses       []PropertyExpense `bson:"expenses,omitempty"`
	Features       []PropertyFeature `bson:"features,omitempty"`
	Photos         []PropertyPhoto   `bson:"photos,omitempty"`
	BrokerIds      []int64           `bson:"broker_ids,omitempty"`
}

type PropertyExpense struct {
	ID           uuid.UUID `bson:"id,omitempty"`
	Type         string    `bson:"type,omitempty"`
	AnnualPrice  float32   `bson:"annual_price,omitempty"`
	MonthlyPrice float32   `bson:"monthly_price,omitempty"`
}

type PropertyFeature struct {
	ID    uuid.UUID `bson:"id,omitempty"`
	Title string    `bson:"title,omitempty"`
	Value string    `bson:"value,omitempty"`
}

type PropertyPhoto struct {
	ID          uuid.UUID `bson:"id,omitempty"`
	Link        string    `bson:"link,omitempty"`
	Description *string   `bson:"description,omitempty"`
}
