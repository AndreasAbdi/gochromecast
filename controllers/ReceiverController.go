package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AndreasAbdi/gochromecast/api"
	"github.com/AndreasAbdi/gochromecast/controllers/receiver"
	"github.com/AndreasAbdi/gochromecast/primitives"
)

//ReceiverController is a chromecast controller for the receiver namespace. This involves
type ReceiverController struct {
	interval time.Duration
	channel  *primitives.Channel
	Incoming chan *receiver.Status
}

//NewReceiverController is for building a new receiver controller
func NewReceiverController(client *primitives.Client, sourceID, destinationID string) *ReceiverController {
	controller := &ReceiverController{
		channel:  client.NewChannel(sourceID, destinationID, receiverControllerNamespace),
		Incoming: make(chan *receiver.Status, 0),
	}

	controller.channel.OnMessage(receiverControllerSystemEventReceiverStatus, controller.onStatus)

	return controller
}

func (c *ReceiverController) onStatus(message *api.CastMessage) {
	response := &receiver.StatusResponse{}

	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)

	if err != nil {
		return
	}

	select {
	case c.Incoming <- response.Status:
	case <-time.After(time.Second):
	}

}

//GetStatus attempts to receive the current status of the controllers chromecast device.
func (c *ReceiverController) GetStatus(timeout time.Duration) (*receiver.Status, error) {
	message, err := c.channel.Request(&primitives.PayloadHeaders{Type: receiverControllerSystemEventGetStatus}, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to get receiver status: %s", err)
	}
	c.onStatus(message)

	response := &receiver.StatusResponse{}

	err = json.Unmarshal([]byte(*message.PayloadUtf8), response)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
	}

	return response.Status, nil
}

/*
LaunchApplication attempts to launch an application on the chromecast.
forceLaunch forces the app to run even if something is already running.
*/
func (c *ReceiverController) LaunchApplication(appID *string, timeout time.Duration, forceLaunch bool) {
	//TODO: test out force launch and actually write it.
	c.channel.Request(&receiver.LaunchRequest{
		PayloadHeaders: primitives.PayloadHeaders{Type: receiverControllerSystemEventLaunch},
		AppID:          appID,
	}, timeout)
}

//TODO: so application termination requires sessionID, need to figure out how to rewrite code to work with that.
//Actually, you know what? we could do it so that there is a wrapper that sends requests to these thingies.
func (c *ReceiverController) StopApplication(sessionID *string, timeout time.Duration) {
	c.channel.Request(&receiver.StopRequest{
		PayloadHeaders: primitives.PayloadHeaders{Type: receiverControllerSystemEventStop},
		SessionID:      sessionID,
	}, timeout)
}

//SetVolume sets the volume on the controller's chromecast.
func (c *ReceiverController) SetVolume(volume *receiver.Volume, timeout time.Duration) (*api.CastMessage, error) {
	return c.channel.Request(&receiver.Status{
		PayloadHeaders: primitives.PayloadHeaders{Type: receiverControllerSystemEventSetVolume},
		Volume:         volume,
	}, timeout)
}

//GetVolume gets the volume on the controller's chromecast.
func (c *ReceiverController) GetVolume(timeout time.Duration) (*receiver.Volume, error) {
	status, err := c.GetStatus(timeout)

	if err != nil {
		return nil, err
	}

	return status.Volume, err
}
