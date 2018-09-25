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
	connection     *mediaConnection
	Incoming       chan []*media.MediaStatus
	MediaSessionID int
}

var getMediaStatus = primitives.PayloadHeaders{Type: "GET_STATUS"}

var commandMediaPlay = primitives.PayloadHeaders{Type: "PLAY"}
var commandMediaPause = primitives.PayloadHeaders{Type: "PAUSE"}
var commandMediaStop = primitives.PayloadHeaders{Type: "STOP"}

const responseTypeMediaStatus = "MEDIA_STATUS"

//NewMediaController is the constructors for the media controller
func NewMediaController(client *primitives.Client, sourceID string, receiverController *ReceiverController) *MediaController {
	mediaConnection := NewMediaConnection(client, receiverController, MediaControllerNamespace, sourceID)
	controller := &MediaController{
		Incoming:   make(chan []*media.MediaStatus, 0),
		connection: mediaConnection,
	}

	controller.connection.OnMessage(responseTypeMediaStatus, func(message *api.CastMessage) {
		controller.onStatus(message)
	})
	return controller
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
		log.Printf("Incoming media status, but we aren't listening. %v", response)
	}

	return response.Status, nil
}

//GetStatus attempts to request the chromecast return the status of the current media controller channel
func (c *MediaController) GetStatus(timeout time.Duration) ([]*media.MediaStatus, error) {

	spew.Dump("getting media Status")

	message, err := c.connection.Request(&getMediaStatus, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to get media status: %s", err)
	}

	spew.Dump("got media Status", message)

	return c.onStatus(message)
}

//TODO
//Load sends a load request to play a generic media event
func (c *MediaController) Load(url string, contentTypeString string, timeout time.Duration) (*api.CastMessage, error) {
	//TODO should do something about messaging with the contenttype, so it works with different media types. so we can attach more metadata
	//TODO also should be sending a message of type media data( should probably actually construct the request)
	//c.GetStatus(defaultTimeout)
	mediaData, err := c.constructMediaData(url, contentTypeString)
	if err != nil {
		return nil, err
	}
	loadCommand := c.constructLoadCommand(mediaData)

	_, err = c.connection.Request(&loadCommand, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send play command: %s", err)
	}
	fmt.Printf("There is no error")
	return nil, nil
}

func (c *MediaController) constructLoadCommand(mediaData *media.MediaData) media.LoadCommand {
	return media.LoadCommand{
		PayloadHeaders: primitives.PayloadHeaders{Type: eventTypeLoad},
		Media:          *mediaData,
		Autoplay:       true,
		CurrentTime:    0,
		CustomData:     nil,
	}
}

func (c *MediaController) constructMediaData(url string, contentTypeString string) (*media.MediaData, error) {
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
	return &mediaData, nil
}

//Play sends the play command so that the chromecast session is resumed
func (c *MediaController) Play(timeout time.Duration) (*api.CastMessage, error) {

	message, err := c.sendMessage(commandMediaPlay, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send play command: %s", err)
	}
	fmt.Printf("There is no error")

	return message, nil
}

//Pause sends the pause command to the chromecast
func (c *MediaController) Pause(timeout time.Duration) (*api.CastMessage, error) {

	message, err := c.sendMessage(commandMediaPause, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send pause command: %s", err)
	}
	spew.Dump("Pause Command:", message.PayloadUtf8)
	return message, nil
}

//Stop sends the stop command to the chromecast
func (c *MediaController) Stop(timeout time.Duration) (*api.CastMessage, error) {

	message, err := c.sendMessage(commandMediaStop, timeout)
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

func (c *MediaController) sendMessage(payload primitives.PayloadHeaders, timeout time.Duration) (*api.CastMessage, error) {
	c.updateForNewSession(timeout)
	return c.connection.Request(&media.MediaCommand{
		PayloadHeaders: payload,
		MediaSessionID: c.MediaSessionID}, timeout)
}

//UpdateForNewSession refreshes the media controller for a new media session that's been executed.
func (c *MediaController) updateForNewSession(timeout time.Duration) {
	waitStatusCh := make(chan bool)
	go func() {
		status := <-c.Incoming
		if len(status) <= 0 {
			waitStatusCh <- false
			return
		}
		c.MediaSessionID = status[0].MediaSessionID
		waitStatusCh <- true
	}()

	c.GetStatus(timeout)
	<-waitStatusCh
}
