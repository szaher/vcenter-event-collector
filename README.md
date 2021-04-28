# vcenter-event-collector

Streaming events from vcenter

> Print help message 

```shell script
go run main.go -h
Usage of main:
  -b duration
        Start time of events to be streamed (default 10m0s)
  -c int
        Number of events to fetch every time. (default 100)
  -e duration
        End time of events to be streamed
  -f    Follow event stream
  -i    Insecure (default true)
  -p string
        Vcenter password
  -u string
        Vcenter Username (default "administrator@vsphere.local")
  -url string
```

> Fetch events

```shell script
go run main.go  -url https://VCENTER_URL/sdk -p VCENTER_PASSWORD   
```


> Stream events in certain period

```shell script
go run main.go  -url https://VCENTER_URL/sdk -p VCENTER_PASSWORD -b 10h -e 9h  
```

> Stream events

```shell script
go run main.go  -url https://VCENTER_URL/sdk -p VCENTER_PASSWORD -f  
```
