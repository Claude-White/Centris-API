-- name: GetAllProperties :many
SELECT * FROM property
LIMIT $1 OFFSET $2;

-- name: CreateProperty :one
INSERT INTO property (id, title, category, civic_number, street_name, apartment_number, city_name, neighbourhood_name, price, description, bedroom_number, room_number, bathroom_number, longitude, latitude)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id;
