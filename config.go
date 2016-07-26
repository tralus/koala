package koala

import (
  	"io/ioutil"
  	"path"

  	"github.com/kardianos/osext"
  	"github.com/olebedev/config"
)

// It loads the config file from current project folder if env is development or exe folder
func LoadConfig(targetEnv string) (*config.Config, error) {
	var c *config.Config
	var err error
	
	if (targetEnv == "development") {
		// file path relative to project folder
		c, err = LoadConfigFromFolder("")
		
		if err != nil {
			return nil, err
		}
	} else {
		binFolder, err := osext.ExecutableFolder()
	  
		if err != nil {
			return nil, err
		}
		
		c, err = LoadConfigFromFolder(binFolder)
		
		if err != nil {
			return nil, err
		}
	}
	
	return c.Get(targetEnv)
}

// It loads the settings.yml file from folder parameter
func LoadConfigFromFolder(folder string) (*config.Config, error) {  
  	filename := path.Join(folder, "settings.yaml")
  
  	file, err := ioutil.ReadFile(filename)

  	if err != nil {
    	return nil, err
  	}
  
  	content := string(file)
  
  	if err != nil {
    	return nil, err
  	}
  
  	return config.ParseYaml(content)
}

// JwtConfig represents JWT Token information
type JwtConfig struct {
	Expire int
	Secret string
}

// ApiClientConfig represents external/internal API information
// It is used for adapters of models on another external/internal domain
type ApiClientConfig struct {
	Host string
	Timeout int
}

// It gets a instance for ApiClientConfig 
func GetApiClientConfig(c *config.Config) ApiClientConfig {
	apiClientConfig := ApiClientConfig{"", 35}
	
	config, err := c.Get("apiClient")
	
	if (err != nil) {
		return apiClientConfig
	}
	
	host, err := config.String("host")
	
	if (err == nil) {
		apiClientConfig.Host = host
	}
	
	timeout, err := config.Int("timeout")
	
	if (err == nil) {
		apiClientConfig.Timeout = timeout
	}
		
	return apiClientConfig
}

// GlobalConfig represents built-in general configs
type GlobalConfig struct {
	EnabledCors bool
}

// It gets a instance for GlobalConfig 
func GetGlobalConfig(c *config.Config) GlobalConfig {
	enabledCors, err := c.Bool("enabledCors")
	
	if (err != nil) {
		enabledCors = false
	}
		
	return GlobalConfig{enabledCors}
}

// It gets a instance for GetJwtConfig
func GetJwtConfig(c *config.Config) (JwtConfig, error) {
	var jwtConfig JwtConfig  
	
	c, err := c.Get("jwt")
	
	if (err != nil) {
		return jwtConfig, err
	}
	
	secret, err := c.String("secret")
	
	if (err != nil) {
		return jwtConfig, err
	}
	
	expire, err := c.Int("expire")
	
	if (err != nil) {
		expire = 72
	}
	
	return JwtConfig{expire, secret}, nil
}

// SessionConfig represents Session information
type SessionConfig struct {
	Secret string
}

// It gets a instance for SessionConfig
func GetSessionConfig(c *config.Config) (SessionConfig, error) {
	var sessionConfig SessionConfig
	
	c, err := c.Get("session")
	
	if (err != nil) {
		return sessionConfig, err
	}
	
	secret, err := c.String("secret")
	
	if (err != nil) {
		return sessionConfig, err
	}
	
	sessionConfig = SessionConfig{secret}
	
	return sessionConfig, nil
}

// DBConfig represents Database information
type DBConfig struct {
	Driver string
	Datasource string
	MaxOpenConns int
}

const DEFAULT_MAX_OPEN_CONNS = 32

// It gets a instance for DBConfig
func GetDBConfig(c *config.Config) (DBConfig, error) {
	var dbConfig DBConfig  
	
	c, err := c.Get("database")
	
	if (err != nil) {
		return dbConfig, err
	}
	
	driver, err := c.String("driver")
	
	if (err != nil) {
		return dbConfig, err
	}
	
	datasource, err := c.String("datasource")
	
	if (err != nil) {
		return dbConfig, err
	}
	
	maxOpenConns, err := c.Int("maxOpenConns")
	
	if (err != nil) {
		maxOpenConns = DEFAULT_MAX_OPEN_CONNS
	}
	
	dbConfig = DBConfig{driver, datasource, maxOpenConns}  
	
	return dbConfig, nil
}