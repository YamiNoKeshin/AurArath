package beacon
import (
	"time"
	"io"
	"os"
)

type Config struct {

	ListenAddress string	//Address, Default 224.0.0.251
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