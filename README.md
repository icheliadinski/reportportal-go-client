# ReportPortal Go Client [![Build Status](https://travis-ci.org/icheliadinski/client-go.svg?branch=master)](https://travis-ci.org/icheliadinski/client-go) [![Go Report Card](https://goreportcard.com/badge/github.com/icheliadinski/client-go)](https://goreportcard.com/report/github.com/icheliadinski/client-go) [![Coverage Status](https://coveralls.io/repos/github/icheliadinski/client-go/badge.svg?branch=master)](https://coveralls.io/github/icheliadinski/client-go?branch=master) [![GoDoc](https://godoc.org/github.com/icheliadinski/client-go/rp?status.svg)](https://godoc.org/github.com/icheliadinski/client-go/rp)
Go client for ReportPortal http://reportportal.io/

## Already implemented listeners:
* EMPTY


## Installation
The latest version is available with:
```cmd
go get github.com/icheliadinski/client-go
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
endpoint  | URL of your RP server.
project   | The name of the project in which the launches will be created.
token     | user's token Report Portal from which you want to send requests. It can be found on the profile page of this user.
version   | API version. Responsible for adding /v1 or /v2 etc to the API endpoint

## Api

### Client

#### CheckConnect
 CheckConnect - method for verifying the correctness of the client connection
```go
if err := c.CheckConnect(); err != nil {
  // handle error
}
```

### Launch

#### NewLaunch
 NewLaunch - creates new launch object. Returns this object
```go
l := rp.NewLaunch(c, "Launch name", "Description", rp.ModeDefault, []string{"tag1", "tag2"})
```

Parameter   | Description
----------- | -----------
client      | ReportPortal client created by NewClient function
name        | Launch name
description | Launch description
mode        | Launch mode (rp.ModeDefault or rp.ModeDebug)
tags        | (optional) Tags list for the launch

#### Start
 Start - starts spcified launch object. Returns error
```go
if err := l.Start(); err != nil {
  // handle error
}
```

#### Finish
 Finish - finishes specified launch object. Returns error
```go
if err := l.Finish(rp.StatusPassed); err != nil {
  // handle error
}
```

Parameter | Description
--------- | -----------
status    | Status with which one launch should be finished (all statuses accessible with `rp.Status...` constant)

#### Stop
 Stop - stops specified launch object. Returns error
```go
if err := l.Stop(rp.StatusFailed); err != nil {
  // handle error
}
```

Parameter | Description
--------- | -----------
status    | Status with which one launch should be stopped (all statuses accessible with `rp.Status...` constant)


#### Delete
 Delete - deletes specified launch object. Returns error
```go
if err := l.Delete(); err != nil {
  // handle error
}
```

#### Update
 Delete - deletes specified launch object. Returns error
```go
if err := l.Update("new description", rp.ModeDebug, []string{"new", "tags"}); err != nil {
  // handle error
}
```

Parameter    | Description
------------ | -----------
descsription | New launch description
mode         | New launch mode (all modes accessible with `rp.Mode...` constant)
tags         | New launch tags

### TestItem

#### NewTestItem
 NewTestItem - creates new test item object. Returns this object
```go
ti := rp.NewTestItem(launch, "Test item", "child description", rp.TestItemTest, []string{"tag1", "tag2"}, nil)
```

Parameter   | Description
----------- | -----------
launch      | ReportPortal launch created by NewLaunch function
name        | Test item name
description | Test item description
type        | Test item type (all types accessible through `rp.TestItem...` constants)
tags        | (optional) Tags list for the test item
