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

// type LogFile struct {
// 	File            bytes.Buffer `json:"file"`
// 	JsonRequestPart bytes.Buffer `json:"json_request_part"`
// }

// ProjectSettings defines project settings
// type ProjectSettings struct {
// 	StatisticsStrategy string                 `json:"statisticsStrategy"`
// 	Name               string                 `json:"project"`
// 	SubTypes           map[string]interface{} `json:"subTypes"`
// }

// type Attachment struct {
// 	Name    string
// 	Type    string
// 	Content []byte
// }

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

// // Log sends a log to report portal server
// func (c *Client) Log(id, message, level string, startTime time.Time) error {
// 	url := fmt.Sprintf("%s/%s/log", c.Endpoint, c.Project)
// 	data := struct {
// 		ItemId  string `json:"item_id"`
// 		Time    int64  `json:"time"`
// 		Message string `json:"message"`
// 		Level   string `json:"level"`
// 	}{
// 		ItemId:  id,
// 		Time:    startTime.Unix() * int64(time.Microsecond),
// 		Message: message,
// 		Level:   level,
// 	}

// 	b, err := json.Marshal(&data)
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to marshal object, %v", data)
// 	}

// 	r := bytes.NewReader(b)
// 	req, err := http.NewRequest(http.MethodPost, url, r)
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to create POST request to %s", url)
// 	}

// 	auth := fmt.Sprintf("Bearer %s", c.Token)
// 	req.Header.Set("Authorization", auth)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	defer func() {
// 		if err := resp.Body.Close(); err != nil {
// 			log.Println("[WARN] failed to close response body")
// 		}
// 	}()
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to execute POST request %s", req.URL)
// 	}
// 	if resp.StatusCode != http.StatusCreated {
// 		return errors.Errorf("failed with status %s", resp.Status)
// 	}
// 	return nil
// }

// // LogWithFile sends log with file as attachment
// func (c *Client) LogWithFile(id, message, level, filename string, startTime time.Time) error {
// 	url := fmt.Sprintf("%s/%s/log", c.Endpoint, c.Project)
// 	bodyBuf := &bytes.Buffer{}
// 	bodyWriter := multipart.NewWriter(bodyBuf)

// 	// json request part
// 	h := make(textproto.MIMEHeader)
// 	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, "json_request_part"))
// 	h.Set("Content-Type", "application/json")
// 	reqWriter, err := bodyWriter.CreatePart(h)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to create form file")
// 	}

// 	s := fmt.Sprintf(`[{"file":{"name": "%s"}, "item_id": "%s", "level":"%s", "message": "%s", "time": %d}]`, filepath.Base(filename), id, level, message, startTime.Unix()*int64(time.Microsecond))
// 	reqReader := strings.NewReader(s)

// 	_, err = io.Copy(reqWriter, reqReader)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to copy reader")
// 	}

// 	// file
// 	h = make(textproto.MIMEHeader)
// 	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(filename)))
// 	h.Set("Content-Type", "img/jpeg")
// 	fileWriter, err := bodyWriter.CreatePart(h)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to create form file")
// 	}

// 	fh, err := os.Open(filename)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to open file")
// 	}
// 	defer fh.Close()

// 	_, err = io.Copy(fileWriter, fh)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to copy file writer")
// 	}

// 	bodyWriter.Close()

// 	req, err := http.NewRequest(http.MethodPost, url, bodyBuf)
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to create POST request to %s", url)
// 	}

// 	auth := fmt.Sprintf("Bearer %s", c.Token)
// 	req.Header.Set("Authorization", auth)
// 	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	defer func() {
// 		if err := resp.Body.Close(); err != nil {
// 			log.Println("[WARN] failed to close response body")
// 		}
// 	}()
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to execute POST request %s", req.URL)
// 	}
// 	if resp.StatusCode != http.StatusCreated {
// 		return errors.Errorf("failed with status %s", resp.Status)
// 	}
// 	return nil
// }

// // finalizeLaunch finalizes exact match with specific action
// func (c *Client) finalizeLaunch(id, action, status string, endTime time.Time) error {
// 	url := fmt.Sprintf("%s/%s/launch/%s/%s", c.Endpoint, c.Project, id, action)
// 	data := struct {
// 		Status  string `json:"status"`
// 		EndTime int64  `json:"end_time"`
// 	}{
// 		Status:  status,
// 		EndTime: endTime.Unix() * int64(time.Microsecond),
// 	}

// 	b, err := json.Marshal(&data)
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to marshal object, %v", data)
// 	}

// 	r := bytes.NewReader(b)
// 	req, err := http.NewRequest(http.MethodPut, url, r)
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to create request to %s", url)
// 	}

// 	auth := fmt.Sprintf("Bearer %s", c.Token)
// 	req.Header.Set("Authorization", auth)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	defer func() {
// 		if err := resp.Body.Close(); err != nil {
// 			log.Println("[WARN] Failed to close response body")
// 		}
// 	}()
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to execute PUT request %s", req.URL)
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return errors.Errorf("failed with status %s", resp.Status)
// 	}
// 	return nil
// }
