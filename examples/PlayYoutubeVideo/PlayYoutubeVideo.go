package main

import (
	"time"

	castv2 "github.com/AndreasAbdi/go-castv2"
	"github.com/AndreasAbdi/go-castv2/configs"
)

// A simple example, showing how to play a youtube video.
func main() {
	devices := make(chan *castv2.Device, 100)
	castv2.FindDevices(time.Second*5, devices)
	for device := range devices {
		appID := configs.YoutubeAppID
		device.ReceiverController.LaunchApplication(&appID, time.Second*5, false)
		device.YoutubeController.PlayVideo("cwQgjq0mCdE")
	}

}
