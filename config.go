package koala

import (
  	"io/ioutil"
  	"path"

  	"github.com/kardianos/osext"
  	"github.com/olebedev/config"
)

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

type JwtConfig struct {
	Expire int
	Secret string
}

type ApiClientConfig struct {
	Host string
	Timeout int
}

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

type GlobalConfig struct {
	EnabledCors bool
}

func GetGlobalConfig(c *config.Config) GlobalConfig {
	enabledCors, err := c.Bool("enabledCors")
	
	if (err != nil) {
		enabledCors = false
	}
		
	return GlobalConfig{enabledCors}
}

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

type SessionConfig struct {
	Secret string
}

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

type DBConfig struct {
	Driver string
	Datasource string
	MaxOpenConns int
}

const DEFAULT_MAX_OPEN_CONNS = 32

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