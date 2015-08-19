package beacon
import (
	"time"
	"io"
	"os"
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