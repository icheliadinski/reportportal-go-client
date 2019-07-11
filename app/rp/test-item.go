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

const (
	TestItemSuite        = "SUITE"
	TestItemStory        = "STORY"
	TestItemTest         = "TEST"
	TestItemScenario     = "SCENARIO"
	TestItemStep         = "STEP"
	TestItemBeforeClass  = "BEFORE_CLASS"
	TestItemBeforeGroups = "BEFORE_GROUPS"
	TestItemBeforeMethod = "BEFORE_METHOD"
	TestItemBeforeSuite  = "BEFORE_SUITE"
	TestItemBeforeTest   = "BEFORE_TEST"
	TestItemAfterClass   = "AFTER_CLASS"
	TestItemAfterGroups  = "AFTER_GROUPS"
	TestItemAfterMethod  = "AFTER_METHOD"
	TestItemAfterSuite   = "AFTER_SUITE"
	TestItemAfterTest    = "AFTER_TEST"
)

type TestItem struct {
	Id          string
	Name        string
	Description string
	Parent      *TestItem
	Parameters  []struct {
		Key   string
		Value string
	}
	Retry     bool
	StartTime time.Time
	Tags      []string
	Type      string
	UniqueId  string

	launch *Launch
}

// NewTestItem creates new test item
func NewTestItem(launch *Launch, name, description, itemType string, tags []string, parent *TestItem) *TestItem {
	return &TestItem{
		Name:        name,
		Description: description,
		Parent:      parent,
		Tags:        tags,
		Type:        itemType,
		launch:      launch,
	}
}

func (ti *TestItem) Start() error {
	url := fmt.Sprintf("%s/%s/item/%s", ti.launch.client.Endpoint, ti.launch.client.Project, ti.Parent.Id)
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
		Name:        ti.Name,
		Description: ti.Description,
		Tags:        ti.Tags,
		StartTime:   toTimestamp(time.Now()),
		LaunchId:    ti.launch.Id,
		Type:        ti.Type,
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal object %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create POST request to %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", ti.launch.client.Token)
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
		return errors.Wrapf(err, "failed to execute POST request %s", req.URL)
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("failed with status %s", resp.Status)
	}

	v := struct {
		Id       string
		UniqueId string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return errors.Wrapf(err, "failed to decode response from %s", req.URL)
	}
	ti.Id = v.Id
	ti.UniqueId = v.UniqueId
	return nil
}

func (ti *TestItem) Finish(status string) error {
	url := fmt.Sprintf("%s/%s/item/%s", ti.launch.client.Endpoint, ti.launch.client.Project, ti.Id)
	data := struct {
		EndTime int64  `json:"end_time"`
		Status  string `json:"status"`
	}{toTimestamp(time.Now()), status}

	b, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal request data %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create PUT request to %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", ti.launch.client.Token)
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
