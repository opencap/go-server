package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	host     string `json:"host"`
	port     uint16 `json:"port"`
	db       string `json:"db"`
	dbsource string `json:"dbsource"`
}

func Read(file *os.File) (*Config, error) {
	var conf Config
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func (conf *Config) Verify() error {
	return nil
}
