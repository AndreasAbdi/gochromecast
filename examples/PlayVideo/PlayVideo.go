package main

import (
	"time"

	"github.com/AndreasAbdi/gochromecast"
)

// A simple example on how to use the devices and play a video.
func main() {
	devices := make(chan *cast.Device, 100)
	cast.FindDevices(time.Second*5, devices)
	for device := range devices {
		device.PlayMedia("http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4", "video/mp4")
		device.QuitApplication(time.Second * 5)
	}
}
