package dockerhub

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "time"
)

type Client interface {
	CheckRepo(repo string) (bool, error)
    CheckTags(repo string) (bool, error)
}

func MustNewClient() Client {
    return &client{
        httpClient: &http.Client{
            Timeout: time.Second * 5,
        }}
}

type client struct {
    httpClient *http.Client
}

type token struct {
    Token string `json:"token"`
}

type imageTags struct {
    Name string   `json:"name"`
    Tags []string `json:"tags"`
}

func (c *client) sendRequest(endpoint, token string) (*http.Response, error) {
    var (
        req *http.Request
        err error
    )
    req, err = http.NewRequest("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }
    if token != "" {
        req.Header.Add("Authorization", "Bearer "+token)
    }
    res, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    // this check could be improved
    if res.StatusCode > 300 {
        return res, errors.New(fmt.Sprintf("unexpected http status code %d", res.StatusCode))
    }
    return res, nil
}

func (c *client) CheckRepo(repo string) (bool, error) {
		token, err := c.getToken(repo)
		if err != nil {
			return false, err
		}
		res, err := c.sendRequest(fmt.Sprintf("https://index.docker.io/v2/%s/tags/list", repo), token)
		if err != nil {
			// This means that the repo exists, but has no tags. There is no actual endpoint to check for repo existence
			if res != nil && res.StatusCode == 404 {
				return true, nil
			}
			return false, err
		}
	return true, nil
}

func (c *client) CheckTags(repo string) (bool, error) {
	tags, err := c.getRepoTags(repo)
	if err != nil {
		return false, err
	}
	return len(tags.Tags) > 0, nil
}

func (c *client) getRepoTags(repo string) (*imageTags, error) {
    token, err := c.getToken(repo)
    if err != nil {
        return nil, err
    }
    res, err := c.sendRequest(fmt.Sprintf("https://index.docker.io/v2/%s/tags/list", repo), token)
    if err != nil {
        return nil, err
    }
    if res.StatusCode == 404 {
    	return &imageTags{
    		Tags: []string{},
		}, nil
	}
    if res.Body != nil {
        defer res.Body.Close()
        var tags imageTags
        err := unmarshalResponse(res.Body, &tags)
        if err != nil {
            return nil, err
        }
        return &tags, nil
	}
    return nil, nil
}

func (c *client) getToken(repo string) (string, error) {
    res, err := c.sendRequest(fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", repo), "")
    if err != nil {
        return "", err
    }
    // this check could be improved
    if res.StatusCode > 300 {
        return "", errors.New(fmt.Sprintf("unexpected http status code %d", res.StatusCode))
    }
    if res.Body != nil {
        defer res.Body.Close()
        var t token
        err := unmarshalResponse(res.Body, &t)
        if err != nil {
            return "", err
        }
        return t.Token, nil
    } else {
        return "", errors.New("response body is nil")
    }
}

func unmarshalResponse(b io.ReadCloser, target interface{}) error {
    body, err := ioutil.ReadAll(b)
    if err != nil {
        return err
    }
    err = json.Unmarshal(body, target)
    if err != nil {
        return err
    }
    return nil
}
