package rp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
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

// TestItem defines test item structure
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

// Attachment defines file attachment structure
type Attachment struct {
	Name    string
	Type    string
	Content []byte
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

func (ti *TestItem) Log(message, level string, attachment *Attachment) error {
	var req *http.Request
	var err error
	if attachment != nil {
		req, err = ti.getReqForLogWithAttach(message, level, attachment)
	} else {
		req, err = ti.getReqForLog(message, level)
	}
	if err != nil {
		return err
	}

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

func (ti *TestItem) getReqForLogWithAttach(message, level string, attachment *Attachment) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s/log", ti.launch.client.Endpoint, ti.launch.client.Project)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// json request part
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="json_request_part"`)
	h.Set("Content-Type", "application/json")
	reqWriter, err := bodyWriter.CreatePart(h)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create form file")
	}

	s := fmt.Sprintf(`[{"file":{"name": "%s"}, "item_id": "%s", "level":"%s", "message": "%s", "time": %d}]`, filepath.Base(attachment.filename), id, level, message, startTime.Unix()*int64(time.Microsecond))
	reqReader := strings.NewReader(s)

	// _, err = io.Copy(reqWriter, reqReader)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to copy reader")
	// }

	// // file
	// h = make(textproto.MIMEHeader)
	// h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(filename)))
	// h.Set("Content-Type", "img/jpeg")
	// fileWriter, err := bodyWriter.CreatePart(h)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to create form file")
	// }

	// fh, err := os.Open(filename)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to open file")
	// }
	// defer fh.Close()

	// _, err = io.Copy(fileWriter, fh)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to copy file writer")
	// }

	// bodyWriter.Close()

	// req, err := http.NewRequest(http.MethodPost, url, bodyBuf)
	// if err != nil {
	// 	return errors.Wrapf(err, "failed to create POST request to %s", url)
	// }

	// auth := fmt.Sprintf("Bearer %s", c.Token)
	// req.Header.Set("Authorization", auth)
	// req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	// return nil, nil
}

func (ti *TestItem) getReqForLog(message, level string) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s/log", ti.launch.client.Endpoint, ti.launch.client.Project)
	data := struct {
		ItemId  string `json:"item_id"`
		Message string `json:"message"`
		Level   string `json:"level"`
		Time    int64  `json:"time"`
	}{ti.Id, message, level, toTimestamp(time.Now())}

	b, err := json.Marshal(&data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal object, %v", data)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create POST request to %s", url)
	}

	auth := fmt.Sprintf("Bearer %s", ti.launch.client.Token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
