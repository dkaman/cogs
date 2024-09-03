package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

const (
	CONFIG_DELIMITER      = "."
	ENVVAR_PREFIX         = "COGS"
	DefaultConfigDir      = ".cogs"
	DefaultConfigFileName = "config.json"
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

func (e *Config) Get(path string, o interface{}) error {
	return e.config.Unmarshal(path, o)
}

func (e *Config) PutString(path string, val string) error {
	return e.config.Set(path, val)
}

func (e *Config) JSON() string {
	b, _ := e.config.Marshal(json.Parser())
	return string(b)
}

func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, DefaultConfigDir, DefaultConfigFileName)
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
	return func(c *Config) error {
		p := env.Provider(ENVVAR_PREFIX, CONFIG_DELIMITER, func(key string) string {
			t := strings.ToLower(strings.TrimPrefix(key, ENVVAR_PREFIX+"_"))
			return strings.Replace(t, "_", ".", -1)
		})

		err := c.config.Load(p, nil)
		if err != nil {
			return fmt.Errorf("error reading config from env vars: %w", err)
		}

		return nil
	}
}

func WithFlags(fs *flag.FlagSet) configOption {
	return func(c *Config) error {
		p := posflag.ProviderWithFlag(fs, ".", nil, func(f *flag.Flag) (string, interface{}) {
			key := fmt.Sprintf("%s.%s", fs.Name(), f.Name)
			val := posflag.FlagVal(fs, f)
			return key, val
		})

		err := c.config.Load(p, nil)
		if err != nil {
			return fmt.Errorf("error reading cmd line flags: %w", err)
		}

		return nil
	}
}
