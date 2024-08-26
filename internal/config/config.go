package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	CONFIG_DELIMITER = "."
	ENVVAR_PREFIX    = "COGS"
)

var (
	ErrConfigNotExist = errors.New("provided config path does not exist")
)

type Config struct {
	config *koanf.Koanf
}

type configOption func(*Config) error

func New(opts ...configOption) (*Config, error) {
	c := &Config{
		config: koanf.New(CONFIG_DELIMITER),
	}

	for _, o := range opts {
		err := o(c)
		if err != nil {
			return nil, fmt.Errorf("error with Env option: %w", err)
		}
	}

	return c, nil
}

func WithJSONConfigFile(path string) configOption {
	return func(e *Config) error {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return ErrConfigNotExist
		}

		p := file.Provider(path)
		if err := e.config.Load(p, json.Parser()); err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}

		return nil
	}
}

func WithEnvVars() configOption {
	return func(e *Config) error {
		p := env.Provider(ENVVAR_PREFIX, CONFIG_DELIMITER, func(key string) string {
			t := strings.ToLower(strings.TrimPrefix(key, ENVVAR_PREFIX+"_"))
			return strings.Replace(t, "_", ".", -1)
		})

		err := e.config.Load(p, nil)
		if err != nil {
			return fmt.Errorf("error reading config from env vars: %w", err)
		}

		return nil
	}
}

func (e *Config) Get(path string, o interface{}) error {
	return e.config.Unmarshal(path, o)
}

func (e *Config) JSON() string {
	b, _ := e.config.Marshal(json.Parser())
	return string(b)
}
