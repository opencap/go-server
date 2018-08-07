package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	prefix = ""
	indent = "    "
)

type Config interface {
	Hostname() string
	Port() uint16
	Debug() bool
	DatabaseType() string
	DatabaseDataSource() string
}

type config struct {
	Config    `json:"-"`
	JHost     string         `json:"host"`
	JPort     uint16         `json:"port"`
	JDebug    bool           `json:"debug"`
	JDatabase databaseConfig `json:"database,omitempty"`
}

type databaseConfig struct {
	JType       string `json:"type"`
	JDataSource string `json:"dataSource"`
}

func (conf *config) Hostname() string {
	return conf.JHost
}

func (conf *config) Port() uint16 {
	return conf.JPort
}

func (conf *config) Debug() bool {
	return conf.JDebug
}

func (conf *config) DatabaseType() string {
	return conf.JDatabase.JType
}

func (conf *config) DatabaseDataSource() string {
	return conf.JDatabase.JDataSource
}

var defaultConfig = config{
	JHost:  "",
	JPort:  41145,
	JDebug: false,
	JDatabase: databaseConfig{
		JType:       "sqlite",
		JDataSource: "file:opencap.db",
	},
}

func LoadConfig(path string) (Config, error) {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("fileinfo for %s failed: %v", path, err)
	}

	var config config
	if os.IsNotExist(err) {
		config = defaultConfig
	} else {
		file, err := ioutil.ReadFile(path)
		if err != nil && err != os.ErrNotExist {
			return nil, fmt.Errorf("reading %s failed: %v", path, err)
		}

		if err := json.Unmarshal(file, &config); err != nil {
			return nil, fmt.Errorf("parsing failed: %v", err)
		}
	}

	file, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("formatting failed: %v", err)
	}

	if err := ioutil.WriteFile(path, file, 0644); err != nil {
		return nil, fmt.Errorf("writing to %s failed: %v", path, err)
	}
	return &config, nil
}
