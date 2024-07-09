package gateways

import (
	"context"

	"github.com/japhy-tech/backend-test/internal/domain/breeds"
)

// IDatastore
// More repositories could be added
type IDatastore interface {
	Breeds() breeds.Repository
	Close() error
	Reset(context.Context) error
}
