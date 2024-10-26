package storage

import (
	"context"
	"errors"

	"pulsar_alice/internal/models/storage"
)

var ErrInvalidState = errors.New("invalid state")

type Storage interface {
	Meters(ctx context.Context) ([]*storage.Meter, error)
	Log(ctx context.Context, meterID string, level storage.LogLevel, msg string)
}
