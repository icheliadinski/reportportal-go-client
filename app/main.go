package main

import (
	"fmt"
	"os"
	"time"

	"github.com/icheliadinski/reportportal-go-client/app/rp"
	"github.com/jessevdk/go-flags"
)

func main() {
	var opts rp.Client
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}
	c := rp.NewClient(opts.Endpoint, opts.Token, opts.Launch, opts.Project)
	if err := c.CheckConnect(); err != nil {
		panic(err)
	}
	fmt.Println("Connection checked. Trying to start a launch...")
	time.Sleep(2 * time.Second)

	id, err := c.StartLaunch("Go Launch", "Test go launch", rp.ModeDefault, []string{}, time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println("Launch Started. Updating...")
	time.Sleep(2 * time.Second)

	if err := c.UpdateLaunch(id, "updated descr", rp.ModeDebug, []string{"tag1", "tag2"}); err != nil {
		panic(err)
	}
	fmt.Println("Launch updated. Adding item...")
	time.Sleep(2 * time.Second)

	suiteId, err := c.StartTestItem(id, "Item", "Item descr", "SUITE", "", []string{"suite"}, time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println("Item started. Creating subitem...")
	time.Sleep(2 * time.Second)

	subItem, err := c.StartTestItem(id, "Sub Item", "Sub item descr", "TEST", suiteId, []string{"sub"}, time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println("Subitem started. Writing to test item...")
	time.Sleep(2 * time.Second)

	if err := c.Log(subItem, "my super failed message", rp.LevelError, time.Now()); err != nil {
		panic(err)
	}
	fmt.Println("Log sent. Adding file...")
	time.Sleep(2 * time.Second)

	// f, err := os.Open("C:\\Users\\Igor_Cheliadinski\\Downloads\\img.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// a := &reportportal.Attachment{
	// 	Name: "img.jpg",
	// 	Type: "image/jpeg",
	// 	Content: ioutil.ReadFile(f)
	// }
	// c.LogWithFile(subItem, "my", reportportal.LevelError, a, time.Now())
	if err := c.LogWithFile(subItem, "my", "INFO", "C:\\Users\\Igor_Cheliadinski\\Downloads\\img.JPEG", time.Now()); err != nil {
		panic(err)
	}
	fmt.Println("File sent. Failing subitem...")
	time.Sleep(2 * time.Second)

	if err := c.FinishTestItem(subItem, rp.StatusFailed, time.Now()); err != nil {
		panic(err)
	}
	fmt.Println("Subitem failed. Stopping launch...")
	time.Sleep(2 * time.Second)

	if err := c.StopLaunch(id, "", time.Now()); err != nil {
		panic(err)
	}
	fmt.Println("Launch Stopped. Getting project settings...")
	time.Sleep(2 * time.Second)

	s, err := c.GetProjectSettings()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	fmt.Println("Project settings received. Closing...")
}
