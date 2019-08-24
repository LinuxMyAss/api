package config

import (
	"encoding/json"
	"os"
)

// Config .
type Config struct {
	Environment string `json:"environment"`

	Backend struct {
		Secret string `json:"secret"`

		MongoDB struct {
			URI      string `json:"uri"`
			Database string `json:"database"`
		} `json:"mongodb"`

		Redis struct {
			URI      string `json:"uri"`
			Password string `json:"password"`
			Database int    `json:"database"`
		} `json:"redis"`
	} `json:"backend"`

	HTTP struct {
		Address string `json:"address"`
	} `json:"http"`

	SMTP struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		From     string `json:"from"`

		Register struct {
			From    string `json:"from"`
			Subject string `json:"subject"`
		} `json:"register"`
	} `json:"smtp"`
}

var config *Config

// Get returns the loaded config object.
func Get() *Config {
	return config
}

// IsProduction is self explanatory..
func IsProduction() bool {
	return config.Environment == "production"
}

// Load loads the configuration from the disk.
func Load() error {
	file, err := os.Open(".env/config.json")
	defer file.Close()

	if err != nil {
		return err
	}

	parser := json.NewDecoder(file)
	err = parser.Decode(&config)
	if err != nil {
		return err
	}

	if len(config.Backend.MongoDB.URI) < 1 {
		config.Backend.MongoDB.URI = "mongodb://127.0.0.1:27017"
	}

	if len(config.Backend.MongoDB.Database) < 1 {
		config.Backend.MongoDB.Database = "ikuta"
	}

	if len(config.HTTP.Address) < 1 {
		config.HTTP.Address = ":7136"
	}

	return nil
}
