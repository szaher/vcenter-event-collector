/*
@author: Saad Zaher
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"net/url"
	"reflect"
	"time"

	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"

)

var (
	begin time.Duration
	end time.Duration
	follow bool
	vcenterUrl string
	insecure bool
	username string
	password string
	eventCount int
)


func init(){
	flag.DurationVar(&begin, "b", 10 * time.Minute, "Start time of events to be streamed")
	flag.DurationVar(&end, "e", 0, "End time of events to be streamed")
	flag.BoolVar(&follow, "f", false, "Follow event stream")
	flag.StringVar(&vcenterUrl, "url", "", "Vcenter URL. i.e. https://localhost/sdk")
	flag.StringVar(&username, "u", "administrator@vsphere.local", "Vcenter Username")
	flag.StringVar(&password, "p", "", "Vcenter password")
	flag.BoolVar(&insecure, "i", true, "Insecure")
	flag.IntVar(&eventCount, "c", 100,"Number of events to fetch every time.")


}

func main() {
	// example use against simulator: go run main.go -b 8h -f
	// example use against vCenter with optional event filters:
	// go run main.go -url $GOVMOMI_URL -insecure $GOVMOMI_INSECURE -b 8h -f VmEvent UserLoginSessionEvent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	flag.Parse()

	if vcenterUrl == ""{
		fmt.Fprintf(os.Stderr, "-url vCenter url is Required\n")
		os.Exit(1)
	}
	if password == "" {
		fmt.Fprintf(os.Stderr, "-p Password is required\n")
		os.Exit(1)
	}

	u, _ :=url.Parse(vcenterUrl)
	u.User = url.UserPassword(username, password)

	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil{
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	c.Login(ctx, u.User)

	m := event.NewManager(c.Client)

	ref := c.ServiceContent.RootFolder

		now, err := methods.GetCurrentTime(ctx, c) // vCenter server time (UTC)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		filter := types.EventFilterSpec{
			EventTypeId: flag.Args(), // e.g. VmEvent
			Entity: &types.EventFilterSpecByEntity{
				Entity:    ref,
				Recursion: types.EventFilterSpecRecursionOptionAll,
			},
			Time: &types.EventFilterSpecByTime{
				BeginTime: types.NewTime(now.Add(begin * -1)),
			},
		}
		if end != 0 {
			filter.Time.EndTime = types.NewTime(now.Add(end * -1))
		}

		collector, err := m.CreateCollectorForEvents(ctx, filter)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		defer collector.Destroy(ctx)

		for {
			events, err := collector.ReadNextEvents(ctx, int32(eventCount))
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}

			if len(events) == 0 {
				if follow {
					time.Sleep(time.Second)
					continue
				}
				break
			}

			for i := range events {
				event := events[i].GetEvent()
				kind := reflect.TypeOf(events[i]).Elem().Name()
				fmt.Printf("%d [%s] [%s] %s\n", event.Key, event.CreatedTime.Format(time.ANSIC), kind, event.FullFormattedMessage)
			}
		}
}