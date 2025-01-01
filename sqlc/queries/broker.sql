-- name: GetAllBrokers :many
SELECT * FROM broker
LIMIT $1 OFFSET $2;

-- name: GetBroker :one
SELECT * FROM broker 
WHERE broker.id = $1
LIMIT 1;
