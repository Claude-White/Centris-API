// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: property_features.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPropertyFeature = `-- name: CreatePropertyFeature :one
INSERT INTO property_features (id, property_id, title, value, created_at)
values (uuid_generate_v4(), $1, $2, $3, $4)
RETURNING id
`

type CreatePropertyFeatureParams struct {
	PropertyID int64      `json:"property_id"`
	Title      string     `json:"title"`
	Value      string     `json:"value"`
	CreatedAt  *time.Time `json:"created_at"`
}

func (q *Queries) CreatePropertyFeature(ctx context.Context, arg CreatePropertyFeatureParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createPropertyFeature,
		arg.PropertyID,
		arg.Title,
		arg.Value,
		arg.CreatedAt,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
