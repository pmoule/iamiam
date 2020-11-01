package iamiam

import (
	"fmt"

	"golang.org/x/oauth2"
)

// CreateInfoURL creates an URL for requesting profile information by token.
func CreateInfoURL(hostname string, port int, token string) string {
	address := createAddress(hostname, port)
	infoURL := fmt.Sprintf("http://%s/info?access_token=", address)

	return infoURL + token
}

// CreateEndpoint is Iamiam's OAuth endpoint
func CreateEndpoint(hostname string, port int) oauth2.Endpoint {
	address := createAddress(hostname, port)
	authURL := fmt.Sprintf("http://%s/auth", address)
	tokenURL := fmt.Sprintf("http://%s/token", address)

	endpoint := oauth2.Endpoint{
		AuthURL:   authURL,
		TokenURL:  tokenURL,
		AuthStyle: oauth2.AuthStyleAutoDetect,
	}

	return endpoint
}

func createAddress(hostname string, port int) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}
