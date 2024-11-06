package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/run"
	zerolog "github.com/rs/zerolog/log"

	"pulsar_alice/internal/forward"
	"pulsar_alice/internal/log"
	"pulsar_alice/internal/services"
)

type config struct {
	Logger log.Config
	Meters struct {
		Address   string        `env:"METERS_ADDRESS,required"`
		ColdMeter string        `env:"METERS_COLD,required"`
		HotMeter  string        `env:"METERS_HOT,required"`
		Timeout   time.Duration `env:"METERS_TIMEOUT,default=5s"`
	}
	Forward forward.Config
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

	if err := orderRunner.SetupService(ctx, new(cfg), "process", g); err != nil {
		logger.Fatal().Err(err).Msg("Failed setup process service")
	}

	logger.Info().Msg("Running the service...")
	if err := g.Run(); err != nil {
		logger.Fatal().Err(err).Msg("The service has been stopped with error")
	}
	logger.Info().Msg("The service is stopped")
}
