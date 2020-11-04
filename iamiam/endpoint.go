package iamiam

import (
	"fmt"

	"golang.org/x/oauth2"
)

// CreateInfoURL creates an URL for requesting profile information by token.
func CreateInfoURL(address string, token string) string {
	infoURL := fmt.Sprintf("http://%s/info?access_token=", address)

	return infoURL + token
}

// CreateEndpoint is Iamiam's OAuth endpoint
func CreateEndpoint(address string) oauth2.Endpoint {
	authURL := fmt.Sprintf("http://%s/auth", address)
	tokenURL := fmt.Sprintf("http://%s/token", address)

	endpoint := oauth2.Endpoint{
		AuthURL:   authURL,
		TokenURL:  tokenURL,
		AuthStyle: oauth2.AuthStyleAutoDetect,
	}

	return endpoint
}
