package main

import (
	"time"

	castv2 "github.com/AndreasAbdi/go-castv2"
)

// A simple example, showing how to find a Chromecast using mdns, and request its status.
func main() {
	devices := make(chan *castv2.Device, 100)
	castv2.FindDevices(time.Second*5, devices)
	for device := range devices {
		device.PlayMedia("http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4", "video/mp4")
		time.Sleep(time.Second * 5)
		device.MediaController.Pause(time.Second * 5)
		time.Sleep(time.Second * 5)
		device.MediaController.Play(time.Second * 5)
		device.QuitApplication(time.Second * 5)
	}
}
