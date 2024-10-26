package pulsar

import "time"

type Config struct {
	Address string
	Meter   string
	Timeout time.Duration
}
