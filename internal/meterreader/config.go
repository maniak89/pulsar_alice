package meterreader

import "time"

type Config struct {
	RefreshRate time.Duration `env:"METERREADER_REFRESH_RATE,default=5m"`
	ReadTimeout time.Duration `env:"METERREADER_READ_TIMEOUT,default=5s"`
}
