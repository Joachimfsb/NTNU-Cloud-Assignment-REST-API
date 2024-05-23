package util

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config is global access to configurations found in config.yaml
var Config config

// config stores all variables in the server configuration file.
// This allows global access to the configurations.
type config struct {
	Secrets struct {
		FirebaseKey string `yaml:"firebase_key"`
	} `yaml:"secrets"`
	Stubs struct {
		Database      bool `yaml:"database"`
		Currencies    bool `yaml:"currencies"`
		Weather       bool `yaml:"weather"`
		RestCountries bool `yaml:"rest_countries"`
	} `yaml:"stubs"`
}

// InitializeConfig must be run once. It reads variables in the config.yaml configurations file.
func InitializeConfig() error {
	err := cleanenv.ReadConfig("config.yaml", &Config)
	if err != nil {
		return fmt.Errorf("unable to initialize the configurations.\n%v", err)
	}

	return nil
}
