package koala

import (
	"github.com/kelseyhightower/envconfig"
)

// EnvVars represents project environment variables
type EnvVars struct {
	TargetEnv string `default:"development"`
}

// It loads environment variables
func LoadEnvVars() (EnvVars, error) {
	var vars EnvVars
	
	err := envconfig.Process("app", &vars)
	
	return vars, err
}