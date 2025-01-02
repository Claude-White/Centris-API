-- name: GetAllBrokers :many
SELECT * FROM broker
LIMIT @number_of_items::int OFFSET @start_position::int;

-- name: GetBroker :one
SELECT * FROM broker 
WHERE broker.id = $1
LIMIT 1;