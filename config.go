package aurarath

import (
	"io"
	"os"
)

type Config struct {
	LogOutput         string
	NetworkInterfaces []string

	logger io.Writer
}

func DefaultConfig() *Config {
	return &Config{
		LogOutput:         "stderr",
		NetworkInterfaces: []string{"eth0"},

		logger: os.Stderr,
	}
}

func (c *Config) Logger() io.Writer {
	return c.logger
}
