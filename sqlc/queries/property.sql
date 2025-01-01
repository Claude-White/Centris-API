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
SELECT property.* FROM property
INNER JOIN broker_property
ON property.id = broker_property.property_id
WHERE broker_property.broker_id = $1
LIMIT @number_of_items::int OFFSET @start_position::int;

-- name: GetAllAgencyProperties :many
SELECT property.* FROM property
INNER JOIN broker_property
ON property.id = broker_property.property_id
INNER JOIN broker
ON broker_property.broker_id = broker.id
WHERE LOWER(broker.agency_name) = $1
LIMIT @number_of_items::int OFFSET @start_position::int;

-- name: GetAllCategoryProperties :many
SELECT * FROM property
WHERE LOWER(property.category) = $1
LIMIT @number_of_items::int OFFSET @start_position::int;

-- name: GetAllCityProperties :many
SELECT * FROM property
WHERE LOWER(property.city_name) = $1
LIMIT @number_of_items::int OFFSET @start_position::int;

-- name: GetAllNeighbourhoodProperties :many
SELECT * FROM property
WHERE LOWER(property.neighbourhood_name) = $1
LIMIT @number_of_items::int OFFSET @start_position::int;

-- name: GetAllRadiusProperties :many
SELECT *
FROM property
WHERE (
    6371 * acos(
        cos(radians(@latitude::float32)) * cos(radians(property.latitude)) *
        cos(radians(property.longitude) - radians(@longitude::float32)) +
        sin(radians(@latitude::float32)) * sin(radians(latitude))
    )
) <= @radius::float32;

-- name: CreateProperty :one
INSERT INTO property (id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id;
