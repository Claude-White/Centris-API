package server

// Coordinate represents latitude and longitude
type Coordinate struct {
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
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

type OutputData struct {
	Link string `json:"link"`
}

type Property struct {
	Id               int
	Title            string
	CivicNumber      string
	Street           string
	AppartmentNumber string
	City             string
	Neighbourhood    string
	Price            string
	Description      string
	NumberOfBedrooms string
	NumberOfRooms    string
	NumberOfToilets  string
	Longitude        string
	Latitude         string
}

type PropertyFeatures struct {
	Id         int
	PropertyId int
	Title      string
	Value      string
}

type PropertyExpenses struct {
	Id           int
	PropertyId   int
	Type         string
	AnnualPrice  string
	MonthlyPrice string
}

type PropertyPhotos struct {
	Id          int
	PropertyId  int
	PhotoLink   string
	Description string
}

// BrokerInfo holds all the information we want to collect
type Broker struct {
	Id           int
	Name         string
	Title        string
	ProfilePhoto string
	// Add other fields as needed
}

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
	PageIndex int `json:"PageIndex"`
	ZoomLevel int `json:"ZoomLevel"`
	Latitude float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	GeoHash string `json:"GeoHash"`
}