package koala

import (
	"github.com/kelseyhightower/envconfig"
)

type EnvVars struct {
	TargetEnv string `default:"development"`
}

func LoadEnvVars() (EnvVars, error) {
	var vars EnvVars
	
	err := envconfig.Process("app", &vars)
	
	return vars, err
}