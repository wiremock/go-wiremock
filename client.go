package wiremock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	wiremockAdminURN         = "__admin"
	wiremockAdminMappingsURN = "__admin/mappings"
	wiremockAdminFindURN     = "__admin/requests/find"
)

// A Client implements requests to the wiremock server
type Client struct {
	url string
}

// NewClient returns *Client.
func NewClient(url string) *Client {
	return &Client{url: url}
}

// StubFor creates a new stub mapping.
func (c *Client) StubFor(stubRule *StubRule) error {
	requestBody, err := json.Marshal(stubRule)
	if err != nil {
		return fmt.Errorf("build stub request error: %s", err.Error())
	}

	res, err := http.Post(fmt.Sprintf("%s/%s", c.url, wiremockAdminMappingsURN), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("stub request error: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %s", err.Error())
		}

		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}

// Clear deletes all stub mappings.
func (c *Client) Clear() error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", c.url, wiremockAdminMappingsURN), nil)
	if err != nil {
		return fmt.Errorf("build cleare request error: %s", err.Error())
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return fmt.Errorf("clear request error: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response status: %d", res.StatusCode)
	}

	return nil
}

// Reset restores stub mappings to the defaults defined back in the backing store.
func (c *Client) Reset() error {
	res, err := http.Post(fmt.Sprintf("%s/%s/reset", c.url, wiremockAdminMappingsURN), "application/json", nil)
	if err != nil {
		return fmt.Errorf("reset request error: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %s", err.Error())
		}

		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}

// ResetAllScenarios resets back to start of the state of all configured scenarios.
func (c *Client) ResetAllScenarios() error {
	res, err := http.Post(fmt.Sprintf("%s/%s/scenarios/reset", c.url, wiremockAdminURN), "application/json", nil)
	if err != nil {
		return fmt.Errorf("reset all scenarios request error: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %s", err.Error())
		}

		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}

// FindFor finds requests that were made for a StubRule
func (c *Client) FindRequestsFor(stubRule *StubRule) (map[string]interface{}, error) {
	//to find requests matching a stub, we need only the request portion of the stub
	requestBody, err := json.Marshal(&stubRule.request)
	if err != nil {
		return nil, fmt.Errorf("build stub request error: %s", err.Error())
	}

	res, err := http.Post(fmt.Sprintf("%s/%s", c.url, wiremockAdminFindURN), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("stub request error: %s", err.Error())
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		if err != nil {
			return nil, fmt.Errorf("read response error: %s", err.Error())
		}

		return nil, fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	var response map[string]interface{}
	err = json.Unmarshal(bodyBytes, &response)

	if err != nil {
		return nil, fmt.Errorf("error: %s unmarshalling response: %s", err, response)
	}

	return response, nil
}
