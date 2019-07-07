package reportportal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
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

	LevelError = "error"
	LevelTrace = "trace"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelEmpty = ""
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

type LogFile struct {
	File            bytes.Buffer `json:"file"`
	JsonRequestPart bytes.Buffer `json:"json_request_part"`
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
		StartTime:   startTime.Unix() * int64(time.Microsecond),
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
		StartTime:   startTime.Unix() * int64(time.Microsecond),
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
		EndTime: endTime.Unix() * int64(time.Microsecond),
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

// Log sends a log to report portal server
func (c *Client) Log(id, message, level string, startTime time.Time) error {
	url := fmt.Sprintf("%s/%s/log", c.Endpoint, c.Project)
	data := struct {
		ItemId  string `json:"item_id"`
		Time    int64  `json:"time"`
		Message string `json:"message"`
		Level   string `json:"level"`
	}{
		ItemId:  id,
		Time:    startTime.Unix() * int64(time.Microsecond),
		Message: message,
		Level:   level,
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal object, %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return errors.Wrapf(err, "failed to create POST request to %s", url)
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
		return errors.Wrapf(err, "failed to execute POST request %s", req.URL)
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}

// LogWithFile sends log with file as attachment
func (c *Client) LogWithFile(id, message, level string, startTime time.Time) error {
	url := fmt.Sprintf("%s/%s/log", c.Endpoint, c.Project)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	attachFile, err := os.Open("C:\\Users\\Igor_Cheliadinski\\Downloads\\img.jpg")
	if err != nil {
		return errors.Wrap(err, "failed to open file img.png")
	}
	defer attachFile.Close()

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "img.jpg"))
	h.Set("Content-Type", "image/jpeg")
	filePart, err := writer.CreatePart(h)
	if err != nil {
		return errors.Wrap(err, "failed to create form file img.png")
	}

	jsonReqFile, err := os.Create("json_request_part.json")
	if err != nil {
		return errors.Wrap(err, "failed to create filename json_request_part.json")
	}
	defer jsonReqFile.Close()

	name := struct {
		Name string `json:"name"`
	}{"img.jpg"}
	jsonReqData := struct {
		File struct {
			Name string `json:"name"`
		} `json:"file"`
		ItemId   string `json:"item_id"`
		LogLevel string `json:"level"`
		Message  string `json:"message"`
		Time     int64  `json:"time"`
	}{
		File:     name,
		ItemId:   id,
		LogLevel: level,
		Message:  message,
		Time:     startTime.Unix() * int64(time.Microsecond),
	}

	b, err := json.Marshal(&jsonReqData)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json req data")
	}

	_, err = jsonReqFile.Write(b)
	if err != nil {
		return errors.Wrap(err, "failed to write to json_request_part.json")
	}

	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "json_request_part", "json_request_part.json"))
	h.Set("Content-Type", "application/json")
	jsonReqPart, err := writer.CreatePart(h)
	if err != nil {
		return errors.Wrap(err, "failed to create form file for json_request_part.json")
	}

	_, err = io.Copy(filePart, attachFile)
	_, err = io.Copy(jsonReqPart, jsonReqFile)
	defer writer.Close()
	if err != nil {
		return errors.Wrap(err, "failed to copy file content to body")
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return errors.Wrapf(err, "failed to create POST request %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", c.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	bbb, _ := ioutil.ReadAll(req.Body)
	log.Println(string(bbb))
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
		bb, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(bb))
		return errors.Errorf("failed with status %s", resp.Status)
	}
	return nil
}

// finalizeLaunch finalizes exact match with specific action
func (c *Client) finalizeLaunch(id, action, status string, endTime time.Time) error {
	url := fmt.Sprintf("%s/%s/launch/%s/%s", c.Endpoint, c.Project, id, action)
	data := struct {
		Status  string `json:"status"`
		EndTime int64  `json:"end_time"`
	}{
		Status:  status,
		EndTime: endTime.Unix() * int64(time.Microsecond),
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
