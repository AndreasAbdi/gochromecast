package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AndreasAbdi/go-castv2"
	"github.com/AndreasAbdi/go-castv2/api"
	"github.com/davecgh/go-spew/spew"
)

//ReceiverController is a chromecast controller for the receiver namespace. This involves
type ReceiverController struct {
	interval           time.Duration
	currentApplication string
	channel            *castv2.Channel
	Incoming           chan *ReceiverStatus
}

//NewReceiverController is for building a new receiver controller
func NewReceiverController(client *castv2.Client, sourceID, destinationID string) *ReceiverController {
	controller := &ReceiverController{
		channel:  client.NewChannel(sourceID, destinationID, receiverControllerNamespace),
		Incoming: make(chan *ReceiverStatus, 0),
	}

	controller.channel.OnMessage(receiverControllerSystemEventReceiverStatus, controller.onStatus)

	return controller
}

func (c *ReceiverController) onStatus(message *api.CastMessage) {
	spew.Dump("Got status message", message)

	response := &StatusResponse{}

	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)

	if err != nil {
		log.Printf("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
		return
	}

	select {
	case c.Incoming <- response.Status:
		log.Printf("Delivered status")
		/*
			pychromecast presumes the first application in the applications list is the one running.
			Honestly don't know how to run multiple applications so would be hard for me to test.
		*/
		appData := response.Status.Applications[0]
		c.currentApplication = *appData.AppID
	case <-time.After(time.Second):
		log.Printf("Incoming status, but we aren't listening. %v", response.Status)
	}

}

//GetSessionByNamespace attempts to return the first session with a specified namespace.
func (s *ReceiverStatus) GetSessionByNamespace(namespace string) *ApplicationSession {

	for _, app := range s.Applications {
		for _, ns := range app.Namespaces {
			if ns.Name == namespace {
				return app
			}
		}
	}

	return nil
}

//GetStatus attempts to receive the current status of the controllers chromecast device.
func (c *ReceiverController) GetStatus(timeout time.Duration) (*ReceiverStatus, error) {
	message, err := c.channel.Request(&castv2.PayloadHeaders{Type: receiverControllerSystemEventGetStatus}, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to get receiver status: %s", err)
	}
	c.onStatus(message)

	response := &StatusResponse{}

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
	log.Printf("Attempting to launch an application: %s\n", *appID)
	c.channel.Request(&LaunchRequest{
		PayloadHeaders: castv2.PayloadHeaders{Type: receiverControllerSystemEventLaunch},
		AppID:          appID,
	}, timeout)
}

//TODO: so application termination requires sessionID, need to figure out how to rewrite code to work with that.
//Actually, you know what? we could do it so that there is a wrapper that sends requests to these thingies.
func (c *ReceiverController) StopApplication(sessionID *string, timeout time.Duration) {
	log.Println("Attempting to stop the current application")
	c.channel.Request(&StopRequest{
		PayloadHeaders: castv2.PayloadHeaders{Type: receiverControllerSystemEventStop},
		sessionID:      sessionID,
	}, timeout)
}

//SetVolume sets the volume on the controller's chromecast.
func (c *ReceiverController) SetVolume(volume *Volume, timeout time.Duration) (*api.CastMessage, error) {
	return c.channel.Request(&ReceiverStatus{
		PayloadHeaders: castv2.PayloadHeaders{Type: receiverControllerSystemEventSetVolume},
		Volume:         volume,
	}, timeout)
}

//GetVolume gets the volume on the controller's chromecast.
func (c *ReceiverController) GetVolume(timeout time.Duration) (*Volume, error) {
	status, err := c.GetStatus(timeout)

	if err != nil {
		return nil, err
	}

	return status.Volume, err
}
