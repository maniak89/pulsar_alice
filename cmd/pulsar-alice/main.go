package main

import (
	"context"
	"os"
	"os/signal"
	"pulsar_alice/internal/log"
	"pulsar_alice/internal/meterreader"
	"pulsar_alice/internal/services"
	"pulsar_alice/internal/services/rest"
	"pulsar_alice/internal/storage/sql"
	"syscall"

	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/run"
	zerolog "github.com/rs/zerolog/log"
)

type config struct {
	Logger      log.Config
	Rest        rest.Config
	Storage     sql.Config
	MeterReader meterreader.Config
}

const signalChLen = 10

func main() {
	var cfg config
	if err := envdecode.StrictDecode(&cfg); err != nil {
		zerolog.Fatal().Err(err).Msg("Cannot decode config envs")
	}

	logger, err := log.New(cfg.Logger)
	if err != nil {
		zerolog.Fatal().Err(err).Msg("Cannot init logger")
	}

	ctx, cancel := context.WithCancel(logger.WithContext(context.Background()))

	g := &run.Group{}
	{
		stop := make(chan os.Signal, signalChLen)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		g.Add(func() error {
			<-stop
			return nil
		}, func(error) {
			signal.Stop(stop)
			cancel()
			close(stop)
		})
	}

	orderRunner := services.OrderRunner{}

	storage := sql.New(cfg.Storage)
	if err := storage.Connect(ctx); err != nil {
		logger.Panic().Err(err).Msg("Failed connect to db")
	}
	defer func() {
		if err := storage.Disconnect(ctx); err != nil {
			logger.Error().Err(err).Msg("Failed disconnect from db")
		}
	}()

	meterReader, err := meterreader.New(ctx, cfg.MeterReader, storage)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed create meter reader")
	}

	restService, err := rest.New(ctx, cfg.Rest, logger.With().Str("role", "rest").Logger(), meterReader)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed create rest service")
	}
	if err := orderRunner.SetupService(ctx, restService, "rest", g); err != nil {
		logger.Fatal().Err(err).Msg("Failed setup rest service")
	}

	logger.Info().Msg("Running the service...")
	if err := g.Run(); err != nil {
		logger.Fatal().Err(err).Msg("The service has been stopped with error")
	}
	logger.Info().Msg("The service is stopped")
}
