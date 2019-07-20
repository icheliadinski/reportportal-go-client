package main

import "github.com/icheliadinski/client-go/rp"

func main() {
	c := rp.NewClient("your rp endpoint", "project name", "secret token", 1)
	l := rp.NewLaunch(c, "Launch name", "Description", rp.ModeDefault, []string{"tag1", "tag2"})

	if err := c.CheckConnect(); err != nil {
		// handle error
	}

	if err := l.Start(); err != nil {
		// handle error
	}

	parent := rp.NewTestItem(l, "Suite item", "parent description", rp.TestItemSuite, []string{"tag1", "tag2"}, nil)
	if err := parent.Start(); err != nil {
		// handle error
	}

	child := rp.NewTestItem(l, "Test item", "child description", rp.TestItemTest, []string{"tag1", "tag2"}, parent)
	if err := child.Start(); err != nil {
		// handle error
	}

	if err := child.Finish(rp.StatusPassed); err != nil {
		// handle error
	}
	if err := parent.Finish(rp.StatusPassed); err != nil {
		// handle error
	}

	if err := l.Finish(rp.StatusPassed); err != nil {
		// handle error
	}
}
