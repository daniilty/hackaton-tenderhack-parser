package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr    string `json:"httpAddr"`
	PostgresDSN string `json:"postgresDsn"`
}

func getDefault() *Config {
	return &Config{}
}

func Read(name string) (*Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	c := getDefault()

	err = dec.Decode(c)
	if err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return c, nil
}
