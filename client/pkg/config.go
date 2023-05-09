package pkg

import (
	"errors"
	"flag"
	"fmt"
)

// Config stores server config params.
type Config struct {
	MasterKey         string
	SyncTimeoutSec    int
	RequestsPerMinute int
	ServerURL         string
	StorageFile       string
}

// Default config params.
const (
	defMasterKey         = "mytestmasterkey"
	defSyncTimeout       = 2
	defRequestsPerMinute = 100
	defServerURL         = "https://localhost:8085"
	defStorageFile       = "storage.db"
)

// NewConfig inits new config.
// Reads flag params over default params, then redefines  with environment params.
func NewConfig() (*Config, error) {
	cfg := Config{}

	cfg.readFlagConfig()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации конфига: %w", err)
	}
	return &cfg, nil
}

// Validate validates config params.
func (c *Config) Validate() error {
	if len(c.MasterKey) == 0 {
		return errors.New("master key is empty")
	}
	if len(c.ServerURL) == 0 {
		return errors.New("server url is empty")
	}
	if len(c.StorageFile) == 0 {
		return errors.New("storage filename is empty")
	}
	if c.SyncTimeoutSec == 0 {
		return errors.New("sync timeout is 0")
	}
	if c.RequestsPerMinute == 0 {
		return errors.New("requests per minute is 0")
	}

	return nil
}

// readFlagConfig reads flag params over default params.
func (c *Config) readFlagConfig() {
	flagConfig := &Config{}
	flag.StringVar(&flagConfig.MasterKey, "m", defMasterKey, "master key")
	flag.IntVar(&flagConfig.SyncTimeoutSec, "t", defSyncTimeout, "sync timeout in seconds")
	flag.IntVar(&flagConfig.RequestsPerMinute, "r", defRequestsPerMinute, "sync action requests per minute")
	flag.StringVar(&flagConfig.ServerURL, "s", defServerURL, "server address http(s)://<address>:<port>")
	flag.StringVar(&flagConfig.StorageFile, "db", defStorageFile, "storage filename")

	flag.Parse()
	c.redefineConfig(flagConfig)
}

// redefineConfig redefines config with new config.
// if string not empty override
// if bool true override
func (c *Config) redefineConfig(nc *Config) {
	if nc.MasterKey != "" {
		c.MasterKey = nc.MasterKey
	}
	if nc.StorageFile != "" {
		c.StorageFile = nc.StorageFile
	}
	if nc.ServerURL != "" {
		c.ServerURL = nc.ServerURL
	}
	if nc.SyncTimeoutSec != 0 {
		c.SyncTimeoutSec = nc.SyncTimeoutSec
	}
	if nc.RequestsPerMinute != 0 {
		c.RequestsPerMinute = nc.RequestsPerMinute
	}
}
