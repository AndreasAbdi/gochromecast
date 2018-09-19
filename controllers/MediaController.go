package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AndreasAbdi/go-castv2/api"
	"github.com/AndreasAbdi/go-castv2/controllers/media"
	"github.com/AndreasAbdi/go-castv2/primitives"
	"github.com/davecgh/go-spew/spew"
)

/*
MediaController is a type of chromecast controller.
Use it to load, play, pause, stop media.
Use it to enable or disable subtitles.
*/
type MediaController struct {
	interval       time.Duration
	channel        *primitives.Channel
	Incoming       chan []*media.MediaStatus
	DestinationID  string
	MediaSessionID int
}

var getMediaStatus = primitives.PayloadHeaders{Type: "GET_STATUS"}

var commandMediaPlay = primitives.PayloadHeaders{Type: "PLAY"}
var commandMediaPause = primitives.PayloadHeaders{Type: "PAUSE"}
var commandMediaStop = primitives.PayloadHeaders{Type: "STOP"}

const responseTypeMediaStatus = "MEDIA_STATUS"

//NewMediaController is the constructors for the media controller
func NewMediaController(client *primitives.Client, sourceID, destinationID string) *MediaController {
	controller := &MediaController{
		channel:       client.NewChannel(sourceID, destinationID, mediaControllerNamespace),
		Incoming:      make(chan []*media.MediaStatus, 0),
		DestinationID: destinationID,
	}

	controller.channel.OnMessage(responseTypeMediaStatus, func(message *api.CastMessage) {
		controller.onStatus(message)
	})

	return controller
}

//SetDestinationID sets the target destination for the media controller
func (c *MediaController) SetDestinationID(id string) {
	c.channel.DestinationID = id
	c.DestinationID = id
}

func (c *MediaController) onStatus(message *api.CastMessage) ([]*media.MediaStatus, error) {
	spew.Dump("Got media status message", message)

	response := &media.MediaStatusResponse{}

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
func (c *MediaController) GetStatus(timeout time.Duration) ([]*media.MediaStatus, error) {

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
func (c *MediaController) Load(url string, contentTypeString string, timeout time.Duration) (*api.CastMessage, error) {
	//TODO should do something about messaging with the contenttype, so it works with different media types. so we can attach more metadata
	//TODO also should be sending a message of type media data( should probably actually construct the request)
	contentType, err := media.NewContentType(contentTypeString)
	if err != nil {
		return nil, err
	}
	contentID, err := media.NewContentID(url)
	if err != nil {
		return nil, err
	}
	builder, err := media.NewGenericMediaDataBuilder(contentID, contentType, media.NoneStreamType)
	if err != nil {
		return nil, err
	}
	mediaData, err := builder.Build()
	if err != nil {
		return nil, err
	}

	_, err = c.channel.Request(&media.LoadCommand{
		PayloadHeaders: primitives.PayloadHeaders{Type: eventTypeLoad},
		Media:          mediaData,
		Autoplay:       true,
		CurrentTime:    0,
		CustomData:     nil,
	}, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send play command: %s", err)
	}
	fmt.Printf("There is no error")
	return nil, nil
}

//Play sends the play command so that the chromecast session is resumed
func (c *MediaController) Play(timeout time.Duration) (*api.CastMessage, error) {

	message, err := c.channel.Request(&media.MediaCommand{
		PayloadHeaders: commandMediaPlay,
		MediaSessionID: c.MediaSessionID}, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send play command: %s", err)
	}

	return message, nil
}

//Pause sends the pause command to the chromecast
func (c *MediaController) Pause(timeout time.Duration) (*api.CastMessage, error) {

	message, err := c.channel.Request(&media.MediaCommand{
		PayloadHeaders: commandMediaPause,
		MediaSessionID: c.MediaSessionID}, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send pause command: %s", err)
	}

	return message, nil
}

//Stop sends the stop command to the chromecast
func (c *MediaController) Stop(timeout time.Duration) (*api.CastMessage, error) {

	message, err := c.channel.Request(&media.MediaCommand{
		PayloadHeaders: commandMediaStop,
		MediaSessionID: c.MediaSessionID}, timeout)
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
