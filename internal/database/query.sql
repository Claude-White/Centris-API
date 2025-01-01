-- name: GetAllProperties :many
SELECT * FROM property
LIMIT $1 OFFSET $2;

-- name: GetProperty :one
SELECT * FROM property 
WHERE property.id = $1
LIMIT 1;

-- name: GetPropertyByCoordinates :one
SELECT * FROM property
WHERE property.longitude = $1 AND property.latitude = $2
LIMIT 1;

-- name: GetAllBrokerProperties :many
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
LIMIT $2 OFFSET $3;

-- name: GetAllAgencyProperties :many
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
LIMIT $2 OFFSET $3;

-- name: GetAllCategoryProperties :many
SELECT * FROM property
WHERE LOWER(property.category) = $1
LIMIT $2 OFFSET $3;

-- name: GetAllCityProperties :many
SELECT * FROM property
WHERE LOWER(property.city_name) = $1
LIMIT $2 OFFSET $3;

-- name: GetAllNeighbourhoodProperties :many
SELECT * FROM property
WHERE LOWER(property.neighbourhood_name) = $1
LIMIT $2 OFFSET $3;

-- name: GetAllRadiusProperties :many
SELECT *
FROM property
WHERE (
    6371 * acos(
        cos(radians($1)) * cos(radians(property.latitude)) *
        cos(radians(property.longitude) - radians($2)) +
        sin(radians($1)) * sin(radians(latitude))
    )
) <= $3;

-- name: CreateProperty :one
INSERT INTO property (id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id;