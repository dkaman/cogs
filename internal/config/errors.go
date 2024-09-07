package config

import (
	"errors"
	"fmt"
)

var (
	ErrConfigNotExist       = errors.New("provided config path does not exist")
)

type ConfigValidationError struct {
	message         string
	offendingParams []string
}

func (e *ConfigValidationError) Error() string {
	return fmt.Sprintf("%s: %v\n", e.message, e.offendingParams)
}
