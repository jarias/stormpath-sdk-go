package stormpath

import (
	"os"

	"github.com/dmotylev/goproperties"
)

//Credentials represents a set of Stormpath credentials
type Credentials struct {
	ID     string
	Secret string
}

//NewCredentialsFromFile creates a new credentials from a Stormpath key files
func NewCredentialsFromFile(file string) (Credentials, error) {
	c := Credentials{}

	p, err := properties.Load(file)

	if err != nil {
		return Credentials{}, err
	}

	c.ID = p.String("apiKey.id", "")
	c.Secret = p.String("apiKey.secret", "")

	return c, err
}

//NewDefaultCredentials would create a new credentials based on env variables first then the default file location
//at ~/.config/stormpath/apiKey.properties
func NewDefaultCredentials() (Credentials, error) {
	apiKeyID := os.Getenv("STORMPATH_API_KEY_ID")
	apiKeySecret := os.Getenv("STORMPATH_API_KEY_SECRET")
	if apiKeyID != "" && apiKeySecret != "" {
		return Credentials{apiKeyID, apiKeySecret}, nil
	}
	return NewCredentialsFromFile(os.Getenv("HOME") + "/.config/stormpath/apiKey.properties")
}
