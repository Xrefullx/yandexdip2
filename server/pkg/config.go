package pkg

import (
	"errors"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

// Config stores server config params.
type Config struct {
	ServerPort  string `env:"SEC_SERVER_ADDRESS" json:"server_address"  validate:"required,hostname_port"`
	DatabaseDSN string `env:"SEC_DATABASE_DSN" json:"database_dsn"  validate:"-"`
	TableName   string `env:"SEC_TABLE" json:"table"  validate:"-"`
	EnableHTTPS bool   `env:"SEC_ENABLE_HTTPS" json:"enable_https" envDefault:"false" validate:"-"`

	Debug   bool `env:"SEC_DEBUG" json:"-" envDefault:"false" validate:"-"`
	Migrate bool `env:"SEC_MIGRATE" json:"-" envDefault:"false" validate:"-"`
}

// Default config params.
const (
	defServerPort  = ":8080"
	defDatabaseDSN = "postgres://postgres:postgres@localhost:5432/test1?sslmode=disable"
	defTable       = "test1"
	defEnableHTTPS = false

	defDebug   = false
	defMigrate = false
)

// NewConfig inits new config.
// Reads flag params over default params, then redefines  with environment params.
func NewConfig() (*Config, error) {
	cfg := Config{}

	cfg.readFlagConfig()
	if err := cfg.readEnvConfig(); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации конфига: %w", err)
	}
	return &cfg, nil
}

// Validate validates config params.
func (c *Config) Validate() error {
	if len(c.DatabaseDSN) == 0 {
		return errors.New("master key is empty")
	}
	if len(c.ServerPort) == 0 {
		return errors.New("server-port is empty")
	}
	return nil
}

// readFlagConfig reads flag params over default params.
func (c *Config) readFlagConfig() {
	flagConfig := &Config{}

	flag.StringVar(&flagConfig.ServerPort, "p", defServerPort, "port for HTTP-server <:port>")
	flag.StringVar(&flagConfig.DatabaseDSN, "d", defDatabaseDSN, "database connection string")
	flag.BoolVar(&flagConfig.EnableHTTPS, "s", defEnableHTTPS, "enable https")

	flag.BoolVar(&flagConfig.Migrate, "migrate", defMigrate, "enable migrate database")
	flag.StringVar(&flagConfig.TableName, "t", defTable, "table name")

	flag.Parse()
	c.redefineConfig(flagConfig)
}

// redefineConfig redefines config with new config.
// if string not empty override
// if bool true override
func (c *Config) redefineConfig(nc *Config) {
	if nc.ServerPort != "" {
		c.ServerPort = nc.ServerPort
	}

	if nc.DatabaseDSN != "" {
		c.DatabaseDSN = nc.DatabaseDSN
	}

	if nc.TableName != "" {
		c.TableName = nc.TableName
	}

	if nc.EnableHTTPS {
		c.EnableHTTPS = nc.EnableHTTPS
	}

	if nc.Debug {
		c.Debug = nc.Debug
	}

	if nc.Migrate {
		c.Migrate = nc.Migrate
	}
}

// readEnvConfig redefines config params with environment params.
func (c *Config) readEnvConfig() error {
	envConfig := &Config{}

	if err := env.Parse(envConfig); err != nil {
		return fmt.Errorf("ошибка чтения переменных окружения:%w", err)
	}

	c.redefineConfig(envConfig)
	return nil
}
