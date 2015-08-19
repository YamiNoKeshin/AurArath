package config

import (
	"io"
	"net"
	"strings"
	"io/ioutil"
	"os"
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

		logger: ioutil.Discard,
		//logger: os.Stderr,

	}
}

func DefaultLocalhost() *Config {

	return &Config{
		NetworkInterfaces: []string{"145.108.172.104"},

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


