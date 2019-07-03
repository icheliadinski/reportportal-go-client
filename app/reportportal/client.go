package reportportal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	ModeDebug   = "DEBUG"
	ModeDefault = "DEFAULT"

	StatusPassed   = "PASSED"
	StatusFailed   = "FAILED"
	StatusStopped  = "STOPPED"
	StatusSkipped  = "SKIPPED"
	StatusReseted  = "RESETED"
	StatusCanceled = "CANCELED"

	ActionStop   = "stop"
	ActionFinish = "finish"
)

// Client defines a report portal client
type Client struct {
	Endpoint string `short:"e" long:"endpoint" env:"ENDPOINT" description:"report portal endpoint"`
	Token    string `short:"t" long:"token" env:"TOKEN" description:"user token for report portal"`
	Launch   string `short:"l" long:"launch" env:"LAUNCH" description:"launch name"`
	Project  string `short:"p" long:"project" env:"PROJECT" description:"project name"`
}

// LaunchInfo defines launch object
type LaunchInfo struct {
	Id     string `json:"id"`
	Number int64  `json:"number"`
}

// TestItemInfo defines test information
type TestItemInfo struct {
	Id       string `json:"id"`
	UniqueId string `json:"uniqueId"`
}

// ProjectSettings defines project settings
type ProjectSettings struct {
	StatisticsStrategy string                 `json:"statisticsStrategy"`
	Name               string                 `json:"project"`
	SubTypes           map[string]interface{} `json:"subTypes"`
}

// NewClient defines function constructor for client
func NewClient(endpoint, token, launch, project string) *Client {
	e := strings.TrimSuffix(endpoint, "/")
	client := &Client{e, token, launch, project}
	return client
}

// CheckConnect defines check for connection
func (c *Client) CheckConnect() error {
	url := fmt.Sprintf("%s/user", c.Endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrapf(err, "can't create a new request for %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", c.Token)
	req.Header.Set("Authorization", auth)

	client := http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("[WARN] failed to close body for response")
		}
	}()

	if err != nil {
		return errors.Wrapf(err, "failed to execute GET request %s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}

// StartLaunch runs launch
func (c *Client) StartLaunch(name, description string, mode string, tags []string, startTime time.Time) (string, error) {
	url := fmt.Sprintf("%s/%s/launch", c.Endpoint, c.Project)
	launch := struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Mode        string   `json:"mode"`
		Tags        []string `json:"tags,omitempty"`
		StartTime   int64    `json:"start_time"`
	}{
		Name:        name,
		Description: description,
		Mode:        mode,
		Tags:        tags,
		StartTime:   startTime.UnixNano(),
	}

	b, err := json.Marshal(&launch)
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal object, %v", launch)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create request for %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", c.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("[WARN] failed to close body for response")
		}
	}()
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute POST request %s", req.URL)
	}
	if resp.StatusCode != http.StatusCreated {
		return "", errors.Errorf("failed with status %s", resp.Status)
	}

	v := LaunchInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", errors.Wrapf(err, "failed to decode response from %s", req.URL)
	}
	return v.Id, nil
}

// StopLaunch stops the exact launch
func (c *Client) StopLaunch(id, status string, endTime time.Time) error {
	return c.finalizeLaunch(id, ActionStop, status, endTime)
}

// FinishLaunch finishes exact launch
func (c *Client) FinishLaunch(id, status string, endTime time.Time) error {
	return c.finalizeLaunch(id, ActionFinish, status, endTime)
}

// UpdateLaunch updates launch info
func (c *Client) UpdateLaunch(id, description, mode string, tags []string) error {
	url := fmt.Sprintf("%s/%s/launch/%s/update", c.Endpoint, c.Project, id)
	data := struct {
		Description string   `json:"description"`
		Mode        string   `json:"mode"`
		Tags        []string `json:"tags"`
	}{
		Description: description,
		Mode:        mode,
		Tags:        tags,
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal json %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create request for %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", c.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
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

// StartTestItem starts a test item suite/story/test etc
func (c *Client) StartTestItem(launchId, name, description, itemType, parentId string, tags []string, startTime time.Time) (string, error) {
	url := fmt.Sprintf("%s/%s/item/%s", c.Endpoint, c.Project, parentId)
	data := struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
		StartTime   int64    `json:"start_time"`
		LaunchId    string   `json:"launch_id"`
		Type        string   `json:"type"`
		Parameters  []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"parameters"`
	}{
		Name:        name,
		Description: description,
		Tags:        tags,
		StartTime:   startTime.UnixNano(),
		LaunchId:    launchId,
		Type:        itemType,
		Parameters:  nil,
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal object %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create POST request to %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", c.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("[WARN] failed to close response body")
		}
	}()
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute POST request %s", req.URL)
	}
	if resp.StatusCode != http.StatusCreated {
		return "", errors.Errorf("failed with status %s", resp.Status)
	}

	v := LaunchInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", errors.Wrapf(err, "failed to decode response from %s", req.URL)
	}
	return v.Id, nil
}

// FinishTestItem finishes specified test item with specific status
func (c *Client) FinishTestItem(id, status string, endTime time.Time) error {
	url := fmt.Sprintf("%s/%s/item/%s", c.Endpoint, c.Project, id)
	data := struct {
		EndTime int64  `json:"end_time"`
		Status  string `json:"status"`
	}{
		EndTime: endTime.UnixNano(),
		Status:  status,
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal request data %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create PUT request to %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", c.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("[WARN] failed to close response body")
		}
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to execute PUT request to %s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}

func (c *Client) GetProjectSettings() (ProjectSettings, error) {
	url := fmt.Sprintf("%s/%s/settings", c.Endpoint, c.Project)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ProjectSettings{}, errors.Wrapf(err, "failed to create GET request to %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", c.Token)
	req.Header.Set("Authorization", auth)

	client := http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("[WARN] failed to close body from response")
		}
	}()
	if err != nil {
		return ProjectSettings{}, errors.Wrapf(err, "failed to GET to %s", url)
	}

	v := ProjectSettings{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return ProjectSettings{}, errors.Wrapf(err, "failed to decode response from %s", url)
	}
	return v, nil
}

// finalizeLaunch finalizes exact match with specific action
func (c *Client) finalizeLaunch(id, action, status string, endTime time.Time) error {
	url := fmt.Sprintf("%s/%s/launch/%s/%s", c.Endpoint, c.Project, id, action)
	data := struct {
		Status  string `json:"status"`
		EndTime int64  `json:"end_time"`
	}{
		Status:  status,
		EndTime: endTime.UnixNano(),
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal object, %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create request to %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", c.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
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
