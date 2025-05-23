package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DemoMode     bool    `yaml:"demo_mode"`
	SingleMode   bool    `yaml:"single_mode"`
	DirectionOut bool    `yaml:"direction_out"`
	ApiKey       *string `yaml:"api_key,omitempty"`
	RelayPin     int     `yaml:"relay_pin"`
	HighMode     bool    `yaml:"high_mode"`
}

const configPath = "config.yaml"

var GlobalConfig *Config

func GetApiEndpoint() string {
	if GlobalConfig == nil {
		// this is an error but we're already checking for nil-ness before calling this
		return ""
	}
	if GlobalConfig.DemoMode {
		//return "http://localhost:3000/v2/turnstiles/"
		return "https://beta-backend-dev-kpe3ohblca-ew.a.run.app/v2/turnstiles/"
	}
	return "https://api.sendmoregetbeta.com/v2/turnstiles/"
}

func Init() error {
	// get file
	f, err := os.ReadFile(configPath)
	if err != nil {
		// if the error is that the file doesn't exist, create it
		if errors.Is(err, os.ErrNotExist) {
			c := Config{
				DemoMode:     true,
				SingleMode:   false,
				DirectionOut: false,
				ApiKey:       nil,
				RelayPin:     17,
				HighMode:     false,
			}
			GlobalConfig = &c
			err = StoreConfig(c)
			return err
		} else {
			return err
		}
	}

	// Create an empty Config to be the target of unmarshalling
	var c Config

	// Unmarshal our input YAML file into empty Config (var c)
	if err = yaml.Unmarshal(f, &c); err != nil {
		return err
	}
	GlobalConfig = &c
	return nil
}

func StoreConfig(newCfg Config) error {
	// store file
	bytes, err := yaml.Marshal(newCfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, bytes, os.ModeDevice)
	if err != nil {
		return err
	}

	GlobalConfig = &newCfg
	return nil
}
