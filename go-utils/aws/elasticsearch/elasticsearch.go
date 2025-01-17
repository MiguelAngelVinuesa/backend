package elasticsearch

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client represents an Elasticsearch client.
type Client struct {
	url  string
	user string
	pass string
}

// New returns a new Elasticsearch client.
func New(url, user, pass string) (*Client, error) {
	if url == "" {
		return nil, errors.New("empty url provided")
	}

	if l := len(url) - 1; url[l] == '/' {
		url = url[:l]
	}

	return &Client{url: url, user: user, pass: pass}, nil
}

// CreateIndex creates an Elasticsearch index.
func (c *Client) CreateIndex(idx *Index, ignoreIfExists bool) error {
	var reqData putIndexReq
	reqData.Mappings.Properties = make(map[string]putIndexField, len(idx.Fields))

	for _, field := range idx.Fields {
		reqData.Mappings.Properties[field.Name] = putIndexField{
			Type:  field.DataType,
			Index: field.Index,
		}
	}

	reqB, reqBErr := json.Marshal(reqData)
	if reqBErr != nil {
		return fmt.Errorf("failed to marshal request: %v", reqBErr)
	}

	uri := fmt.Sprintf("%s/%s", c.url, idx.Name)
	req, reqErr := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(reqB))
	if reqErr != nil {
		return fmt.Errorf("failed to create new request: %v", reqErr)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.user != "" {
		req.Header.Set("Authorization", "Basic "+basicAuth(c.user, c.pass))
	}

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return fmt.Errorf("failed to issue request: %v", respErr)
	}
	defer resp.Body.Close()

	bodyB, bodyBErr := ioutil.ReadAll(resp.Body)
	if bodyBErr != nil {
		return fmt.Errorf("failed to read request body: %v", bodyBErr)
	}

	if resp.StatusCode != http.StatusOK {
		if !ignoreIfExists || !strings.Contains(string(bodyB), "resource_already_exists_exception") {
			return fmt.Errorf("unexpected status code received: %d", resp.StatusCode)
		}
	}

	return nil
}

// AddDocument adds a document to an Elasticsearch index.
func (c *Client) AddDocument(idxName string, data map[string]interface{}) error {
	if len(data) == 0 {
		return errors.New("empty data map provided")
	}

	reqB, reqBErr := json.Marshal(data)
	if reqBErr != nil {
		return fmt.Errorf("failed to marshal data map: %v", reqBErr)
	}

	uri := fmt.Sprintf("%s/%s/_doc", c.url, idxName)
	req, reqErr := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(reqB))
	if reqErr != nil {
		return fmt.Errorf("failed to create new request: %v", reqErr)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.user != "" {
		req.Header.Set("Authorization", "Basic "+basicAuth(c.user, c.pass))
	}

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return fmt.Errorf("failed to issue request: %v", respErr)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code received: %d", resp.StatusCode)
	}

	return nil
}

// DeleteOldDocuments removes documents from an Elasticsearch index that are older than ageSec,
// and returns the number of removed documents.
func (c *Client) DeleteOldDocuments(idxName string, ageSec int) (int, error) {
	if ageSec < 60 {
		return 0, errors.New("invalid ageSec value provided")
	}

	reqRaw := fmt.Sprintf("{\"query\": {\"range\": {\"ts\": {\"lt\":%d}}}}", time.Now().UnixMilli()-int64(ageSec*1000))
	uri := fmt.Sprintf("%s/%s/_delete_by_query", c.url, idxName)
	req, reqErr := http.NewRequest(http.MethodPost, uri, strings.NewReader(reqRaw))
	if reqErr != nil {
		return 0, fmt.Errorf("failed to create new request: %v", reqErr)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.user != "" {
		req.Header.Set("Authorization", "Basic "+basicAuth(c.user, c.pass))
	}

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return 0, fmt.Errorf("failed to issue request: %v", respErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code received: %d", resp.StatusCode)
	}

	bodyB, bodyBErr := ioutil.ReadAll(resp.Body)
	if bodyBErr != nil {
		return 0, fmt.Errorf("failed to read response body: %v", bodyBErr)
	}

	var respData postDeleteResp
	respDataErr := json.Unmarshal(bodyB, &respData)
	if respDataErr != nil {
		return 0, fmt.Errorf("failed to parse response body: %v", respDataErr)
	}

	return respData.Deleted, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
