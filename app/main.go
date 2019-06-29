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
	fmt.Println("Connection checked. Trying start a launch...")

	id, err := c.StartLaunch("Vika", "Horoshaya", nil, time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println("Launch Started")
	if err := c.StopLaunch(id, "", time.Now()); err != nil {
		panic(err)
	}
	fmt.Println("Launch Stopped")
}
