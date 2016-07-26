package stormpath

import (
	"fmt"
	"os"
	"strings"

	"time"

	"github.com/spf13/viper"
)

/*
---
stormpath:
  client:
    apiKey:
      file: null
      id: null
      secret: null
    cacheManager:
      enabled: true
      defaultTtl: 300 # seconds
      defaultTti: 300
      caches: #Per resource cacehe config
    baseUrl: "https://api.stormpath.com/v1"
    connectionTimeout: 30 # seconds
    authenticationScheme: "SAUTHC1"
    proxy:
      port: null
      host: null
      username: null
      password: null
*/

//ClientConfiguration representd the overall SDK configuration options
type ClientConfiguration struct {
	APIKeyFile           string
	APIKeyID             string
	APIKeySecret         string
	CacheManagerEnabled  bool
	CacheTTL             time.Duration
	CacheTTI             time.Duration
	BaseURL              string
	ConnectionTimeout    int
	AuthenticationScheme string
	ProxyPort            int
	ProxyHost            string
	ProxyUsername        string
	ProxyPassword        string
}

//LoadConfiguration loads the configuration from the default locations
func LoadConfiguration() (ClientConfiguration, error) {
	c := newDefaultClientConfiguration()
	v := viper.New()

	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigName("stormpath")
	v.AddConfigPath(os.Getenv("HOME") + "/.stormpath")
	v.AddConfigPath(".")
	v.ReadInConfig()

	c.APIKeyID = v.GetString("stormpath.client.apiKey.id")
	c.APIKeySecret = v.GetString("stormpath.client.apiKey.secret")
	id, secret, err := loadCredentials(c.APIKeyFile)
	if err == nil {
		c.APIKeyID = id
		c.APIKeySecret = secret
	}

	if c.APIKeyID == "" && c.APIKeySecret == "" {
		return c, fmt.Errorf("API credentials couldn't be loaded")
	}

	if v.Get("stormpath.client.cacheManager.enabled") != nil {
		c.CacheManagerEnabled = v.GetBool("stormpath.client.cacheManager.enabled")
	}
	if v.Get("stormpath.client.cacheManager.defaultTtl") != nil {
		c.CacheTTL = time.Duration(v.GetInt("stormpath.client.cacheManager.defaultTtl")) * time.Second
	}
	if v.Get("stormpath.client.cacheManager.defaultTti") != nil {
		c.CacheTTI = time.Duration(v.GetInt("stormpath.client.cacheManager.defaultTti")) * time.Second
	}

	if v.GetString("stormpath.client.baseUrl") != "" {
		c.BaseURL = v.GetString("stormpath.client.baseUrl")
	}
	if v.Get("stormpath.client.connectionTimeout") != nil {
		c.ConnectionTimeout = v.GetInt("stormpath.client.connectionTimeout")
	}

	if v.GetString("stormpath.client.authenticationScheme") != "" {
		c.AuthenticationScheme = v.GetString("stormpath.client.authenticationScheme")
	}

	c.ProxyHost = v.GetString("stormpath.client.proxy.host")
	c.ProxyPort = v.GetInt("stormpath.client.proxy.port")
	c.ProxyUsername = v.GetString("stormpath.client.proxy.username")
	c.ProxyPassword = v.GetString("stormpath.client.proxy.password")

	return c, nil
}

func LoadConfigurationWithCreds(key string, secret string) ClientConfiguration {
	c := newDefaultClientConfiguration()
	c.APIKeyID = key
	c.APIKeySecret = secret

	return c
}

func loadCredentials(extraFileLocation string) (id string, secret string, err error) {
	id = os.Getenv("STORMPATH_API_KEY_ID")
	secret = os.Getenv("STORMPATH_API_KEY_SECRET")
	if id != "" && secret != "" {
		return id, secret, nil
	}

	id, secret, err = loadCredentialsFromFile(os.Getenv("HOME") + "/.stormpath/apiKey.properties")
	if err == nil && id != "" && secret != "" {
		return
	}

	id, secret, err = loadCredentialsFromFile("./apiKey.properties")
	if err == nil && id != "" && secret != "" {
		return
	}

	if extraFileLocation != "" {
		return loadCredentialsFromFile(extraFileLocation)
	}
	return
}

//GetJWTSigningKey returns the API Key Secret as a []byte to sign JWT tokens
func (config ClientConfiguration) GetJWTSigningKey() []byte {
	return []byte(config.APIKeySecret)
}

func newDefaultClientConfiguration() (c ClientConfiguration) {
	return ClientConfiguration{
		APIKeyFile:           "",
		APIKeyID:             "",
		APIKeySecret:         "",
		CacheManagerEnabled:  true,
		CacheTTI:             300 * time.Second,
		CacheTTL:             300 * time.Second,
		BaseURL:              "https://api.stormpath.com/v1/",
		ConnectionTimeout:    30,
		AuthenticationScheme: "SAUTHC1",
		ProxyHost:            "",
		ProxyPort:            0,
		ProxyUsername:        "",
		ProxyPassword:        "",
	}
}
