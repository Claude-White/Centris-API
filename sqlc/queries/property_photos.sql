-- name: CreatePropertyPhoto :one
INSERT INTO property_photo (id, property_id, link, description, created_at)
values (uuid_generate_v4(), $1, $2, $3, $4)
RETURNING id;

-- name: CreateAllPropertiesPhotos :copyfrom
INSERT INTO property_photo (id, property_id, link, description, created_at)
values ($1, $2, $3, $4, $5);