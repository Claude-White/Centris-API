-- name: GetAllProperties :many
SELECT * FROM property
LIMIT $1 OFFSET $2;
