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

// Client defines a report portal client
type Client struct {
	Endpoint string `short:"e" long:"endpoint" env:"ENDPOINT" description:"report portal endpoint"`
	Token    string `short:"t" long:"token" env:"TOKEN" description:"user token for report portal"`
	Launch   string `short:"l" long:"launch" env:"LAUNCH" description:"launch name"`
	Project  string `short:"p" long:"project" env:"PROJECT" description:"project name"`
}

type LaunchInfo struct {
	Name        string
	StartTime   time.Time
	Description string
	tags        []string
}

// NewClient defines function constructor for client
func NewClient(endpoint string, token string, launch string, project string) *Client {
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

// StartLaunch defines launch start
func (c *Client) StartLaunch(name string, description string, tags []string, startTime time.Time) error {
	url := fmt.Sprintf("%s/%s/launch", c.Endpoint, c.Project)
	launch := struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Tags        []string  `json:"tags"`
		StartTime   time.Time `json:"start_time"`
	}{
		Name:        name,
		Description: description,
		Tags:        tags,
		StartTime:   startTime,
	}

	b, err := json.Marshal(&launch)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal object, %v", launch)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create request for %s", url)
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
		return errors.Wrapf(err, "failed to execute POST request %s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}
