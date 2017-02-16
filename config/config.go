package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"gopkg.in/yaml.v2"
)

// ConfigFilename settings filename
var ConfigFilename string

func init() {
	// Gets the config filename from env var
	ConfigFilename = os.Getenv("CONFIG_FILENAME")

	if len(ConfigFilename) == 0 {
		ConfigFilename = "app.yml"
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

var tplFuncMap = template.FuncMap{
	"env": getEnvVar,
}

// getEnvVar gets an environment variable or the default value
func getEnvVar(v string, d string) string {
	r := os.Getenv(v)

	if r != "" {
		return r
	}

	return d
}

// ParseSettingsFile parse the settings file
// The parse can override some values as environment variables
func ParseSettingsFile(text string) ([]byte, error) {
	template := template.New("koalasettingstpl")

	template, err := template.Funcs(tplFuncMap).Parse(text)

	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	template.Execute(&buffer, nil)

	return buffer.Bytes(), nil
}

// LoadConfigFromBytes loads the application settings from bytes
// It allows that binary data can be used
func LoadConfigFromBytes(b []byte) (Config, error) {
	var c Config

	if err := yaml.Unmarshal(b, &c); err != nil {
		return c, err
	}

	return c, nil
}

// LoadConfig loads the application settings
func LoadConfig() (c Config, err error) {
	bytes, err := ReadConfigFile()

	if err != nil {
		return c, err
	}

	bytes, err = ParseSettingsFile(string(bytes))

	if err != nil {
		return c, err
	}

	return LoadConfigFromBytes(bytes)
}

// ReadConfigFile reads the settings file per convention
// The file should be stored inside config folder in the project root
func ReadConfigFile() ([]byte, error) {
	return ioutil.ReadFile(path.Join("config", ConfigFilename))
}
