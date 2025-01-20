-- name: CreatePropertyFeature :one
INSERT INTO property_features (id, property_id, title, value, created_at)
values (uuid_generate_v4(), $1, $2, $3, $4)
RETURNING id;

-- name: CreateAllPropertiesFeatures :copyfrom
INSERT INTO property_features (id, property_id, title, value, created_at)
values ($1, $2, $3, $4, $5);