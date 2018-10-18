package cast

import (
	"net"
	"time"

	"github.com/AndreasAbdi/gochromecast/configs"
	"github.com/AndreasAbdi/gochromecast/controllers"
	"github.com/AndreasAbdi/gochromecast/controllers/media"
	"github.com/AndreasAbdi/gochromecast/controllers/receiver"
	"github.com/AndreasAbdi/gochromecast/primitives"
)

const defaultTimeout = time.Second * 10

//Device Object to run basic chromecast commands
type Device struct {
	client               *primitives.Client
	heartbeatController  *controllers.HeartbeatController
	connectionController *controllers.ConnectionController
	ReceiverController   *controllers.ReceiverController
	MediaController      *controllers.MediaController
	YoutubeController    *controllers.YoutubeController
}

//NewDevice is constructor for Device struct
func NewDevice(host net.IP, port int) (Device, error) {
	var device Device

	client, err := primitives.NewClient(host, port)
	if err != nil {
		return device, err
	}
	device.client = client

	device.heartbeatController = controllers.NewHeartbeatController(client, defaultChromecastSenderID, defaultChromecastReceiverID)
	device.heartbeatController.Start()

	device.connectionController = controllers.NewConnectionController(client, defaultChromecastSenderID, defaultChromecastReceiverID)
	device.connectionController.Connect()

	device.ReceiverController = controllers.NewReceiverController(client, defaultChromecastSenderID, defaultChromecastReceiverID)

	device.MediaController = controllers.NewMediaController(client, defaultChromecastSenderID, device.ReceiverController)

	device.YoutubeController = controllers.NewYoutubeController(client, defaultChromecastSenderID, device.ReceiverController)
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

//QuitApplication that is currently running on the device
func (device *Device) QuitApplication(timeout time.Duration) {
	status, err := device.ReceiverController.GetStatus(timeout)
	if err != nil {
		return
	}
	for _, appSessions := range status.Applications {
		session := appSessions.SessionID
		device.ReceiverController.StopApplication(session, timeout)
	}
}

//PlayYoutubeVideo launches the youtube app and tries to play the video based on its id.
func (device *Device) PlayYoutubeVideo(videoID string) {
	appID := configs.YoutubeAppID
	device.ReceiverController.LaunchApplication(&appID, defaultTimeout, false)
	device.YoutubeController.PlayVideo(videoID, "")
}

//GetMediaStatus of current media controller
func (device *Device) GetMediaStatus(timeout time.Duration) []*media.MediaStatus {
	response, err := device.MediaController.GetStatus(time.Second * 5)
	if err != nil {
		emptyStatus := make([]*media.MediaStatus, 0)
		return emptyStatus
	}
	return response
}

//GetStatus of the device.
func (device *Device) GetStatus(timeout time.Duration) *receiver.Status {
	response, err := device.ReceiverController.GetStatus(time.Second * 5)
	if err != nil {
		return nil
	}
	return response
}
