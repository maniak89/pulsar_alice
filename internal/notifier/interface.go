package notifier

import (
	"context"
	"pulsar_alice/internal/models/common"
)

type Notifier interface {
	NotifyMetersChanged(ctx context.Context, meter []*common.Meter) error
}
