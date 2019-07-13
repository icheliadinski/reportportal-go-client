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

// NewLaunch creates new launch for specified client
func NewLaunch(client *Client, name, description, mode string, tags []string) *Launch {
	return &Launch{
		Name:        name,
		Description: description,
		Mode:        mode,
		Tags:        tags,
		client:      client,
	}
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

	auth := fmt.Sprintf("Bearer %s", l.client.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := doRequest(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("[WARN] failed to close body for response")
		}
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to execute POST request %s", req.URL)
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("failed with status %s", resp.Status)
	}

	v := struct {
		Id string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return errors.Wrapf(err, "failed to decode response from %s", req.URL)
	}
	l.Id = v.Id
	return nil
}

// Stop stops the launch
func (l *Launch) Stop(status string) error {
	return l.finalize(status, ActionStop)
}

// Finish finishes launch
func (l *Launch) Finish(status string) error {
	return l.finalize(status, ActionFinish)
}

// Delete delete launch
func (l *Launch) Delete() error {
	url := fmt.Sprintf("%s/%s/launch/%s", l.client.Endpoint, l.client.Project, l.Id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create DELETE request for %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", l.client.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := doRequest(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("[WARN] failed to close body response")
		}
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to execute PUT request %s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}

// Update updates launch
func (l *Launch) Update(description, mode string, tags []string) error {
	url := fmt.Sprintf("%s/%s/launch/%s/update", l.client.Endpoint, l.client.Project, l.Id)
	data := struct {
		Description string   `json:"description"`
		Mode        string   `json:"mode"`
		Tags        []string `json:"tags"`
	}{description, mode, tags}

	b, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal json %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create PUT request for %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", l.client.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := doRequest(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("[WARN] failed to close body response")
		}
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to execute PUT request %s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}

func (l *Launch) finalize(status, action string) error {
	url := fmt.Sprintf("%s/%s/launch/%s/%s", l.client.Endpoint, l.client.Project, l.Id, action)
	data := struct {
		Status  string `json:"status"`
		EndTime int64  `json:"end_time"`
	}{status, toTimestamp(time.Now())}

	b, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal object, %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create PUT request to %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", l.client.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := doRequest(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("[WARN] Failed to close response body")
		}
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to execute PUT request %s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}
