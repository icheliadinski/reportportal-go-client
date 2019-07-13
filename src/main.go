package main

import (
	"fmt"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/icheliadinski/reportportal-go-client/src/rp"
)

func main() {
	var opts rp.Client
	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}

	c := rp.NewClient(opts.Endpoint, opts.Project, opts.Token)
	c.CheckConnect()
	fmt.Println("Connection checked!")

	fmt.Println("Trying to start a launch...")
	time.Sleep(2 * time.Second)
	l := rp.NewLaunch(c, "Go Launch", "Go launch info", rp.ModeDefault, []string{"tag1", "tag2"})
	if err := l.Start(); err != nil {
		panic(err)
	}
	fmt.Println("Launch started!")

	fmt.Println("Trying to add test suite...")
	time.Sleep(2 * time.Second)
	ts := rp.NewTestItem(l, "Suite", "Suite descr", rp.TestItemSuite, []string{"suite"}, nil)
	if err := ts.Start(); err != nil {
		panic(err)
	}
	fmt.Println("Test suite created!")

	fmt.Println("Trying to add test item...")
	time.Sleep(2 * time.Second)
	ti := rp.NewTestItem(l, "Test", "Test descr", rp.TestItemTest, []string{"test"}, ts)
	if err := ti.Start(); err != nil {
		panic(err)
	}
	fmt.Println("Test item created!")

	fmt.Println("Trying to send log...")
	time.Sleep(2 * time.Second)
	if err := ti.Log("super mega message", rp.LevelInfo, ""); err != nil {
		panic(err)
	}
	fmt.Println("Log sent!")

	fmt.Println("Trying to send log with attach...")
	time.Sleep(2 * time.Second)
	if err := ti.Log("super mega message", rp.LevelError, "C:\\Users\\Igor_Cheliadinski\\Downloads\\cat.jpg"); err != nil {
		panic(err)
	}
	fmt.Println("Log sent!")

	fmt.Println("Trying to finish test...")
	time.Sleep(2 * time.Second)
	if err := ti.Finish(rp.StatusFailed); err != nil {
		panic(err)
	}
	fmt.Println("Test finished!")

	fmt.Println("Trying to finish suite...")
	time.Sleep(2 * time.Second)
	if err := ts.Finish(rp.StatusFailed); err != nil {
		panic(err)
	}
	fmt.Println("Suite finished!")

	fmt.Println("Tryng to stop the launch...")
	time.Sleep(2 * time.Second)
	if err := l.Stop(rp.StatusFailed); err != nil {
		panic(err)
	}
	fmt.Println("Launch stopped!")
}
