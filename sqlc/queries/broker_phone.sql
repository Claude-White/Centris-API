-- name: CreateBrokerPhone :one
INSERT INTO broker_phone (id, broker_id, type, number, created_at)
values ($1, $2, $3, $4, $5)
RETURNING id;