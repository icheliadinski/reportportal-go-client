package reportportal

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	Token    string
	Endpoint string
	Launch   string
	Project  string
}

func (c *Client) CheckConnect() error {
	url := fmt.Sprintf("%s/user", c.Endpoint)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrapf(err, "can't create a new request for %s", url)
	}
	r.Header.Set("Authorization", "Bearer "+c.Token)

	client := http.Client{}

	resp, err := client.Do(r)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		return errors.Wrapf(err, "failed to execute GET request %s", r.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}
