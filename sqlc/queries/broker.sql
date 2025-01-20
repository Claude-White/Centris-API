-- name: GetAllBrokers :many
SELECT * 
FROM broker
WHERE (@broker_name::text IS NULL OR name = @broker_name::text)
    AND (@agency::text IS NULL OR agency = @agency::text)
    AND (@area::text IS NULL OR area = @area::text)
    AND (@language::text IS NULL OR language = @language::text)
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
