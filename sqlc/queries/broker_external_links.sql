-- name: GetAllBrokerLinksByBrokerId :many
SELECT * FROM broker_external_links
WHERE broker_external_links.broker_id = $1;

-- name: CreateBrokerExternalLink :one
INSERT INTO broker_external_links (id, broker_id, type, link, created_at)
values (uuid_generate_v4(), $1, $2, $3, $4)
RETURNING id;

-- name: CreateAllBrokerExternalLink :copyfrom
INSERT INTO broker_external_links (id, broker_id, type, link, created_at)
values ($1, $2, $3, $4, $5);
