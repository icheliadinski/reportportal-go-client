package main

import (
	"fmt"
	"os"
	"time"

	"github.com/icheliadinski/reportportal-go-client/app/reportportal"
	"github.com/jessevdk/go-flags"
)

func main() {
	var opts reportportal.Client
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}
	c := reportportal.NewClient(opts.Endpoint, opts.Token, opts.Launch, opts.Project)
	if err := c.CheckConnect(); err != nil {
		panic(err)
	}
	fmt.Println("Connection checked. Trying to start a launch...")
	time.Sleep(2 * time.Second)

	_, err := c.StartLaunch("Go Launch", "Test go launch", reportportal.ModeDefault, []string{}, time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println("Launch Started. Updating...")
	time.Sleep(2 * time.Second)

	if err := c.UpdateLaunch(id, "updated descr", reportportal.ModeDebug, []string{"tag1", "tag2"}); err != nil {
		panic(err)
	}
	fmt.Println("Launch updated. Stopping...")
	time.Sleep(2 * time.Second)

	if err := c.StopLaunch(id, "", time.Now()); err != nil {
		panic(err)
	}
	fmt.Println("Launch Stopped")
}
