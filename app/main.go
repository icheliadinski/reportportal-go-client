package main

import (
	"fmt"

	"github.com/icheliadinski/reportportal-go-client/app/reportportal"
)

func main() {
	c, err := reportportal.NewClient("", "")
	if err := c.CheckConnect(); err != nil {
		fmt.Println(err)
	}
}
