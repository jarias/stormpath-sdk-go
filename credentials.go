package stormpath

import "github.com/dmotylev/goproperties"

func loadCredentialsFromFile(file string) (id string, secret string, err error) {
	p, err := properties.Load(file)

	if err != nil {
		return "", "", err
	}

	return p.String("apiKey.id", ""), p.String("apiKey.secret", ""), err
}
