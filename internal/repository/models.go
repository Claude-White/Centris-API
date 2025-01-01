// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BrokerLinkType string

const (
	BrokerLinkTypeWebsite   BrokerLinkType = "website"
	BrokerLinkTypeLinkedin  BrokerLinkType = "linkedin"
	BrokerLinkTypeFacebook  BrokerLinkType = "facebook"
	BrokerLinkTypeTwitter   BrokerLinkType = "twitter"
	BrokerLinkTypeInstagram BrokerLinkType = "instagram"
	BrokerLinkTypeYoutube   BrokerLinkType = "youtube"
)

func (e *BrokerLinkType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BrokerLinkType(s)
	case string:
		*e = BrokerLinkType(s)
	default:
		return fmt.Errorf("unsupported scan type for BrokerLinkType: %T", src)
	}
	return nil
}

type NullBrokerLinkType struct {
	BrokerLinkType BrokerLinkType `json:"broker_link_type"`
	Valid          bool           `json:"valid"` // Valid is true if BrokerLinkType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBrokerLinkType) Scan(value interface{}) error {
	if value == nil {
		ns.BrokerLinkType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BrokerLinkType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBrokerLinkType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.BrokerLinkType), nil
}

type ExpenseType string

const (
	ExpenseTypeTax         ExpenseType = "tax"
	ExpenseTypeInsurance   ExpenseType = "insurance"
	ExpenseTypeMaintenance ExpenseType = "maintenance"
	ExpenseTypeUtilities   ExpenseType = "utilities"
	ExpenseTypeOther       ExpenseType = "other"
)

func (e *ExpenseType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ExpenseType(s)
	case string:
		*e = ExpenseType(s)
	default:
		return fmt.Errorf("unsupported scan type for ExpenseType: %T", src)
	}
	return nil
}

type NullExpenseType struct {
	ExpenseType ExpenseType `json:"expense_type"`
	Valid       bool        `json:"valid"` // Valid is true if ExpenseType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullExpenseType) Scan(value interface{}) error {
	if value == nil {
		ns.ExpenseType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ExpenseType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullExpenseType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ExpenseType), nil
}

type PhoneType string

const (
	PhoneTypeMobile PhoneType = "mobile"
	PhoneTypeOffice PhoneType = "office"
	PhoneTypeHome   PhoneType = "home"
	PhoneTypeFax    PhoneType = "fax"
)

func (e *PhoneType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PhoneType(s)
	case string:
		*e = PhoneType(s)
	default:
		return fmt.Errorf("unsupported scan type for PhoneType: %T", src)
	}
	return nil
}

type NullPhoneType struct {
	PhoneType PhoneType `json:"phone_type"`
	Valid     bool      `json:"valid"` // Valid is true if PhoneType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPhoneType) Scan(value interface{}) error {
	if value == nil {
		ns.PhoneType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PhoneType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPhoneType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PhoneType), nil
}

type PropertyCategory string

const (
	PropertyCategoryHouse      PropertyCategory = "house"
	PropertyCategoryApartment  PropertyCategory = "apartment"
	PropertyCategoryCondo      PropertyCategory = "condo"
	PropertyCategoryLand       PropertyCategory = "land"
	PropertyCategoryCommercial PropertyCategory = "commercial"
)

func (e *PropertyCategory) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PropertyCategory(s)
	case string:
		*e = PropertyCategory(s)
	default:
		return fmt.Errorf("unsupported scan type for PropertyCategory: %T", src)
	}
	return nil
}

type NullPropertyCategory struct {
	PropertyCategory PropertyCategory `json:"property_category"`
	Valid            bool             `json:"valid"` // Valid is true if PropertyCategory is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPropertyCategory) Scan(value interface{}) error {
	if value == nil {
		ns.PropertyCategory, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PropertyCategory.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPropertyCategory) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PropertyCategory), nil
}

type PropertyStatus string

