package main

import (
	"time"

	castv2 "github.com/AndreasAbdi/go-castv2"
)

// A simple example on how to use the devices and
func main() {
	devices := make(chan *castv2.Device, 100)
	castv2.FindDevices(time.Second*5, devices)
	for device := range devices {
		device.PlayMedia("http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4", "video/mp4")
		device.QuitApplication(time.Second * 5)
	}
}
