
# gochromecast

[![GoDoc](https://godoc.org/github.com/AndreasAbdi/gochromecast?status.png)](http://godoc.org/github.com/AndreasAbdi/gochromecast)
[![Go Report Card](https://goreportcard.com/badge/github.com/AndreasAbdi/gochromecast)](https://goreportcard.com/report/github.com/AndreasAbdi/gochromecast)
[![Build Status](https://travis-ci.org/AndreasAbdi/gochromecast.svg?branch=master)](https://travis-ci.org/AndreasAbdi/gochromecast)

## Description

Library for using chromecast. Contains partial implementations for media player controls, and youtube controls.

## Usage

To install the library, run
`go get github.com/AndreasAbdi/gochromecast`

## Examples

```go
// A simple example, showing how to find a Chromecast using mdns, and request its status.
package main

import (
"time"
"github.com/AndreasAbdi/gochromecast"
)

func main() {
    devices := make(chan *cast.Device, 100)
    cast.FindDevices(time.Second*5, devices)
    for device := range devices {
        device.PlayMedia("http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4", "video/mp4")
        time.Sleep(time.Second * 5)
        device.MediaController.Pause(time.Second * 5)
        device.QuitApplication(time.Second * 5)
    }
}
```

For more examples of how to use it, see the [examples folder](https://github.com/AndreasAbdi/gochromecast/tree/master/examples).

## References

References listed in docs.
Ported from

- [node-castv2](https://github.com/thibauts/node-castv2)
- [node-castv2-client](https://github.com/thibauts/node-castv2-client)
- [go-castv2](https://github.com/ninjasphere/go-castv2)