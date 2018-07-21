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

type ReceiverController struct {
	interval time.Duration
	channel  *castv2.Channel
	Incoming chan *ReceiverStatus
}

func NewReceiverController(client *castv2.Client, sourceId, destinationId string) *ReceiverController {
	controller := &ReceiverController{
		channel:  client.NewChannel(sourceId, destinationId, receiverControllerNamespace),
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
	case <-time.After(time.Second):
		log.Printf("Incoming status, but we aren't listening. %v", response.Status)
	}

}

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

func (c *ReceiverController) SetVolume(volume *Volume, timeout time.Duration) (*api.CastMessage, error) {
	return c.channel.Request(&ReceiverStatus{
		PayloadHeaders: castv2.PayloadHeaders{Type: receiverControllerSystemEventSetVolume},
		Volume:         volume,
	}, timeout)
}

func (c *ReceiverController) GetVolume(timeout time.Duration) (*Volume, error) {
	status, err := c.GetStatus(timeout)

	if err != nil {
		return nil, err
	}

	return status.Volume, err
}
