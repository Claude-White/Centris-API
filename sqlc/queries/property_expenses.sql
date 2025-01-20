-- name: CreatePropertyExpenses :one
INSERT INTO property_expenses (id, property_id, type, annual_price, monthly_price, created_at)
values (uuid_generate_v4(), $1, $2, $3, $4, $5)
RETURNING id;

-- name: CreateAllPropertiesExpenses :copyfrom
INSERT INTO property_expenses (id, property_id, type, annual_price, monthly_price, created_at)
values ($1, $2, $3, $4, $5, $6);