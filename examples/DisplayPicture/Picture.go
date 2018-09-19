package main

import (
	"fmt"
	"time"

	"github.com/AndreasAbdi/go-castv2"
)

// A simple example, showing how to create a device and use it.
func main() {
	deviceCh := make(chan *castv2.Device, 100)
	castv2.FindDevices(time.Second*30, deviceCh)
	for device := range deviceCh {
		fmt.Print(device)
		device.PlayMedia(
			"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/BigBuckBunny.jpg",
			"image/jpeg")
	}

}
