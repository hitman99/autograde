package rest

import (
	"net/http"
	"time"
)

type client struct {
	httpClient *http.Client
}

type Client interface {
	CheckEndpoint(endpoint string) (bool, error)
	CallEndpoint(endpoint, expectedResult string) (bool, error)
}

func MustNewClient() Client {
	return &client{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *client) CheckEndpoint(endpoint string) (bool, error) {
	return false, nil
}

func (c *client) CallEndpoint(endpoint, expectedResult string) (bool, error) {
	return false, nil
}
