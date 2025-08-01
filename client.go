package wiremock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/wiremock/go-wiremock/journal"
)

const (
	wiremockAdminURN         = "__admin"
	wiremockAdminMappingsURN = "__admin/mappings"
	wiremockAdminRequestsURN = "__admin/requests"
)

// A Client implements requests to the wiremock server.
type Client struct {
	url string
}

// NewClient returns *Client.
func NewClient(url string) *Client {
	return &Client{url: url}
}

// StubFor creates a new stub mapping.
func (c *Client) StubFor(stubRule *StubRule) error {
	requestBody, err := stubRule.MarshalJSON()
	if err != nil {
		return fmt.Errorf("build stub request error: %w", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/%s", c.url, wiremockAdminMappingsURN), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("stub request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %w", err)
		}

		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}

// Clear deletes all stub mappings.
func (c *Client) Clear() error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", c.url, wiremockAdminMappingsURN), nil)
	if err != nil {
		return fmt.Errorf("build cleare Request error: %w", err)
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return fmt.Errorf("clear Request error: %w", err)
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
		return fmt.Errorf("reset Request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %w", err)
		}

		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}

// ResetAllScenarios resets back to start of the state of all configured scenarios.
func (c *Client) ResetAllScenarios() error {
	res, err := http.Post(fmt.Sprintf("%s/%s/scenarios/reset", c.url, wiremockAdminURN), "application/json", nil)
	if err != nil {
		return fmt.Errorf("reset all scenarios Request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %w", err)
		}

		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}

// GetCountRequests gives count requests by criteria.
func (c *Client) GetCountRequests(r *Request) (int64, error) {
	requestBody, err := r.MarshalJSON()
	if err != nil {
		return 0, fmt.Errorf("get count requests: build error: %w", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/%s/requests/count", c.url, wiremockAdminURN), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, fmt.Errorf("get count requests: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("get count requests: read response error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("get count requests: bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	var countRequestsResponse struct {
		Count int64 `json:"count"`
	}

	err = json.Unmarshal(bodyBytes, &countRequestsResponse)
	if err != nil {
		return 0, fmt.Errorf("get count requests: read json error: %w", err)
	}

	return countRequestsResponse.Count, nil
}

// Verify checks count of request sent.
func (c *Client) Verify(r *Request, expectedCount int64) (bool, error) {
	actualCount, err := c.GetCountRequests(r)
	if err != nil {
		return false, err
	}

	return actualCount == expectedCount, nil
}

// GetAllRequests returns all requests logged in the journal.
func (c *Client) GetAllRequests() (*journal.GetAllRequestsResponse, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", c.url, wiremockAdminRequestsURN))
	if err != nil {
		return nil, fmt.Errorf("get all requests: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("get all requests: read response error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get all requests: bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	var response journal.GetAllRequestsResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("get all requests: error unmarshalling response: %w", err)
	}
	return &response, nil
}

// GetRequestByID retrieves a single request from the journal, by its ID.
func (c *Client) GetRequestByID(requestID string) (*journal.GetRequestResponse, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s/%s", c.url, wiremockAdminRequestsURN, requestID))
	if err != nil {
		return nil, fmt.Errorf("get request by id: build request error: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("get request by id: read response error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get request by id: bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	var response journal.GetRequestResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("get request by id: error unmarshalling response: %w", err)
	}
	return &response, nil
}

// FindRequestsByCriteria returns all requests in the journal matching the criteria.
func (c *Client) FindRequestsByCriteria(r *Request) (*journal.FindRequestsByCriteriaResponse, error) {
	requestBody, err := r.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("find requests by criteria: build error: %w", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/%s/find", c.url, wiremockAdminRequestsURN), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("find requests by criteria: request error: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("find requests by criteria: read response error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("find requests by criteria: bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	var requests journal.FindRequestsByCriteriaResponse
	err = json.Unmarshal(bodyBytes, &requests)
	if err != nil {
		return nil, fmt.Errorf("find requests by criteria: read json error: %w", err)
	}
	return &requests, nil
}

// FindUnmatchedRequests returns all requests in the journal matching the criteria.
func (c *Client) FindUnmatchedRequests() (*journal.FindUnmatchedRequestsResponse, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s/unmatched", c.url, wiremockAdminRequestsURN))
	if err != nil {
		return nil, fmt.Errorf("find unmatched requests: request error: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("find unmatched requests: read response error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("find unmatched requests: bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	var requests journal.FindUnmatchedRequestsResponse
	err = json.Unmarshal(bodyBytes, &requests)
	if err != nil {
		return nil, fmt.Errorf("find unmatched requests: read json error: %w", err)
	}
	return &requests, nil
}

// DeleteAllRequests deletes all the requests in the journal.
func (c *Client) DeleteAllRequests() error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", c.url, wiremockAdminRequestsURN), nil)
	if err != nil {
		return fmt.Errorf("delete all requests: build error: %w", err)
	}
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return fmt.Errorf("delete all requests: request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %w", err)
		}
		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}
	return nil
}

// DeleteRequestByID deletes a single request from the journal, by its ID.
func (c *Client) DeleteRequestByID(requestID string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s", c.url, wiremockAdminRequestsURN, requestID), nil)
	if err != nil {
		return fmt.Errorf("delete request by id: build request error: %w", err)
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return fmt.Errorf("delete request by id: request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("delete request by id: read response error: %w", err)
		}
		return fmt.Errorf("delete request by id: bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}
	return nil
}

// DeleteRequestsByCriteria deletes all requests in the journal matching the criteria.
func (c *Client) DeleteRequestsByCriteria(r *Request) (*journal.DeleteRequestByCriteriaResponse, error) {
	requestBody, err := r.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("delete requests by criteria: build error: %w", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/%s/remove", c.url, wiremockAdminRequestsURN), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("delete requests by criteria: request error: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("delete requests by criteria: read response error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("delete requests by criteria: bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	var requests journal.DeleteRequestByCriteriaResponse
	err = json.Unmarshal(bodyBytes, &requests)
	if err != nil {
		return nil, fmt.Errorf("delete requests by criteria: error unmarshalling response: %w", err)
	}
	return &requests, nil
}

// DeleteStubByID deletes stub by id.
func (c *Client) DeleteStubByID(id string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s", c.url, wiremockAdminMappingsURN, id), nil)
	if err != nil {
		return fmt.Errorf("delete stub by id: build request error: %w", err)
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return fmt.Errorf("delete stub by id: request error: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %w", err)
		}

		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}

// DeleteStub deletes stub mapping.
func (c *Client) DeleteStub(s *StubRule) error {
	return c.DeleteStubByID(s.UUID())
}

// StartRecording starts a recording.
func (c *Client) StartRecording(targetBaseUrl string) error {
	requestBody := fmt.Sprintf(`{"targetBaseUrl":"%s"}`, targetBaseUrl)
	res, err := http.Post(
		fmt.Sprintf("%s/%s/recordings/start", c.url, wiremockAdminURN),
		"application/json",
		strings.NewReader(requestBody),
	)
	if err != nil {
		return fmt.Errorf("start recording error: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %w", err)
		}

		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	return err
}

// StopRecording stops a recording.
func (c *Client) StopRecording() error {
	res, err := http.Post(
		fmt.Sprintf("%s/%s/recordings/stop", c.url, wiremockAdminURN),
		"application/json",
		nil,
	)
	if err != nil {
		return fmt.Errorf("stop recording error: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read response error: %w", err)
		}

		return fmt.Errorf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}
