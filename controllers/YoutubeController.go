package controllers

import (
	"net/http"

	"github.com/AndreasAbdi/go-castv2/generic"

	"github.com/AndreasAbdi/go-castv2/api"
	"github.com/AndreasAbdi/go-castv2/primitives"
)

/*
Helps with playing the youtube chromecast app.
See
https://github.com/balloob/pychromecast/blob/master/pychromecast/controllers/youtube.py
and
https://github.com/ur1katz/casttube/blob/master/casttube/YouTubeSession.py.
Essentially, you start a session with the website, and the website/session handles  like any other receiver app.
*/

const defaultTimeout = 10

var messageTypeGetSessionID = "getMdxSessionStatus"
var messageTypeStatus = "mdxSessionStatus"

const actionSetPlaylist = "setPlaylist"
const actionRemoveVideo = "removeVideo"
const actionInsertVideo = "insertVideo"
const actionAdd = "addVideo"

//TODO. This will handle the internal operations of executing commands against the chromecast.
type youtubeSession struct {
}

//YoutubeController is the controller for the commands unique to the dashcast.
type YoutubeController struct {
	channel        *primitives.Channel
	screenID       int
	youtubeSession youtubeSession
	loungeID       int
	counter        generic.Counter
}

//NewYoutubeController is a constructor for a dash cast controller.
func NewYoutubeController(client primitives.Client, sourceID, destinationID string) *DashCastController {
	return &DashCastController{
		channel: client.NewChannel(sourceID, destinationID, dashcastControllerNamespace),
	}
}

type youtubeCommand struct {
	primitives.PayloadHeaders
}

//TODO
func (c *YoutubeController) PlayVideo(videoID int) {
	id, err := c.getLoungeID(c.screenID)
	if err != nil {
		c.loungeID = id
	}
	c.bind()
}

//TODO
func (c *YoutubeController) bind() {
	request := c.createBindRequest()
	_, err := c.sendBindRequest(request)
	if err != nil {
		panic(nil)
	}
}

//TODO
func (c *YoutubeController) createBindRequest() http.Request {
	return http.Request{}
}

//TODO
func (c *YoutubeController) sendBindRequest(bindRequest http.Request) (http.Response, error) {
	return http.Response{}, nil
}

//PlayNext adds a video to be played next in the current playlist TODO
func (c *YoutubeController) PlayNext(videoID int) {
	c.sendAction(actionInsertVideo, videoID)
}

//AddToQueue adds the video to the end of the current playlist TODO
func (c *YoutubeController) AddToQueue(videoID int) {
	c.sendAction(actionAdd, videoID)

}

//RemoveFromQueue removes a video from the videoplaylist TODO
func (c *YoutubeController) RemoveFromQueue(videoID int) {
	c.sendAction(actionRemoveVideo, videoID)
}

//TODO: send a request for an action.
func (c *YoutubeController) sendAction(actionType string, videoID int) {

}

type InitializeQueueRequest struct {
}

func (c *YoutubeController) initializeQueue(videoID int) {
	request := createInitializeQueueRequest(videoID)
	sendInitializeQueueRequest(request)
}

//TODO
func createInitializeQueueRequest(videoID int) http.Request {
	return http.Request{}
}

//TODO
func sendInitializeQueueRequest(request http.Request) {

}

type GetScreenIDRequest struct {
	primitives.PayloadHeaders
	screenID int `json:"screen_ids"`
}

func (c *YoutubeController) getScreenID() (int, error) {
	message, err := c.channel.Request(&primitives.PayloadHeaders{Type: messageTypeGetSessionID}, defaultTimeout)
	if err != nil {
		return 0, err
	}
	id, err := toScreenID(message)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//TODO
func (c *YoutubeController) getLoungeID(screenID int) (int, error) {
	request := createGetLoungeIDRequest(screenID)
	response := sendMessage(request)
	defer response.Body.Close()
	id, err := toLoungeID(response)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//TODO
func createGetLoungeIDRequest(screenID int) http.Request {
	return http.Request{}
}

//TODO
func sendMessage(request http.Request) http.Response {
	return http.Response{}
}

//TODO
func toLoungeID(http.Response) (int, error) {
	return 0, nil
}

//TODO
func toScreenID(message *api.CastMessage) (int, error) {
	return 0, nil
}
