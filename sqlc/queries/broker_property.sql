-- name: CreateBrokerProperty :one
INSERT INTO broker_property (id, broker_id, property_id, created_at)
values (uuid_generate_v4(), $1, $2, $3)
RETURNING id;