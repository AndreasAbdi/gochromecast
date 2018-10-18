package main

import (
	"time"

	"github.com/AndreasAbdi/gochromecast"
	"github.com/AndreasAbdi/gochromecast/configs"
)

// A simple example, showing how to control youtube videos.
func main() {
	devices := make(chan *cast.Device, 100)
	cast.FindDevices(time.Second*5, devices)
	for device := range devices {
		appID := configs.YoutubeAppID
		device.ReceiverController.LaunchApplication(&appID, time.Second*5, false)

		device.YoutubeController.PlayVideo("F1B9Fk_SgI0", "")
		time.Sleep(time.Second * 20)

		device.YoutubeController.AddToQueue("0q-aR6XNZDg")
		time.Sleep(time.Second * 2)

		device.YoutubeController.RemoveFromQueue("0q-aR6XNZDg")
		device.MediaController.Seek(float64(0), time.Second*5)

		time.Sleep(time.Second * 2)
		device.MediaController.Pause(time.Second * 10)
	}
}
