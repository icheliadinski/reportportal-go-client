package rp

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	Endpoint string
	Token    string
	Project  string
}

// NewClient creates new client for ReportPortal endpoint
func NewClient(endpoint, project, token string, apiVersion int) *Client {
	endpoint = strings.TrimSuffix(endpoint, "/")

	var esb strings.Builder
	if !strings.HasPrefix(endpoint, "https://") && !strings.HasPrefix(endpoint, "http://") {
		esb.WriteString("https://")
	}
	esb.WriteString(endpoint)

	if apiVersion < 1 {
		apiVersion = 1
	}

	if !strings.Contains(endpoint, "/api/v") {
		esb.WriteString("/api/v")
		esb.WriteString(strconv.Itoa(apiVersion))
	}

	return &Client{
		Endpoint: esb.String(),
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

	resp, err := doRequest(req, c.Token)
	defer resp.Body.Close()

	if err != nil {
		return errors.Wrapf(err, "failed to execute GET request %s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}
