package meter

import (
	"context"
)

type ValueProvider interface {
	Init(ctx context.Context) error
	Value(ctx context.Context) Value
}

type Value struct {
	Address string
	Value   float64
	Error   error
}
