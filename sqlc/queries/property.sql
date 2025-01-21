-- name: GetAllProperties :many
SELECT * FROM property
LIMIT @number_of_items::int OFFSET @start_position::int;

-- name: GetProperty :one
SELECT * FROM property 
WHERE property.id = @mls::int
LIMIT 1;

-- name: GetPropertyByCoordinates :many
SELECT * FROM property
WHERE
    ABS(property.longitude - $1) < 0.000001 
    AND ABS(property.latitude - $2) < 0.000001
LIMIT @number_of_items::int OFFSET @start_position::int;

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

-- name: GetAllRadiusProperties :many
SELECT *
FROM property
WHERE (
    6371 * acos(
        cos(radians(@latitude::numeric)) * cos(radians(property.latitude)) *
        cos(radians(property.longitude) - radians(@longitude::numeric)) +
        sin(radians(@latitude::numeric)) * sin(radians(property.latitude))
    )
) <= @radius::numeric
LIMIT @number_of_items::int OFFSET @start_position::int;

-- name: CreateProperty :one
INSERT INTO property (id, title, category, address, city_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id;

-- name: CreateAllProperties :copyfrom
INSERT INTO property (id, title, category, address, city_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: DeleteAllProperties :exec
DELETE FROM property;
