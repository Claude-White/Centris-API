// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: property.sql

package repository

import (
	"context"
)

const createProperty = `-- name: CreateProperty :one
INSERT INTO property (id, title, category, address, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING id
`

type CreatePropertyParams struct {
	ID                int64   `json:"mls"`
	Title             string  `json:"title"`
	Category          string  `json:"category"`
	Address           string  `json:"address"`
	CityName          string  `json:"city_name"`
	NeighbourhoodName *string `json:"neighbourhood_name"`
	Price             float32 `json:"price"`
	Description       *string `json:"description"`
	BedroomNumber     *int32  `json:"bedroom_number"`
	RoomNumber        *int32  `json:"room_number"`
	BathroomNumber    *int32  `json:"bathroom_number"`
	Longitude         float32 `json:"longitude"`
	Latitude          float32 `json:"latitude"`
}

func (q *Queries) CreateProperty(ctx context.Context, arg CreatePropertyParams) (int64, error) {
	row := q.db.QueryRow(ctx, createProperty,
		arg.ID,
		arg.Title,
		arg.Category,
		arg.Address,
		arg.CityName,
		arg.NeighbourhoodName,
		arg.Price,
		arg.Description,
		arg.BedroomNumber,
		arg.RoomNumber,
		arg.BathroomNumber,
		arg.Longitude,
		arg.Latitude,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getAllAgencyProperties = `-- name: GetAllAgencyProperties :many
SELECT property.id, property.title, property.category, property.address, property.city_name, property.neighbourhood_name, property.price, property.description, property.bedroom_number, property.room_number, property.bathroom_number, property.latitude, property.longitude, property.created_at, property.updated_at FROM property
INNER JOIN broker_property
ON property.id = broker_property.property_id
INNER JOIN broker
ON broker_property.broker_id = broker.id
WHERE LOWER(broker.agency_name) = $1
LIMIT $3::int OFFSET $2::int
`

type GetAllAgencyPropertiesParams struct {
	AgencyName    string `json:"agency_name"`
	StartPosition int32  `json:"start_position"`
	NumberOfItems int32  `json:"number_of_items"`
}

func (q *Queries) GetAllAgencyProperties(ctx context.Context, arg GetAllAgencyPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllAgencyProperties, arg.AgencyName, arg.StartPosition, arg.NumberOfItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Address,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllBrokerProperties = `-- name: GetAllBrokerProperties :many
SELECT property.id, property.title, property.category, property.address, property.city_name, property.neighbourhood_name, property.price, property.description, property.bedroom_number, property.room_number, property.bathroom_number, property.latitude, property.longitude, property.created_at, property.updated_at FROM property
INNER JOIN broker_property
ON property.id = broker_property.property_id
WHERE broker_property.broker_id = $1
LIMIT $3::int OFFSET $2::int
`

type GetAllBrokerPropertiesParams struct {
	BrokerID      int64 `json:"broker_id"`
	StartPosition int32 `json:"start_position"`
	NumberOfItems int32 `json:"number_of_items"`
}

func (q *Queries) GetAllBrokerProperties(ctx context.Context, arg GetAllBrokerPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllBrokerProperties, arg.BrokerID, arg.StartPosition, arg.NumberOfItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Address,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllCategoryProperties = `-- name: GetAllCategoryProperties :many
SELECT id, title, category, address, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, latitude, longitude, created_at, updated_at FROM property
WHERE LOWER(property.category) = $1
LIMIT $3::int OFFSET $2::int
`

type GetAllCategoryPropertiesParams struct {
	Category      string `json:"category"`
	StartPosition int32  `json:"start_position"`
	NumberOfItems int32  `json:"number_of_items"`
}

func (q *Queries) GetAllCategoryProperties(ctx context.Context, arg GetAllCategoryPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllCategoryProperties, arg.Category, arg.StartPosition, arg.NumberOfItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Address,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllCityProperties = `-- name: GetAllCityProperties :many
SELECT id, title, category, address, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, latitude, longitude, created_at, updated_at FROM property
WHERE LOWER(property.city_name) = $1
LIMIT $3::int OFFSET $2::int
`

type GetAllCityPropertiesParams struct {
	CityName      string `json:"city_name"`
	StartPosition int32  `json:"start_position"`
	NumberOfItems int32  `json:"number_of_items"`
}

func (q *Queries) GetAllCityProperties(ctx context.Context, arg GetAllCityPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllCityProperties, arg.CityName, arg.StartPosition, arg.NumberOfItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Address,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllNeighbourhoodProperties = `-- name: GetAllNeighbourhoodProperties :many
SELECT id, title, category, address, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, latitude, longitude, created_at, updated_at FROM property
WHERE LOWER(property.neighbourhood_name) = $1
LIMIT $3::int OFFSET $2::int
`

type GetAllNeighbourhoodPropertiesParams struct {
	NeighbourhoodName *string `json:"neighbourhood_name"`
	StartPosition     int32   `json:"start_position"`
	NumberOfItems     int32   `json:"number_of_items"`
}

func (q *Queries) GetAllNeighbourhoodProperties(ctx context.Context, arg GetAllNeighbourhoodPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllNeighbourhoodProperties, arg.NeighbourhoodName, arg.StartPosition, arg.NumberOfItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Address,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllProperties = `-- name: GetAllProperties :many
SELECT id, title, category, address, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, latitude, longitude, created_at, updated_at FROM property
LIMIT $1 OFFSET $2
`

type GetAllPropertiesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllProperties(ctx context.Context, arg GetAllPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllProperties, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Address,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllRadiusProperties = `-- name: GetAllRadiusProperties :many
SELECT id, title, category, address, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, latitude, longitude, created_at, updated_at
FROM property
WHERE (
    6371 * acos(
        cos(radians($1::float32)) * cos(radians(property.latitude)) *
        cos(radians(property.longitude) - radians($2::float32)) +
        sin(radians($1::float32)) * sin(radians(latitude))
    )
) <= $3::float32
`

type GetAllRadiusPropertiesParams struct {
	Latitude  interface{} `json:"latitude"`
	Longitude interface{} `json:"longitude"`
	Radius    interface{} `json:"radius"`
}

func (q *Queries) GetAllRadiusProperties(ctx context.Context, arg GetAllRadiusPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllRadiusProperties, arg.Latitude, arg.Longitude, arg.Radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.Address,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProperty = `-- name: GetProperty :one
SELECT id, title, category, address, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, latitude, longitude, created_at, updated_at FROM property 
WHERE property.id = $1
LIMIT 1
`

func (q *Queries) GetProperty(ctx context.Context, id int64) (Property, error) {
	row := q.db.QueryRow(ctx, getProperty, id)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Category,
		&i.Address,
		&i.CityName,
		&i.NeighbourhoodName,
		&i.Price,
		&i.Description,
		&i.BedroomNumber,
		&i.RoomNumber,
		&i.BathroomNumber,
		&i.Latitude,
		&i.Longitude,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPropertyByCoordinates = `-- name: GetPropertyByCoordinates :one
SELECT id, title, category, address, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, latitude, longitude, created_at, updated_at FROM property
WHERE property.longitude = $1 AND property.latitude = $2
LIMIT 1
`

type GetPropertyByCoordinatesParams struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

func (q *Queries) GetPropertyByCoordinates(ctx context.Context, arg GetPropertyByCoordinatesParams) (Property, error) {
	row := q.db.QueryRow(ctx, getPropertyByCoordinates, arg.Longitude, arg.Latitude)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Category,
		&i.Address,
		&i.CityName,
		&i.NeighbourhoodName,
		&i.Price,
		&i.Description,
		&i.BedroomNumber,
		&i.RoomNumber,
		&i.BathroomNumber,
		&i.Latitude,
		&i.Longitude,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
