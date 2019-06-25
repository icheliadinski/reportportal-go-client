package main

import (
	"fmt"
	"os"

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
	fmt.Println("Success")
}
