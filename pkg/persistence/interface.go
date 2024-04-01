package persistence

import (
	"context"
	"errors"
)

var ClockNotFound = errors.New("persistence clock not found")

type Interface interface {
	Persist(ctx context.Context, in int64) error
	Load(ctx context.Context) (int64, error)
}
