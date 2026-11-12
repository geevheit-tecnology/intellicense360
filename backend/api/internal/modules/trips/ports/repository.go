package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/domain"
)

type Query struct {
	Search    string
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	Filters   map[string]string
}

type Page[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type TripRepository interface {
	Create(ctx context.Context, item domain.Trip) (domain.Trip, error)
	FindByID(ctx context.Context, tenantID string, id string) (domain.Trip, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.Trip], error)
	Update(ctx context.Context, item domain.Trip) (domain.Trip, error)
	Delete(ctx context.Context, tenantID string, id string) error
	GetSummary(ctx context.Context, tenantID string) (domain.TripSummary, error)
}

