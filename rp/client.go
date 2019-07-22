package rp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// History defines activity history
type ActivityHistory struct {
	Field    string `json:"field"`
	NewValue string `json:"newValue"`
	OldValue string `json:"oldValue"`
}

// Content defines content history
type ActivityContent struct {
	ActionType       string             `json:"actionType"`
	ActivityId       string             `json:"activityId"`
	History          []*ActivityHistory `json:"history"`
	LastModifiedDate time.Time          `json:"lastModifiedDate"`
	LoggedObjectRef  string             `json:"loggedObjectRef"`
	ObjectName       string             `json:"objectName"`
	ObjectType       string             `json:"objectType"`
	ProjectRef       string             `json:"projectRef"`
	UserRef          string             `json:"userRef"`
}

// Page defines page info for activity
type ActivityPage struct {
	Number        int `json:"number"`
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}

// Activity defines users activity on the project
type Activity struct {
	Content []*ActivityContent `json:"content"`
	Page    *ActivityPage      `json:"page"`
}

// Widget defines widget info
type Widget struct {
	Id       string `json:"widgetId"`
	Size     []int  `json:"widgetSize"`
	Position []int  `json:"widgetPosition"`
}

// Dashboard defines dashoard info
type Dashboard []struct {
	Owner   string    `json:"owner"`
	Share   bool      `json:"share"`
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Widgets []*Widget `json:"widgets"`
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

// GetDashboard gets all dashboard resources for project
func (c *Client) GetDashboard() (*Dashboard, error) {
	url := fmt.Sprintf("%s/%s/dashboard", c.Endpoint, c.Project)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create request for %s", url)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := doRequest(req, c.Token)
	defer resp.Body.Close()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute GET request for %s", url)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed with status %s", resp.Status)
	}

	var d *Dashboard
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, errors.Wrap(err, "failed to decode response for dashboard")
	}
	return d, nil
}

// GetActivity gets all activity info for project
func (c *Client) GetActivity() (*Activity, error) {
	url := fmt.Sprintf("%s/%s/activity", c.Endpoint, c.Project)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create GET request for %s", url)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := doRequest(req, c.Token)
	defer resp.Body.Close()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute GET request for %s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed with status %s", resp.Status)
	}

	var a *Activity
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, errors.Wrapf(err, "failed to decode response for dashboard")
	}
	return a, nil
}
