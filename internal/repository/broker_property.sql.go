// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: broker_property.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createBrokerProperty = `-- name: CreateBrokerProperty :one
INSERT INTO broker_property (id, broker_id, property_id, created_at)
values (uuid_generate_v4(), $1, $2, $3)
RETURNING id
`

type CreateBrokerPropertyParams struct {
	BrokerID   int64      `json:"broker_id"`
	PropertyID int64      `json:"property_id"`
	CreatedAt  *time.Time `json:"created_at"`
}

func (q *Queries) CreateBrokerProperty(ctx context.Context, arg CreateBrokerPropertyParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createBrokerProperty, arg.BrokerID, arg.PropertyID, arg.CreatedAt)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
