package rp

import (
	"fmt"
	"log"
	"net/http"

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
	StatusCanceled = "CANCELLED"

	ActionStop   = "stop"
	ActionFinish = "finish"

	LevelError   = "error"
	LevelWarn    = "warn"
	LevelTrace   = "trace"
	LevelInfo    = "info"
	LevelDebug   = "debug"
	LevelFatal   = "fatal"
	LevelUnknown = "unknown"
)

// Client defines a report portal client
type Client struct {
	Endpoint string `short:"e" long:"endpoint" env:"ENDPOINT" description:"report portal endpoint"`
	Token    string `short:"t" long:"token" env:"TOKEN" description:"user token for report portal"`
	Project  string `short:"p" long:"project" env:"PROJECT" description:"project name"`
}

// NewClient creates new client for ReportPortal endpoint
func NewClient(endpoint, project, token string) *Client {
	return &Client{
		Endpoint: endpoint,
		Project:  project,
		Token:    token,
	}
}

// CheckConnect checks connection to ReportPortal
func (c *Client) CheckConnect() error {
	url := fmt.Sprintf("%s/user", c.Endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrapf(err, "can't create a new request for %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", c.Token)
	req.Header.Set("Authorization", auth)

	resp, err := doRequest(req)
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

// func (c *Client) GetProjectSettings() (ProjectSettings, error) {
// 	url := fmt.Sprintf("%s/%s/settings", c.Endpoint, c.Project)

// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return ProjectSettings{}, errors.Wrapf(err, "failed to create GET request to %s", url)
// 	}

// 	auth := fmt.Sprintf("Bearer %s", c.Token)
// 	req.Header.Set("Authorization", auth)

// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	defer func() {
// 		if err := resp.Body.Close(); err != nil {
// 			fmt.Println("[WARN] failed to close body from response")
// 		}
// 	}()
// 	if err != nil {
// 		return ProjectSettings{}, errors.Wrapf(err, "failed to GET to %s", url)
// 	}

// 	v := ProjectSettings{}
// 	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
// 		return ProjectSettings{}, errors.Wrapf(err, "failed to decode response from %s", url)
// 	}
// 	return v, nil
// }
