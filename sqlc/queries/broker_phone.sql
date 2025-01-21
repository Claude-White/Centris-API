-- name: GetAllBrokerPhonesByBrokerId :many
SELECT * FROM broker_phone
WHERE broker_phone.broker_id = $1;

-- name: CreateBrokerPhone :one
INSERT INTO broker_phone (id, broker_id, type, number, created_at)
values (uuid_generate_v4(), $1, $2, $3, $4)
RETURNING id;

-- name: CreateAllBrokerPhone :copyfrom
INSERT INTO broker_phone (id, broker_id, type, number, created_at)
values ($1, $2, $3, $4, $5);
