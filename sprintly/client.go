package sprintly

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	LibraryVersion = "0.0.1"

	DefaultBaseURL   = "https://sprint.ly/api/"
	DefaultUserAgent = "go-sprintly/" + LibraryVersion
)

type Client struct {
	// Sprintly username to be used to authenticate API calls.
	username string

	// Sprintly access token to be used to authenticate API calls.
	token string

	// HTTP client to be used for communication with the Sprintly API.
	client *http.Client

	// Base URL of the Sprintly API that is to be used to form endpoint URLs.
	baseURL *url.URL

	// User-Agent header to use when making API calls.
	userAgent string

	// The Items service.
	Items *ItemsService

	// The Deploys service.
	Deploys *DeploysService
}

// NewClient returns a new API client instance that uses
// the given username and token to authenticate the API calls.
func NewClient(username, token string) *Client {
	baseURL, _ := url.Parse(DefaultBaseURL)
	client := &Client{
		token:     token,
		client:    http.DefaultClient,
		baseURL:   baseURL,
		userAgent: DefaultUserAgent,
	}
	client.Items = newItemsService(client)
	client.Deploys = newDeploysService(client)
	return client
}

// SetBaseURL can be used to overwrite the default API base URL,
// which is the Sprintly API - https://sprint.ly/api/.
func (c *Client) SetBaseURL(baseURL string) error {
	// Parse the URL.
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	// Make sure the trailing slash is there.
	if u.Path != "" && u.Path[len(u.Path)-1] != '/' {
		u.Path += "/"
	}

	c.baseURL = u
	return nil
}

// SetUserAgent can be used to overwrite the default user agent string.
func (c *Client) SetUserAgent(agent string) {
	c.userAgent = agent
}

// SetHttpClient can be used to really customize the API client behaviour
// by replacing the underlying HTTP client that is being used to carry out
// all the API calls.
func (c *Client) SetHttpClient(client *http.Client) {
	c.client = client
}

// NewRequest returns a new API requests for the given method and relative URL.
//
// In case the body object is not nil, it is marshaled and send in the request body.
func (c *Client) NewRequest(method, urlPath string, body interface{}) (*http.Request, error) {
	path, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(path)

	var bodyBuffer bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&bodyBuffer).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), &bodyBuffer)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	return req, nil
}

// Do carries out the given API request.
//
// In case the interface passed into Do is not nil, it is filled from the response body.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return resp, &ErrAPI{
			Response: resp,
		}
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}
