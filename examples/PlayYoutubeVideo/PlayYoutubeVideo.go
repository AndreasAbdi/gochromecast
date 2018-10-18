package main

import (
	"time"

	"github.com/AndreasAbdi/gochromecast"
	"github.com/AndreasAbdi/gochromecast/configs"
)

// A simple example, showing how to play a youtube video.
func main() {
	devices := make(chan *cast.Device, 100)
	cast.FindDevices(time.Second*5, devices)
	for device := range devices {
		appID := configs.YoutubeAppID
		device.ReceiverController.LaunchApplication(&appID, time.Second*5, false)
		device.YoutubeController.PlayVideo("F1B9Fk_SgI0", "")
	}

}
