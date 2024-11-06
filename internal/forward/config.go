package forward

import "time"

type Config struct {
	Address  string        `env:"FORWARD_ADDRESS,required"`
	Login    string        `env:"FORWARD_LOGIN,required"`
	Password string        `env:"FORWARD_PASSWORD,required"`
	Timout   time.Duration `env:"FORWARD_TIMEOUT,default=30s"`
}
