package rp

import (
	"time"
)

// Launch defines launch info
type Launch struct {
	Id          string
	Name        string
	Description string
	Mode        string
	StartTime   time.Time
	Tags        []string

	client *Client
}

// NewLaunch creates new launch for specified client
func NewLaunch(client *Client, name, description, mode string, tags []string) *Launch {
	return &Launch{
		Name:        name,
		Description: description,
		Mode:        mode,
		Tags:        tags,
		client:      client,
	}
}

// Start starts the launch
func (l *Launch) Start() error {
	return nil
}

// Stop stops the launch
func (l *Launch) Stop() error {
	return nil
}

// Finish finishes launch
func (l *Launch) Finish() error {
	return nil
}

// Delete delete launch
func (l *Launch) Delete() error {
	return nil
}

// Update updates launch
func (l *Launch) Update() error {
	return nil
}

// Compare compares with specified launch
func (l *Launch) Compare(launchToCompare *Launch) error {
	return nil
}

// Analyze launches auto-analyzer on demand
func (l *Launch) Analyze() error {
	return nil
}

// Import imports junit xml report
func (l *Launch) Import() error {
	return nil
}

// Latest gets list of latest project launches by filter
func (l *Launch) Latest() error {
	return nil
}
