-- name: CreateBrokerExternalLink :one
INSERT INTO broker_external_links (id, broker_id, type, link, created_at)
values (uuid_generate_v4(), $1, $2, $3, $4)
RETURNING id;