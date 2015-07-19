package aurarath

import (
	"io"
	"net"
	"os"
)

type Config struct {
	LogOutput         string
	NetworkInterfaces []string

	logger io.Writer
}

func DefaultConfig() *Config {
	config := Config{
		LogOutput: "stderr",
		logger:    os.Stderr,
	}

	interfaces, err := net.Interfaces()

	if err != nil {
		return nil
	}

	var networkInterfaces []string

	for _, iface := range interfaces {
		append(networkInterfaces, iface.Name)
	}

	config.NetworkInterfaces = networkInterfaces

	return &config
}

func (c *Config) Logger() io.Writer {
	return c.logger
}
