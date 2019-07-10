package rp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Launch defines launch info
type Launch struct {
	Id          string
	Name        string
	Description string
	Mode        string
	StartTime   time.Time
	Tags        []string

	client *Client
}

// Start starts the launch
func (l *Launch) Start() error {
	url := fmt.Sprintf("%s/%s/launch", l.client.Endpoint, l.client.Project)
	launch := struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Mode        string   `json:"mode"`
		Tags        []string `json:"tags,omitempty"`
		StartTime   int64    `json:"start_time"`
	}{l.Name, l.Description, l.Mode, l.Tags, toTimestamp(time.Now())}

	b, err := json.Marshal(&launch)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal object, %v", launch)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create request for %s", url)
	}

	addContentTypeJSON(req)

	resp, err := doRequest(req, l.client.Token)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[WARN] failed to close body for response to %s", req.URL)
		}
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to execute POST request %s", req.URL)
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("failed with status %s", resp.Status)
	}

	v := struct {
		Id string `json:"id"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return errors.Wrapf(err, "failed to decode response from %s", req.URL)
	}

	l.Id = v.Id
	return nil
}

// Stop stops the launch
func (l *Launch) Stop() error {
	return nil
}

// Finish finishes launch
func (l *Launch) Finish() error {
	return nil
}

// Delete delete launch
func (l *Launch) Delete() error {
	return nil
}

// Update updates launch
func (l *Launch) Update() error {
	return nil
}
