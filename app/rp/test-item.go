package rp

import "time"

type TestItem struct {
	Id          string
	Name        string
	Description string
	Parent      *TestItem
	Parameters  []map[string]string
	Retry       bool
	StartTime   time.Time
	Tags        []string
	Type        string
	UniqueId    string

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
	return nil
}

func (ti *TestItem) Finish() error {
	return nil
}
