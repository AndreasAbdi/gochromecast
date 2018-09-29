package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/imdario/mergo"

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
	currentStatus  *media.MediaStatus
	Incoming       chan []*media.MediaStatus
	MediaSessionID int
}

//TODO: probably should make this not vars. Or pull them into smaller scope.
var getMediaStatus = primitives.PayloadHeaders{Type: "GET_STATUS"}
var commandMediaPlay = primitives.PayloadHeaders{Type: "PLAY"}
var commandMediaPause = primitives.PayloadHeaders{Type: "PAUSE"}
var commandMediaStop = primitives.PayloadHeaders{Type: "STOP"}
var commandMediaNext = primitives.PayloadHeaders{Type: "NEXT"}
var commandMediaPrevious = primitives.PayloadHeaders{Type: "PREVIOUS"}
var commandMediaSeek = primitives.PayloadHeaders{Type: "SEEK"}
var commandSetSubtitles = primitives.PayloadHeaders{Type: "EDIT_TRACKS_INFO"}

const responseTypeMediaStatus = "MEDIA_STATUS"
const skipTimeBuffer = -5

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
	spew.Dump("[MEDIA] Got media status message", message)

	response := &media.MediaStatusResponse{}

	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
	}
	c.MediaSessionID = response.Status[0].MediaSessionID
	mergo.Merge(c.currentStatus, response.Status[0], mergo.WithOverride)
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

//Play sends the play command so that the chromecast session is resumed
func (c *MediaController) Play(timeout time.Duration) (*api.CastMessage, error) {
	return c.sendCommand(commandMediaPlay, timeout)
}

//Pause sends the pause command to the chromecast
func (c *MediaController) Pause(timeout time.Duration) (*api.CastMessage, error) {
	return c.sendCommand(commandMediaPause, timeout)
}

//Stop sends the stop command to the chromecast
func (c *MediaController) Stop(timeout time.Duration) (*api.CastMessage, error) {
	return c.sendCommand(commandMediaStop, timeout)
}

//Next goes to the next video
func (c *MediaController) Next(timeout time.Duration) (*api.CastMessage, error) {
	return c.sendCommand(commandMediaNext, timeout)
}

//Previous goes to the previous video
func (c *MediaController) Previous(timeout time.Duration) (*api.CastMessage, error) {
	return c.sendCommand(commandMediaPrevious, timeout)
}

//Rewind to the beginning.
func (c *MediaController) Rewind(timeout time.Duration) (*api.CastMessage, error) {
	//
	return c.Seek(0, timeout)
}

//Skip to the end
func (c *MediaController) Skip(timeout time.Duration) (*api.CastMessage, error) {
	if c.currentStatus == nil || c.currentStatus.Media.Duration == nil {
		return nil, errors.New("No media playing, can't skip")
	}
	return c.Seek(*c.currentStatus.Media.Duration-5, timeout)
}

//Seek to some time in the video
func (c *MediaController) Seek(seconds float64, timeout time.Duration) (*api.CastMessage, error) {
	seekCommand := media.CreateSeekCommand(seconds)
	_, err := c.connection.Request(&seekCommand, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send play command: %s", err)
	}
	fmt.Printf("There is no error")
	return nil, nil
}

func (c *MediaController) sendCommand(command primitives.PayloadHeaders, timeout time.Duration) (*api.CastMessage, error) {
	message, err := c.sendMessage(command, timeout)
	if err != nil {
		return nil, fmt.Errorf("Failed to send %v command: %s", command.Type, err)
	}
	spew.Dump("%v Command: \n %s", command.Type, message.PayloadUtf8)
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
		c.currentStatus = status[0]
		waitStatusCh <- true
	}()

	c.GetStatus(timeout)
	<-waitStatusCh
}
