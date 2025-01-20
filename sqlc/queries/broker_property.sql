-- name: CreateBrokerProperty :one
INSERT INTO broker_property (id, broker_id, property_id, created_at)
values (uuid_generate_v4(), $1, $2, $3)
RETURNING id;

-- name: CreateAllBrokersProperties :copyfrom
INSERT INTO broker_property (id, broker_id, property_id, created_at)
values ($1, $2, $3, $4);