const (
	PropertyStatusAvailable PropertyStatus = "available"
	PropertyStatusPending   PropertyStatus = "pending"
	PropertyStatusSold      PropertyStatus = "sold"
	PropertyStatusRented    PropertyStatus = "rented"
	PropertyStatusOffMarket PropertyStatus = "off_market"
)

func (e *PropertyStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PropertyStatus(s)
	case string:
		*e = PropertyStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PropertyStatus: %T", src)
	}
	return nil
}

type NullPropertyStatus struct {
	PropertyStatus PropertyStatus `json:"property_status"`
	Valid          bool           `json:"valid"` // Valid is true if PropertyStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPropertyStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PropertyStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PropertyStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPropertyStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PropertyStatus), nil
}

type Broker struct {
	ID                int64      `json:"id"`
	FirstName         string     `json:"first_name"`
	MiddleName        *string    `json:"middle_name"`
	LastName          string     `json:"last_name"`
	Title             string     `json:"title"`
	ProfilePhoto      *string    `json:"profile_photo"`
	ComplementaryInfo *string    `json:"complementary_info"`
	ServedAreas       *string    `json:"served_areas"`
	Presentation      *string    `json:"presentation"`
	CorporationName   *string    `json:"corporation_name"`
	AgencyName        string     `json:"agency_name"`
	AgencyAddress     string     `json:"agency_address"`
	AgencyLogo        *string    `json:"agency_logo"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

type BrokerExternalLink struct {
	ID        uuid.UUID  `json:"id"`
	BrokerID  int64      `json:"broker_id"`
	Type      string     `json:"type"`
	Link      string     `json:"link"`
	CreatedAt *time.Time `json:"created_at"`
}

type BrokerPhone struct {
	ID        uuid.UUID  `json:"id"`
	BrokerID  int64      `json:"broker_id"`
	Type      string     `json:"type"`
	Number    string     `json:"number"`
	IsPrimary *bool      `json:"is_primary"`
	CreatedAt *time.Time `json:"created_at"`
}

type BrokerProperty struct {
	ID              uuid.UUID  `json:"id"`
	BrokerID        int64      `json:"broker_id"`
	PropertyID      int64      `json:"property_id"`
	IsPrimaryBroker *bool      `json:"is_primary_broker"`
	CreatedAt       *time.Time `json:"created_at"`
}

type Property struct {
	// MLS number
	ID                int64      `json:"mls"`
	Title             string     `json:"title"`
	Category          string     `json:"category"`
	CivicNumber       *string    `json:"civic_number"`
	StreetName        *string    `json:"street_name"`
	ApartmentNumber   *string    `json:"apartment_number"`
	CityName          *string    `json:"city_name"`
	NeighbourhoodName *string    `json:"neighbourhood_name"`
	Price             string     `json:"price"`
	Description       *string    `json:"description"`
	BedroomNumber     *int32     `json:"bedroom_number"`
	RoomNumber        *int32     `json:"room_number"`
	BathroomNumber    *int32     `json:"bathroom_number"`
	Latitude          string     `json:"latitude"`
	Longitude         string     `json:"longitude"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

type PropertyExpense struct {
	ID           uuid.UUID  `json:"id"`
	PropertyID   int64      `json:"property_id"`
	Type         string     `json:"type"`
	AnnualPrice  string     `json:"annual_price"`
	MonthlyPrice string     `json:"monthly_price"`
	CreatedAt    *time.Time `json:"created_at"`
}

type PropertyFeature struct {
	ID         uuid.UUID  `json:"id"`
	PropertyID int64      `json:"property_id"`
	Title      string     `json:"title"`
	Value      string     `json:"value"`
	CreatedAt  *time.Time `json:"created_at"`
}

type PropertyPhoto struct {
	ID          uuid.UUID  `json:"id"`
	PropertyID  int64      `json:"property_id"`
	Link        string     `json:"link"`
	Description *string    `json:"description"`
	IsPrimary   *bool      `json:"is_primary"`
	CreatedAt   *time.Time `json:"created_at"`
}
