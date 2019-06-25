package main

import (
	"fmt"

	"github.com/icheliadinski/reportportal-go-client/app/reportportal"
)

func main() {
	c := reportportal.Client{
		Endpoint: "",
		Token:    "",
	}
	if err := c.CheckConnect(); err != nil {
		fmt.Println(err)
	}
}
