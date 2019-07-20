# ReportPortal Go Client [![Build Status](https://travis-ci.org/icheliadinski/client-go.svg?branch=master)](https://travis-ci.org/icheliadinski/client-go) [![Go Report Card](https://goreportcard.com/badge/github.com/icheliadinski/client-go)](https://goreportcard.com/report/github.com/icheliadinski/client-go) [![Coverage Status](https://coveralls.io/repos/github/icheliadinski/client-go/badge.svg?branch=master)](https://coveralls.io/github/icheliadinski/client-go?branch=master) [![GoDoc](https://godoc.org/github.com/icheliadinski/client-go/rp?status.svg)](https://godoc.org/github.com/icheliadinski/client-go/rp)
Go client for ReportPortal http://reportportal.io/

## Already implemented listeners:
* EMPTY


## Installation
The latest version is available with:
```cmd
go get github.com/icheliadinski/client-go/rp
```

## Example
```go
	c := rp.NewClient("your rp endpoint", "project name", "secret token", 1)
	if err := c.CheckConnect(); err != nil {
		// handle error
	}
```

## Settings
When creating a client instance, you need to specify the following parameters:

Parameter | Description
--------- | -----------
endpoint  | URL of your server. For example, if you visit the page at 'https://server:8080/ui', then endpoint will be equal to 'https://server:8080/api/v1'
project   | The name of the project in which the launches will be created.
token     | user's token Report Portal from which you want to send requests. It can be found on the profile page of this user.

## Api

### CheckConnect
 CheckConnect - method for verifying the correctness of the client connection
```go
if err := c.CheckConnect(); err != nil {
  // handle error
}
```
