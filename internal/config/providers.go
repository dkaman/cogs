package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"

	flag "github.com/spf13/pflag"
)

const (
	EnvVarPrefix = "COGS"
)

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
		p := env.Provider(EnvVarPrefix, ConfigDelimiter, func(key string) string {
			t := strings.ToLower(strings.TrimPrefix(key, EnvVarPrefix+"_"))
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
