package castv2

import (
	"net"
	"time"

	"github.com/AndreasAbdi/go-castv2/controllers"
	"github.com/AndreasAbdi/go-castv2/primitives"
)

const defaultTimeout = time.Second * 10

//Device Object to run basic chromecast commands
type Device struct {
	client               *primitives.Client
	heartbeatController  *controllers.HeartbeatController
	connectionController *controllers.ConnectionController
	receiverController   *controllers.ReceiverController
	mediaController      *controllers.MediaController
}

//NewDevice is constructor for Device struct
func NewDevice(host net.IP, port int) (Device, error) {
	var device Device

	client, err := primitives.NewClient(host, port)
	if err != nil {
		return device, err
	}
	device.client = client

	device.heartbeatController = controllers.NewHeartbeatController(client, defaultChromecastSenderID, defaultChromecastSenderID)
	device.heartbeatController.Start()

	device.connectionController = controllers.NewConnectionController(client, defaultChromecastSenderID, defaultChromecastSenderID)
	device.connectionController.Connect()

	device.receiverController = controllers.NewReceiverController(client, defaultChromecastSenderID, defaultChromecastSenderID)
	device.mediaController = controllers.NewMediaController(client, defaultChromecastSenderID, defaultChromecastSenderID)

	return device, nil
}

//PlayMedia plays a video via the media controller.
func (device *Device) PlayMedia(URL string, MIMEType string) {
	appID := mediaReceiverAppID
	device.receiverController.LaunchApplication(&appID, defaultTimeout, false)
	device.mediaController.Load(URL, MIMEType, defaultTimeout)
}
