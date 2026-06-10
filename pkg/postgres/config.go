package postgres

import (
	"errors"
	"fmt"
)

var ErrInvalidConfig = errors.New("postgres: invalid config")

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
	Options  map[string]string
}

func NewConfig() *Config {
	return &Config{Options: make(map[string]string)}
}

func (c *Config) WithHost(host string) *Config {
	c.Host = host
	return c
}

func (c *Config) WithPort(port int) *Config {
	c.Port = port
	return c
}

func (c *Config) WithUser(user string) *Config {
	c.User = user
	return c
}

func (c *Config) WithPassword(password string) *Config {
	c.Password = password
	return c
}

func (c *Config) WithDatabase(name string) *Config {
	c.Database = name
	return c
}

func (c *Config) WithSSLMode(mode string) *Config {
	c.SSLMode = mode
	return c
}

func (c *Config) WithOption(key, value string) *Config {
	if c.Options == nil {
		c.Options = make(map[string]string)
	}
	c.Options[key] = value
	return c
}

func (c *Config) Validate() error {
	if c == nil {
		return ErrInvalidConfig
	}
	if c.Host == "" {
		return fmt.Errorf("%w: host is required", ErrInvalidConfig)
	}
	if c.User == "" || c.Database == "" {
		return fmt.Errorf("%w: user and database are required", ErrInvalidConfig)
	}
	return nil
}

func (c *Config) DSN() string {
	host := c.Host
	if host == "" {
		host = "localhost"
	}
	port := c.Port
	if port == 0 {
		port = 5432
	}
	sslmode := c.SSLMode
	if sslmode == "" {
		if v, ok := c.Options["sslmode"]; ok && v != "" {
			sslmode = v
		} else {
			sslmode = "disable"
		}
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, c.User, c.Password, c.Database, sslmode,
	)
	for k, v := range c.Options {
		if k == "sslmode" {
			continue
		}
		dsn += fmt.Sprintf(" %s=%s", k, v)
	}
	return dsn
}
