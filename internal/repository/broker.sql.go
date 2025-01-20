// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: broker.sql

package repository

import (
	"context"
	"time"
)

const createBroker = `-- name: CreateBroker :one
INSERT INTO broker (id, first_name, middle_name, last_name, title, profile_photo, complementary_info, served_areas, presentation, corporation_name, agency_name, agency_address, agency_logo, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id
`

type CreateBrokerParams struct {
	ID                int64      `json:"id"`
	FirstName         string     `json:"first_name"`
	MiddleName        *string    `json:"middle_name"`
	LastName          string     `json:"last_name"`
	Title             string     `json:"title"`
	ProfilePhoto      *string    `json:"profile_photo"`
	ComplementaryInfo *string    `json:"complementary_info"`
	ServedAreas       *string    `json:"served_areas"`
	Presentation      *string    `json:"presentation"`
	CorporationName   *string    `json:"corporation_name"`
	AgencyName        string     `json:"agency_name"`
	AgencyAddress     string     `json:"agency_address"`
	AgencyLogo        *string    `json:"agency_logo"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

func (q *Queries) CreateBroker(ctx context.Context, arg CreateBrokerParams) (int64, error) {
	row := q.db.QueryRow(ctx, createBroker,
		arg.ID,
		arg.FirstName,
		arg.MiddleName,
		arg.LastName,
		arg.Title,
		arg.ProfilePhoto,
		arg.ComplementaryInfo,
		arg.ServedAreas,
		arg.Presentation,
		arg.CorporationName,
		arg.AgencyName,
		arg.AgencyAddress,
		arg.AgencyLogo,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteAllBrokers = `-- name: DeleteAllBrokers :exec
DELETE FROM broker
`

func (q *Queries) DeleteAllBrokers(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deleteAllBrokers)
	return err
}

const getAllBrokers = `-- name: GetAllBrokers :many
SELECT id, first_name, middle_name, last_name, title, profile_photo, complementary_info, served_areas, presentation, corporation_name, agency_name, agency_address, agency_logo, created_at, updated_at 
FROM broker
WHERE ($1::text IS NULL OR name = $1::text)
    AND ($2::text IS NULL OR agency = $2::text)
    AND ($3::text IS NULL OR area = $3::text)
    AND ($4::text IS NULL OR language = $4::text)
ORDER BY broker.first_name, broker.last_name
LIMIT $6::int OFFSET $5::int
`

type GetAllBrokersParams struct {
	BrokerName    string `json:"broker_name"`
	Agency        string `json:"agency"`
	Area          string `json:"area"`
	Language      string `json:"language"`
	StartPosition int32  `json:"start_position"`
	NumberOfItems int32  `json:"number_of_items"`
}

func (q *Queries) GetAllBrokers(ctx context.Context, arg GetAllBrokersParams) ([]Broker, error) {
	rows, err := q.db.Query(ctx, getAllBrokers,
		arg.BrokerName,
		arg.Agency,
		arg.Area,
		arg.Language,
		arg.StartPosition,
		arg.NumberOfItems,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Broker
	for rows.Next() {
		var i Broker
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.MiddleName,
			&i.LastName,
			&i.Title,
			&i.ProfilePhoto,
			&i.ComplementaryInfo,
			&i.ServedAreas,
			&i.Presentation,
			&i.CorporationName,
			&i.AgencyName,
			&i.AgencyAddress,
			&i.AgencyLogo,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBroker = `-- name: GetBroker :one
SELECT id, first_name, middle_name, last_name, title, profile_photo, complementary_info, served_areas, presentation, corporation_name, agency_name, agency_address, agency_logo, created_at, updated_at FROM broker 
WHERE broker.id = $1::int
LIMIT 1
`

func (q *Queries) GetBroker(ctx context.Context, borkerID int32) (Broker, error) {
	row := q.db.QueryRow(ctx, getBroker, borkerID)
	var i Broker
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.MiddleName,
		&i.LastName,
		&i.Title,
		&i.ProfilePhoto,
		&i.ComplementaryInfo,
		&i.ServedAreas,
		&i.Presentation,
		&i.CorporationName,
		&i.AgencyName,
		&i.AgencyAddress,
		&i.AgencyLogo,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
