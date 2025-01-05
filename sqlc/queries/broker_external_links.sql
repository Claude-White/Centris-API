-- name: CreateBrokerExternalLink :one
INSERT INTO broker_external_links (id, broker_id, type, link, created_at)
values ($1, $2, $3, $4, $5)
RETURNING id;