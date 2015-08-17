package config

import (
	"io"
	"os"
	"net"
	"strings"
)

func Default() *Config {
	ifaces,_ := net.Interfaces()
	var ifaceNames []string
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		if len(addrs) != 0 {
			addr := strings.Split(addrs[0].String(),"/")[0]
		ifaceNames = append(ifaceNames,addr)
		}
	}
	return &Config{
		NetworkInterfaces: ifaceNames,

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


