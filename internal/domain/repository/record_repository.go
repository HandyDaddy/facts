package repository

import (
	"context"

	"github.com/HandyDaddy/facts/internal/domain/entities"
)

// FactRepository is an interface for persistence
type FactRepository interface {
	Save(ctx context.Context, fact *entities.Fact) error
	Get(ctx context.Context, id string) (*entities.Fact, error)
}
