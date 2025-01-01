// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createProperty = `-- name: CreateProperty :one
INSERT INTO property (id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id
`

type CreatePropertyParams struct {
	ID                int64
	Title             string
	Category          string
	CivicNumber       pgtype.Text
	StreetName        pgtype.Text
	ApartmentNumber   pgtype.Text
	CityName          pgtype.Text
	NeighbourhoodName pgtype.Text
	Price             pgtype.Numeric
	Description       pgtype.Text
	BedroomNumber     pgtype.Int4
	RoomNumber        pgtype.Int4
	BathroomNumber    pgtype.Int4
	Longitude         pgtype.Numeric
	Latitude          pgtype.Numeric
}

func (q *Queries) CreateProperty(ctx context.Context, arg CreatePropertyParams) (int64, error) {
	row := q.db.QueryRow(ctx, createProperty,
		arg.ID,
		arg.Title,
		arg.Category,
		arg.CivicNumber,
		arg.StreetName,
		arg.ApartmentNumber,
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
SELECT 
    property.id, 
    property.title, 
    property.category, 
    property.civic_number, 
    property.street_name, 
    property.apartment_number, 
    property.city_name, 
    property.neighbourhood_name, 
    property.price, 
    property.description, 
    property.bedroom_number, 
    property.room_number, 
    property.bathroom_number, 
    property.longitude, 
    property.latitude, 
    property.created_at, 
    property.updated_at 
FROM property
INNER JOIN broker_property
ON property.id = broker_property.property_id
INNER JOIN broker
ON broker_property.broker_id = broker.id
WHERE LOWER(broker.agency_name) = $1
LIMIT $2 OFFSET $3
`

type GetAllAgencyPropertiesParams struct {
	AgencyName string
	Limit      int32
	Offset     int32
}

func (q *Queries) GetAllAgencyProperties(ctx context.Context, arg GetAllAgencyPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllAgencyProperties, arg.AgencyName, arg.Limit, arg.Offset)
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
			&i.CivicNumber,
			&i.StreetName,
			&i.ApartmentNumber,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Longitude,
			&i.Latitude,
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
SELECT 
    property.id, 
    property.title, 
    property.category, 
    property.civic_number, 
    property.street_name, 
    property.apartment_number, 
    property.city_name, 
    property.neighbourhood_name, 
    property.price, 
    property.description, 
    property.bedroom_number, 
    property.room_number, 
    property.bathroom_number, 
    property.longitude, 
    property.latitude, 
    property.created_at, 
    property.updated_at 
FROM property
INNER JOIN broker_property
ON property.id = broker_property.property_id
WHERE broker_property.broker_id = $1
LIMIT $2 OFFSET $3
`

type GetAllBrokerPropertiesParams struct {
	BrokerID int64
	Limit    int32
	Offset   int32
}

func (q *Queries) GetAllBrokerProperties(ctx context.Context, arg GetAllBrokerPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllBrokerProperties, arg.BrokerID, arg.Limit, arg.Offset)
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
			&i.CivicNumber,
			&i.StreetName,
			&i.ApartmentNumber,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Longitude,
			&i.Latitude,
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
SELECT id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude, created_at, updated_at FROM property
WHERE LOWER(property.category) = $1
LIMIT $2 OFFSET $3
`

type GetAllCategoryPropertiesParams struct {
	Category string
	Limit    int32
	Offset   int32
}

func (q *Queries) GetAllCategoryProperties(ctx context.Context, arg GetAllCategoryPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllCategoryProperties, arg.Category, arg.Limit, arg.Offset)
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
			&i.CivicNumber,
			&i.StreetName,
			&i.ApartmentNumber,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Longitude,
			&i.Latitude,
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
SELECT id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude, created_at, updated_at FROM property
WHERE LOWER(property.city_name) = $1
LIMIT $2 OFFSET $3
`

type GetAllCityPropertiesParams struct {
	CityName pgtype.Text
	Limit    int32
	Offset   int32
}

func (q *Queries) GetAllCityProperties(ctx context.Context, arg GetAllCityPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllCityProperties, arg.CityName, arg.Limit, arg.Offset)
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
			&i.CivicNumber,
			&i.StreetName,
			&i.ApartmentNumber,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Longitude,
			&i.Latitude,
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
SELECT id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude, created_at, updated_at FROM property
WHERE LOWER(property.neighbourhood_name) = $1
LIMIT $2 OFFSET $3
`

type GetAllNeighbourhoodPropertiesParams struct {
	NeighbourhoodName pgtype.Text
	Limit             int32
	Offset            int32
}

func (q *Queries) GetAllNeighbourhoodProperties(ctx context.Context, arg GetAllNeighbourhoodPropertiesParams) ([]Property, error) {
	rows, err := q.db.Query(ctx, getAllNeighbourhoodProperties, arg.NeighbourhoodName, arg.Limit, arg.Offset)
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
			&i.CivicNumber,
			&i.StreetName,
			&i.ApartmentNumber,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Longitude,
			&i.Latitude,
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
SELECT id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude, created_at, updated_at FROM property
LIMIT $1 OFFSET $2
`

type GetAllPropertiesParams struct {
	Limit  int32
	Offset int32
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
			&i.CivicNumber,
			&i.StreetName,
			&i.ApartmentNumber,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Longitude,
			&i.Latitude,
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
SELECT id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude, created_at, updated_at
FROM property
WHERE (
    6371 * acos(
        cos(radians($1)) * cos(radians(property.latitude)) *
        cos(radians(property.longitude) - radians($2)) +
        sin(radians($1)) * sin(radians(latitude))
    )
) <= $3
`

type GetAllRadiusPropertiesParams struct {
	Latitude   float64
	Longitude float64
	Radius  pgtype.Numeric
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
			&i.CivicNumber,
			&i.StreetName,
			&i.ApartmentNumber,
			&i.CityName,
			&i.NeighbourhoodName,
			&i.Price,
			&i.Description,
			&i.BedroomNumber,
			&i.RoomNumber,
			&i.BathroomNumber,
			&i.Longitude,
			&i.Latitude,
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
SELECT id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude, created_at, updated_at FROM property 
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
		&i.CivicNumber,
		&i.StreetName,
		&i.ApartmentNumber,
		&i.CityName,
		&i.NeighbourhoodName,
		&i.Price,
		&i.Description,
		&i.BedroomNumber,
		&i.RoomNumber,
		&i.BathroomNumber,
		&i.Longitude,
		&i.Latitude,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPropertyByCoordinates = `-- name: GetPropertyByCoordinates :one
SELECT id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude, created_at, updated_at FROM property
WHERE property.longitude = $1 AND property.latitude = $2
LIMIT 1
`

type GetPropertyByCoordinatesParams struct {
	Longitude pgtype.Numeric
	Latitude  pgtype.Numeric
}

func (q *Queries) GetPropertyByCoordinates(ctx context.Context, arg GetPropertyByCoordinatesParams) (Property, error) {
	row := q.db.QueryRow(ctx, getPropertyByCoordinates, arg.Longitude, arg.Latitude)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Category,
		&i.CivicNumber,
		&i.StreetName,
		&i.ApartmentNumber,
		&i.CityName,
		&i.NeighbourhoodName,
		&i.Price,
		&i.Description,
		&i.BedroomNumber,
		&i.RoomNumber,
		&i.BathroomNumber,
		&i.Longitude,
		&i.Latitude,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
