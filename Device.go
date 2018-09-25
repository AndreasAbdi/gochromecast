package castv2

import (
	"log"
	"net"
	"time"

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
	receiverController   *controllers.ReceiverController
	mediaController      *controllers.MediaController
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

	device.receiverController = controllers.NewReceiverController(client, defaultChromecastSenderID, defaultChromecastReceiverID)

	device.mediaController = controllers.NewMediaController(client, defaultChromecastSenderID, "")
	setupMediaCh := make(chan bool)
	go func() {
		status := <-device.receiverController.Incoming
		session := status.GetSessionByNamespace(controllers.MediaControllerNamespace)
		var transportId *string
		if session != nil {
			transportId = session.TransportId
			connectionController := controllers.NewConnectionController(client, defaultChromecastSenderID, *transportId)
			connectionController.Connect()
			device.mediaController.SetDestinationID(*transportId)
			setupMediaCh <- true
		}
		setupMediaCh <- false

	}()
	_, err = device.receiverController.GetStatus(defaultTimeout)
	setup := <-setupMediaCh
	if !setup {
		spew.Dump("Failed to setup the media controller")
	}
	return device, nil
}

//Play just plays.
func (device *Device) Play() {
	device.mediaController.Play(defaultTimeout)
}

//PlayMedia plays a video via the media controller.
func (device *Device) PlayMedia(URL string, MIMEType string) {
	appID := mediaReceiverAppID
	device.receiverController.LaunchApplication(&appID, defaultTimeout, false)
	device.mediaController.Load(URL, MIMEType, defaultTimeout)
	//device.mediaController.Play(defaultTimeout)
}

//PlayMedia plays a video via the media controller.
func (device *Device) TestYoutube(URL string) {
	appID := youtubeAppID
	device.receiverController.LaunchApplication(&appID, defaultTimeout, false)
}

func (device *Device) GetMediaStatus(timeout time.Duration) {
	response, err := device.mediaController.GetStatus(time.Second * 5)
	spew.Dump("Status Media response", response, err)

}

func (device *Device) GetStatus(timeout time.Duration) {
	response, err := device.receiverController.GetStatus(time.Second * 5)
	spew.Dump("Status response", response, err)

}
