package rest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type client struct {
	httpClient *http.Client
}

type Client interface {
	CheckEndpointExists(endpoint string) (bool, error)
	CheckEndpointResult(endpoint, expectedResult string) (bool, error)
}

func MustNewClient() Client {
	return &client{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *client) sendRequest(method, endpoint string) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)
	req, err = http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// this check could be improved
	if res.StatusCode > 300 {
		return nil, errors.New(fmt.Sprintf("unexpected http status code %d", res.StatusCode))
	}
	return res, nil
}

func (c *client) CheckEndpointResult(endpoint, expectedResult string) (bool, error) {
	r, err := c.sendRequest("GET", endpoint)
	if err != nil {
		return false, err
	}
	if r.Body != nil {
		defer r.Body.Close()
		response, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return false, fmt.Errorf("failed to read response body", err)
		}
		return string(response) == expectedResult, nil
	} else {
		return false, errors.New("response body is nil")
	}
}

func (c *client) CheckEndpointExists(endpoint string) (bool, error) {
	r, err := c.sendRequest("GET", endpoint)
	if err != nil {
		return false, err
	}
	if r.StatusCode == 200 {
		return true, nil
	}
	return true, nil
}
