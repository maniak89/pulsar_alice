package meterreader

import (
	"context"
	"time"

	"pulsar_alice/internal/meter"
	"pulsar_alice/internal/meter/pulsar"
	"pulsar_alice/internal/meter/wrap_logger"
	"pulsar_alice/internal/models/common"
	storageModels "pulsar_alice/internal/models/storage"
	"pulsar_alice/internal/notifier"
	"pulsar_alice/internal/storage"

	"github.com/rs/zerolog/log"
)

type client struct {
	provider   meter.ValueProvider
	meter      *storageModels.Meter
	lastUpdate time.Time
	lastValue  float64
}

type reader struct {
	config     Config
	clients    map[string][]*client
	storage    storage.Storage
	cancelFunc context.CancelFunc
	notifier   notifier.Notifier
}

func New(ctx context.Context, config Config, storage storage.Storage, notifier notifier.Notifier) (*reader, error) {
	logger := log.Ctx(ctx)

	meters, err := storage.Meters(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed fetch meters")

		return nil, err
	}

	result := reader{
		clients:  make(map[string][]*client, len(meters)),
		config:   config,
		storage:  storage,
		notifier: notifier,
	}

	for _, m := range meters {
		logger := logger.With().Str("serial_number", m.SerialNumber).Str("meter_id", m.ID).Logger()
		pCl := pulsar.New(pulsar.Config{
			Address: m.Address,
			Meter:   m.SerialNumber,
			Timeout: config.ReadTimeout,
		})

		cl := wrap_logger.New(pCl, m.ID, storage)

		if err := cl.Init(ctx); err != nil {
			logger.Error().Err(err).Msg("Failed init")

			return nil, err
		}

		result.clients[m.UserID] = append(result.clients[m.UserID], &client{
			provider: cl,
			meter:    m,
		})
	}

	return &result, nil
}

func (r *reader) Run(ctx context.Context, ready func()) error {
	logger := log.Ctx(ctx).With().Str("role", "checker").Logger()
	ctx = logger.WithContext(ctx)
	ctx, r.cancelFunc = context.WithCancel(ctx)
	defer r.cancelFunc()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(r.config.RefreshRate):
			if err := r.checkUpdates(ctx); err != nil {
				logger.Error().Err(err).Msg("check updates")

				return err
			}
		}
	}
}

func (r *reader) Shutdown(ctx context.Context) error {
	if r.cancelFunc != nil {
		r.cancelFunc()
	}

	return nil
}

func (r *reader) checkUpdates(ctx context.Context) error {
	for _, cls := range r.clients {
		var changed []*common.Meter
		for _, cl := range cls {
			previosVal := cl.lastValue
			if cl.readMeter(ctx, time.Duration(r.config.RefreshRate)).Value != previosVal {
				changed = append(changed, cl.makeMeter())
			}
		}

		if len(changed) > 0 {
			if err := r.notifier.NotifyMetersChanged(ctx, changed); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Failed notify")
			}
		}
	}

	return nil
}

func (r *reader) Meters(ctx context.Context, userID string) []*common.Meter {
	clients := r.clients[userID]
	if len(clients) == 0 {
		return nil
	}

	result := make([]*common.Meter, len(clients))

	for i, cl := range clients {
		result[i] = cl.readMeter(ctx, time.Duration(r.config.RefreshRate))
	}

	return result
}

func (cl *client) readMeter(ctx context.Context, refreshRate time.Duration) *common.Meter {
	if time.Since(cl.lastUpdate) < refreshRate {
		return cl.makeMeter()
	}

	val := cl.provider.Value(ctx)
	if val.Error != nil {
		return cl.makeMeter()
	}
	cl.lastValue = val.Value
	cl.lastUpdate = time.Now()

	return cl.makeMeter()
}

func (cl *client) makeMeter() *common.Meter {
	return &common.Meter{
		ID:           cl.meter.ID,
		UserID:       cl.meter.UserID,
		SerailNumber: cl.meter.SerialNumber,
		Name:         cl.meter.Name,
		Model:        "Light DU15",
		Manufacturer: "pulsar",
		Cold:         cl.meter.Cold,
		Updated:      cl.meter.UpdatedAt,
		Changed:      cl.lastUpdate,
		Value:        cl.lastValue,
	}
}
