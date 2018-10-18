package main

import (
	"time"

	"github.com/AndreasAbdi/gochromecast"
)

// A simple example, showing how to create a device and use it to display a picture.
func main() {
	deviceCh := make(chan *cast.Device, 100)
	cast.FindDevices(time.Second*30, deviceCh)
	for device := range deviceCh {
		device.PlayMedia(
			"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/BigBuckBunny.jpg",
			"image/jpeg")
	}

}
