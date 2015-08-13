package config

import (
	"io"
	"os"
	"net"
)

func Default() *Config {
	ifaces,_ := net.Interfaces()
	return &Config{
		NetworkInterfaces: ifaces,

		logger: os.Stderr,

	}
}

type Config struct {
		NetworkInterfaces []string

		logger io.Writer

		}

func (c *Config) Logger() io.Writer {
	return c.logger
}


