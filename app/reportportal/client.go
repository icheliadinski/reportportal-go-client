package reportportal

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Client defines a report portal client
type Client struct {
	Token    string `short:"t" long:"token" env:"TOKEN" description:"user token for report portal"`
	Endpoint string `short:"e" long:"endpoint" env:"ENDPOINT" description:"report portal endpoint"`
	Launch   string `short:"l" long:"launch" env:"LAUNCH" description:"launch name"`
	Project  string `short:"p" long:"project" env:"PROJECT" description:"project name"`
}

// NewClient defines function constructor for client
func NewClient(endpoint string, token string, launch string, project string) *Client {
	e := strings.TrimSuffix(endpoint, "/")
	t := token
	return &Client{
		Endpoint: e,
		Token:    t,
		Launch:   launch,
		Project:  project,
	}
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
func (c *Client) StartLaunch(name string, description string, tags []string) error {
	return nil
}
