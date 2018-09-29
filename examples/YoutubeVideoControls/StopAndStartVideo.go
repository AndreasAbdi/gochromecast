package main

import (
	"time"

	castv2 "github.com/AndreasAbdi/go-castv2"
	"github.com/AndreasAbdi/go-castv2/configs"
)

// A simple example, showing how to control youtube videos.
func main() {
	devices := make(chan *castv2.Device, 100)
	castv2.FindDevices(time.Second*5, devices)
	for device := range devices {
		appID := configs.YoutubeAppID
		device.ReceiverController.LaunchApplication(&appID, time.Second*5, false)
		device.YoutubeController.PlayVideo("F1B9Fk_SgI0")
		time.Sleep(time.Second * 20)
		device.MediaController.Seek(float64(0), time.Second*5)
		//time.Sleep(time.Second * 5)
		// device.QuitApplication(time.Second * 5)
		// time.Sleep(time.Second * 2)
		// device.YoutubeController.AddToQueue("0q-aR6XNZDg")
		// time.Sleep(time.Second * 2)
		// device.YoutubeController.RemoveFromQueue("rEq1Z0bjdwc")
		// device.MediaController.Pause(time.Second * 10)
	}

}
