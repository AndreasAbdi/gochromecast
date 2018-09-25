package castv2

import (
	"log"
	"net"
	"time"

	"github.com/AndreasAbdi/go-castv2/configs"
	"github.com/AndreasAbdi/go-castv2/controllers"
	"github.com/AndreasAbdi/go-castv2/primitives"
	"github.com/davecgh/go-spew/spew"
)

const defaultTimeout = time.Second * 10

//Device Object to run basic chromecast commands
type Device struct {
	client               *primitives.Client
	heartbeatController  *controllers.HeartbeatController
	connectionController *controllers.ConnectionController
	ReceiverController   *controllers.ReceiverController
	MediaController      *controllers.MediaController
}

//NewDevice is constructor for Device struct
func NewDevice(host net.IP, port int) (Device, error) {
	var device Device

	client, err := primitives.NewClient(host, port)
	if err != nil {
		log.Fatalf("Failed to connect to chromecast %s", host)
		return device, err
	}
	device.client = client

	device.heartbeatController = controllers.NewHeartbeatController(client, defaultChromecastSenderID, defaultChromecastReceiverID)
	device.heartbeatController.Start()

	device.connectionController = controllers.NewConnectionController(client, defaultChromecastSenderID, defaultChromecastReceiverID)
	device.connectionController.Connect()

	device.ReceiverController = controllers.NewReceiverController(client, defaultChromecastSenderID, defaultChromecastReceiverID)

	device.MediaController = controllers.NewMediaController(client, defaultChromecastSenderID, device.ReceiverController)
	return device, nil
}

//Play just plays.
func (device *Device) Play() {
	device.MediaController.Play(defaultTimeout)
}

//PlayMedia plays a video via the media controller.
func (device *Device) PlayMedia(URL string, MIMEType string) {
	appID := configs.MediaReceiverAppID
	device.ReceiverController.LaunchApplication(&appID, defaultTimeout, false)
	device.MediaController.Load(URL, MIMEType, defaultTimeout)
}

func (device *Device) QuitApplication(timeout time.Duration) {
	status, err := device.ReceiverController.GetStatus(timeout)
	if err != nil {
		spew.Dump("Failed to quit application", err)
	}
	for _, appSessions := range status.Applications {
		session := appSessions.SessionID
		device.ReceiverController.StopApplication(session, timeout)
	}
}

//PlayMedia plays a video via the media controller.
func (device *Device) TestYoutube(URL string) {
	appID := configs.SpotifyAppID
	device.ReceiverController.LaunchApplication(&appID, defaultTimeout, false)
}

func (device *Device) GetMediaStatus(timeout time.Duration) {
	response, err := device.MediaController.GetStatus(time.Second * 5)
	spew.Dump("Status Media response", response, err)

}

func (device *Device) GetStatus(timeout time.Duration) {
	response, err := device.ReceiverController.GetStatus(time.Second * 5)
	spew.Dump("Status response", response, err)

}
