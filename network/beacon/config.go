package beacon

import (
	"io"
	"os"
	"time"
)

type Config struct {
	PingAddresses []string
	//
	Port int

	PingInterval time.Duration

	Logger io.Writer
}

func (c *Config) init() {
	if c.Logger == nil {
		c.Logger = os.Stderr
	}
}
