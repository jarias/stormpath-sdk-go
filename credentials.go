package stormpath

import (
	"github.com/dmotylev/goproperties"
	"os"
)

type Credentials struct {
	Id     string
	Secret string
}

func NewCredentialsFromFile(file string) (*Credentials, error) {
	c := new(Credentials)

	p, err := properties.Load(file)

	if err != nil {
		return nil, err
	}

	c.Id = p.String("apiKey.id", "")
	c.Secret = p.String("apiKey.secret", "")

	return c, err
}

func NewDefaultCredentials() (*Credentials, error) {
	apiKeyId := os.Getenv("STORMPATH_API_KEY_ID")
	apiKeySecret := os.Getenv("STORMPATH_API_KEY_SECRET")
	if apiKeyId != "" && apiKeySecret != "" {
		return &Credentials{apiKeyId, apiKeySecret}, nil
	} else {
		defaultFilePath := os.Getenv("HOME") + "/.config/stormpath/apiKey.properties"
		return NewCredentialsFromFile(defaultFilePath)
	}
}
