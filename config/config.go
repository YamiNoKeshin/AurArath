package config

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func Default() *Config {
	ifaces, _ := net.Interfaces()
	var ifaceNames []string
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		if len(addrs) != 0 {
			addr := strings.Split(addrs[0].String(), "/")[0]
			ifaceNames = append(ifaceNames, addr)
		}
	}
	return &Config{
		NetworkInterfaces: ifaceNames,

		logger: ioutil.Discard,
		//logger: os.Stderr,

	}
}

func DefaultLocalhost() *Config {

	return &Config{
		NetworkInterfaces: []string{"127.0.0.1"},

		//logger: ioutil.Discard,
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
