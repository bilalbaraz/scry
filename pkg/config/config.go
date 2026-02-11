package config

import (
	"errors"
	"os"
)

type Config struct {
	Path  string
	Raw   string
	Found bool
}

func Load(path string, required bool) (Config, error) {
	cfgPath := path
	if cfgPath == "" {
		cfgPath = ".scry.yml"
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && !required {
			return Config{Path: cfgPath, Found: false}, nil
		}
		return Config{}, err
	}

	return Config{Path: cfgPath, Raw: string(data), Found: true}, nil
}
