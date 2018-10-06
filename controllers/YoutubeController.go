package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/AndreasAbdi/go-castv2/controllers/youtube"

	"github.com/AndreasAbdi/go-castv2/api"
	"github.com/AndreasAbdi/go-castv2/primitives"
)

/*
Helps with playing the youtube chromecast app.
See
https://github.com/balloob/pychromecast/blob/master/pychromecast/controllers/youtube.py
and
https://github.com/ur1katz/casttube/blob/master/casttube/YouTubeSession.py.
https://github.com/CBiX/gotubecast/blob/master/main.go
https://github.com/mutantmonkey/youtube-remote/blob/master/remote.py
Essentially, you start a session with the website, and the website/session handles  like any other receiver app.
*/

const loungeIDHeader = "X-YouTube-LoungeId-Token"

var messageTypeGetSessionID = "getMdxSessionStatus"
var responseTypeSessionStatus = "mdxSessionStatus"

//YoutubeController is the controller for the commands unique to the dashcast.
type YoutubeController struct {
	connection *mediaConnection
	screenID   *string
	incoming   chan *string
	session    *youtube.Session
}

//NewYoutubeController is a constructor for a dash cast controller.
func NewYoutubeController(client *primitives.Client, sourceID string, receiver *ReceiverController) *YoutubeController {
	connection := NewMediaConnection(client, receiver, youtubeControllerNamespace, sourceID)
	controller := YoutubeController{
		connection: connection,
		incoming:   make(chan *string, 0),
	}
	connection.OnMessage(responseTypeSessionStatus, controller.onStatus)
	return &controller
}

type youtubeCommand struct {
	primitives.PayloadHeaders
}

//PlayVideo initializes a new queue and plays the video
func (c *YoutubeController) PlayVideo(videoID string, listID string) error {
	err := c.session.InitializeQueue(videoID, listID)
	if err == nil {
		return nil
	}

	isActive := c.ensureSessionActive()
	if isActive {
		c.session.InitializeQueue(videoID, listID)
		return nil
	}
	return InitializationError{}
}

//ClearPlaylist of current videos in the chromecast playlist
func (c *YoutubeController) ClearPlaylist() {
	c.runFast(func() error {
		c.session.ClearQueue()
		return nil
	})
}

//PlayNext adds a video to be played next in the current playlist
func (c *YoutubeController) PlayNext(videoID string) {
	c.runFast(func() error {
		c.session.PlayNext(videoID)
		return nil
	})
}

//AddToQueue adds the video to the end of the current playlist
func (c *YoutubeController) AddToQueue(videoID string) {
	c.runFast(func() error {
		c.session.AddToQueue(videoID)
		return nil
	})
}

//RemoveFromQueue removes a video from the video playlist
func (c *YoutubeController) RemoveFromQueue(videoID string) {
	c.runFast(func() error {
		c.session.RemoveFromQueue(videoID)
		return nil
	})
}

func (c *YoutubeController) runFast(command func() error) {
	err := command()
	if err == nil {
		return
	}
	isActive := c.ensureSessionActive()
	if isActive {
		command()
	}
}

func (c *YoutubeController) ensureSessionActive() bool {
	newScreenID, err := c.updateScreenID()
	if err != nil {
		log.Print("Failed to get screenID")
		return false
	}
	if c.screenID == newScreenID {
		return true
	}
	c.updateYoutubeSession(newScreenID)
	c.screenID = newScreenID
	return true
}

func (c *YoutubeController) updateScreenID() (*string, error) {
	screenID, err := c.getScreenID(time.Second * 5)
	if err != nil {
		return nil, err
	}
	return screenID, nil

}

func (c *YoutubeController) updateYoutubeSession(newScreenID *string) error {
	c.session = youtube.NewSession(*newScreenID)
	return c.session.StartSession()
}

func (c *YoutubeController) onStatus(message *api.CastMessage) {
	response := &youtube.ScreenStatus{}
	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)
	if err != nil {
		return
	}
	select {
	case c.incoming <- &response.Data.ScreenID:
	case <-time.After(time.Second):
	}
}

func (c *YoutubeController) getScreenID(timeout time.Duration) (*string, error) {

	waitCh := make(chan bool)
	var screenID *string
	go func() {
		screenID = <-c.incoming
		waitCh <- true
	}()

	c.connection.Request(
		&primitives.PayloadHeaders{Type: messageTypeGetSessionID},
		0)
	select {
	case <-waitCh:
		return screenID, nil
	case <-time.After(timeout):
		return nil, errors.New("Failed to get screen ID, timed out")
	}
}
