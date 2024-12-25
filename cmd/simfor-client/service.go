package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"pulsar_alice/internal/forward"
	"pulsar_alice/internal/meter/pulsar"
)

type service struct {
	config     config
	cancelFunc context.CancelFunc
	doneCh     chan struct{}
}

func new(config config) *service {
	return &service{
		config: config,
		doneCh: make(chan struct{}),
	}
}

func (s *service) Run(ctx context.Context, ready func()) error {
	defer close(s.doneCh)
	ctx, s.cancelFunc = context.WithCancel(ctx)

	ready()

	if err := ctx.Err(); err != nil {
		return err
	}

	coldValue, err := s.meterValue(ctx, s.config.Meters.ColdMeter)
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	hotValue, err := s.meterValue(ctx, s.config.Meters.HotMeter)
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	forwardClient := forward.New(s.config.Forward)
	session, err := forwardClient.StartSession(ctx)
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := session.SendMetrics(ctx, coldValue, hotValue); err != nil {
		return err
	}

	return nil
}

func (s *service) Shutdown(ctx context.Context) error {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	<-s.doneCh

	return nil
}

func (s *service) meterValue(ctx context.Context, meterNumber string) (float64, error) {
	logger := log.Ctx(ctx).With().Str("meter", meterNumber).Logger()

	meter := pulsar.New(pulsar.Config{
		Address: s.config.Meters.Address,
		Meter:   meterNumber,
		Timeout: s.config.Meters.Timeout,
	})

	if err := meter.Init(ctx); err != nil {
		logger.Error().Err(err).Msg("init meter")

		return 0, err
	}

	val := meter.Value(ctx)
	if val.Error != nil {
		logger.Error().Err(val.Error).Msg("read meter")

		return 0, val.Error
	}

	logger.Debug().Float64("value", val.Value).Msg("read success")

	return val.Value, nil
}
