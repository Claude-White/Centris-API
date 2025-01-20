-- name: GetAllBrokers :many
SELECT * 
FROM broker
WHERE
    (broker.first_name IS NULL OR broker.first_name ILIKE '%' || coalesce(sqlc.narg('first_name'), first_name) || '%') AND
    (broker.middle_name IS NULL OR broker.middle_name ILIKE '%' || coalesce(sqlc.narg('middle_name'), middle_name) || '%') AND
    (broker.last_name IS NULL OR broker.last_name ILIKE '%' || coalesce(sqlc.narg('last_name'), last_name) || '%') AND
    (broker.agency_name IS NULL OR broker.agency_name ILIKE '%' || coalesce(sqlc.narg('agency'), agency_name) || '%') AND
    (broker.served_areas IS NULL OR broker.served_areas ILIKE '%' || coalesce(sqlc.narg('area'), served_areas) || '%') AND
    (broker.complementary_info IS NULL OR broker.complementary_info ILIKE '%' || coalesce(sqlc.narg('language'), complementary_info) || '%')
ORDER BY broker.first_name, broker.last_name
LIMIT @number_of_items::int OFFSET @start_position::int;

-- name: GetBroker :one
SELECT * FROM broker 
WHERE broker.id = @borker_id::int
LIMIT 1;

-- name: CreateBroker :one
INSERT INTO broker (id, first_name, middle_name, last_name, title, profile_photo, complementary_info, served_areas, presentation, corporation_name, agency_name, agency_address, agency_logo, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id;

-- name: DeleteAllBrokers :exec
DELETE FROM broker;
