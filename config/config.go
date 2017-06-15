package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// ConfigFilename settings filename
var ConfigFilename string

func init() {
	// Gets the config filename from env var
	ConfigFilename = os.Getenv("CONFIG_FILENAME")

	if len(ConfigFilename) == 0 {
		ConfigFilename = "app"
	}

	// Name of config file (without extension)
	viper.SetConfigName(ConfigFilename)
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

// Config represents the application settings
type Config struct {
	Debug bool

	Session struct {
		Secret string
	}

	Jwt struct {
		Secret string
		Exp    int
	}

	DB struct {
		Driver string
		DSN    string
	}
}

// Viper returns the viper instance
func (c Config) Viper() *viper.Viper {
	return viper.GetViper()
}

// LoadConfig loads the application settings
func LoadConfig() (Config, error) {
	var c Config

	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}
