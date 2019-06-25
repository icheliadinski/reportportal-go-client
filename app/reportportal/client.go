package reportportal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Client struct {
	Token    string
	Endpoint string
	Launch   string
	Project  string
}

func NewClient(endpoint string, token string) (*Client, error) {
	c := &Client{}
	c.Endpoint = strings.TrimSuffix(endpoint, "/")
	c.Token = token
	return c, nil
}

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
	defer func() { _ = resp.Body.Close() }()

	if err != nil {
		return errors.Wrapf(err, "failed to execute GET request %s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}
