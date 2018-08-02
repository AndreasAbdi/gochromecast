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

/*
MediaController is a type of chromecast controller.
Use it to load, play, pause, stop media.
Use it to enable or disable subtitles.
*/
type MediaController struct {
	interval       time.Duration
	channel        *castv2.Channel
	Incoming       chan []*MediaStatus
	DestinationID  string
	MediaSessionID int
}

var getMediaStatus = castv2.PayloadHeaders{Type: "GET_STATUS"}

var commandMediaPlay = castv2.PayloadHeaders{Type: "PLAY"}
var commandMediaPause = castv2.PayloadHeaders{Type: "PAUSE"}
var commandMediaStop = castv2.PayloadHeaders{Type: "STOP"}
var commandMediaLoad = castv2.PayloadHeaders{Type: "LOAD"}

//NewMediaController is the constructors for the media controller
func NewMediaController(client *castv2.Client, sourceID, destinationID string) *MediaController {
	controller := &MediaController{
		channel:       client.NewChannel(sourceID, destinationID, mediaControllerNamespace),
		Incoming:      make(chan []*MediaStatus, 0),
		DestinationID: destinationID,
	}

	controller.channel.OnMessage("MEDIA_STATUS", func(message *api.CastMessage) {
		controller.onStatus(message)
	})

	return controller
}

//SetDestinationID sets the target destination for the media controller
func (c *MediaController) SetDestinationID(id string) {
	c.channel.DestinationID = id
	c.DestinationID = id
}

func (c *MediaController) onStatus(message *api.CastMessage) ([]*MediaStatus, error) {
	spew.Dump("Got media status message", message)

	response := &MediaStatusResponse{}

	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
	}

	select {
	case c.Incoming <- response.Status:
	default:
		log.Printf("Incoming status, but we aren't listening. %v", response)
	}

	return response.Status, nil
}

//GetStatus attempts to request the chromecast return the status of the current media controller channel
func (c *MediaController) GetStatus(timeout time.Duration) ([]*MediaStatus, error) {

	spew.Dump("getting media Status")

	message, err := c.channel.Request(&getMediaStatus, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to get receiver status: %s", err)
	}

	spew.Dump("got media Status", message)

	return c.onStatus(message)
}

//TODO
//Load sends a load request to play a generic media event
func (c *MediaController) Load(url string, contentType string, timeout time.Duration) (*api.CastMessage, error) {
	//TODO should do something about messaging with the request type so we can attach more metadata
	//TODO also should be sending a message of type media data( should probably actually construct the request)
	builder := GenericMediaDataBuilder{}
	builder.SetContentID(url)
	builder.SetContentType(contentType)
	mediaData, err := builder.Build()
	if err != nil {
		return nil, err
	}
	message, err := c.channel.Request(&LoadCommand{
		commandMediaLoad,
		mediaData,
		true,
		0,
		nil,
	}, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send play command: %s", err)
	}

	return message, nil
}

//Play sends the play command so that the chromecast session is resumed
func (c *MediaController) Play(timeout time.Duration) (*api.CastMessage, error) {

	message, err := c.channel.Request(&MediaCommand{commandMediaPlay, c.MediaSessionID}, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send play command: %s", err)
	}

	return message, nil
}

//Pause sends the pause command to the chromecast
func (c *MediaController) Pause(timeout time.Duration) (*api.CastMessage, error) {

	message, err := c.channel.Request(&MediaCommand{commandMediaPause, c.MediaSessionID}, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send pause command: %s", err)
	}

	return message, nil
}

//Stop sends the stop command to the chromecast
func (c *MediaController) Stop(timeout time.Duration) (*api.CastMessage, error) {

	message, err := c.channel.Request(&MediaCommand{commandMediaStop, c.MediaSessionID}, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send stop command: %s", err)
	}

	return message, nil
}

//TODO
//EnableSubtitles sends the enable subtitles command to the chromecast
func (c *MediaController) EnableSubtitles(timeout time.Duration) (*api.CastMessage, error) {
	return nil, nil
}

//TODO
//DisableSubtitles sends the disable subtitles command to the chromecast
func (c *MediaController) DisableSubtitles(timeout time.Duration) (*api.CastMessage, error) {
	return nil, nil
}
