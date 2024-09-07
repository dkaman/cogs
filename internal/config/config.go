package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/v2"
)

const (
	ConfigDelimiter      = "."
	DefaultConfigDir      = ".cogs"
	DefaultConfigFileName = "config.json"
)

func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, DefaultConfigDir, DefaultConfigFileName)
}

type Config struct {
	config *koanf.Koanf
}

type configOption func(*Config) error

func New(opts ...configOption) (*Config, error) {
	c := &Config{
		config: koanf.New(ConfigDelimiter),
	}

	for _, o := range opts {
		err := o(c)
		if err != nil {
			return nil, fmt.Errorf("error with Env option: %w", err)
		}
	}

	return c, nil
}

func (c *Config) Get(path string, o interface{}) error {
	return c.config.Unmarshal(path, o)
}

func (c *Config) PutString(path string, val string) error {
	return c.config.Set(path, val)
}

func (c *Config) JSON() string {
	b, _ := c.config.Marshal(json.Parser())
	return string(b)
}

func (c *Config) Validate(required []string) (err error) {
	var invalid []string

	for _, r := range required {
		if val := c.config.String(r); val == "" {
			invalid = append(invalid, r)
		}
	}

	if len(invalid) != 0 {
		err = &ConfigValidationError{
			message: "required parameter not set",
			offendingParams: invalid,
		}
	}

	return
}

func (c *Config) Merge(fs ...configOption) (err error) {
	original := *c.config

	for _, f := range fs {
		err = f(c)
		if err != nil {
			c.config = &original
			return
		}
	}

	return
}
