package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/AndreasAbdi/go-castv2/controllers/youtube"
	"github.com/AndreasAbdi/go-castv2/generic"
	"github.com/davecgh/go-spew/spew"

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

const youtubeBaseURL = "https://www.youtube.com/"
const bindURL = youtubeBaseURL + "api/lounge/bc/bind"
const loungeTokenURL = youtubeBaseURL + "api/lounge/pairing/get_lounge_token_batch"

const loungeIDHeader = "X-YouTube-LoungeId-Token"

var messageTypeGetSessionID = "getMdxSessionStatus"
var responseTypeSessionStatus = "mdxSessionStatus"

const actionSetPlaylist = "setPlaylist"
const actionRemoveVideo = "removeVideo"
const actionInsertVideo = "insertVideo"
const actionAdd = "addVideo"

const bindVersion = "8"
const bindCVersion = "1"

const gSessionIDKey = "gsessionid"
const cVersionKey = "CVER"
const requestIDKey = "RID"
const sessionIDKey = "SID"
const versionKey = "VER"
const actionKey = "__sc"
const countKey = "count"
const videoIDKey = "_videoId"

const defaultCount = 1

//YoutubeController is the controller for the commands unique to the dashcast.
type YoutubeController struct {
	connection     *mediaConnection
	screenID       string
	sessionID      string
	gSessionID     string
	loungeID       string
	incoming       chan *string
	requestCounter generic.Counter
	sessionCounter generic.Counter
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

func (c *YoutubeController) Test() {
	screenID, err := c.getScreenID(time.Second * 5)
	if err != nil {
		spew.Dump("Failed to get screen ID")
		return
	}
	spew.Dump("Got screen ID", screenID)

	loungeToken, err := youtube.GetLoungeToken(screenID)
	c.loungeID = loungeToken
	if err != nil {
		spew.Dump("Failed to get lounge token")
		return
	}

	c.bind(screenID, loungeToken)
	c.initializeQueue("cwQgjq0mCdE", "")
}

func (c *YoutubeController) bind(screenID *string, loungeToken string) error {
	c.requestCounter.Reset()
	c.sessionCounter.Reset()
	requestID := c.requestCounter.GetAndIncrement()
	request := youtube.CreateBindRequest(requestID, loungeToken)
	response, err := request.Post()
	sessionID, gSessionID, err := youtube.ParseResponse(response)
	if err != nil {
		return err
	}
	c.sessionID = sessionID
	c.gSessionID = gSessionID
	//assign self data to sid and gsiddata.
	return nil
}

//TODO
func (c *YoutubeController) PlayVideo(videoID int) {

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
	request, err := c.createActionRequest(actionType, videoID)
	if err != nil {
		//TODO
		return
	}
	c.sendRequest(&request)
}

func (c *YoutubeController) createActionRequest(actionType string, videoID int) (http.Request, error) {
	request := http.Request{}
	requestID := c.requestCounter.GetAndIncrement()
	message := map[string]interface{}{
		actionKey:  actionType,
		videoIDKey: videoID,
		countKey:   defaultCount,
	}

	messageInBytes, err := json.Marshal(message)
	if err != nil {
		return request, err
	}
	req, err := http.NewRequest("POST", bindURL, bytes.NewBuffer(messageInBytes))
	if err != nil {
		return request, err
	}
	req.Header.Set(loungeIDHeader, c.loungeID)
	req.URL.Query().Add(sessionIDKey, c.sessionID)
	req.URL.Query().Add(requestIDKey, strconv.Itoa(requestID))
	req.URL.Query().Add(gSessionIDKey, c.gSessionID)
	req.URL.Query().Add(cVersionKey, bindCVersion)
	req.URL.Query().Add(versionKey, bindVersion)
	return request, nil

}

func (c *YoutubeController) sendRequest(request *http.Request) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

}

func (c *YoutubeController) initializeQueue(videoID string, listID string) {
	requestParams := c.CreateInitializeQueueRequestParams(videoID, listID)
	request := youtube.CreateInitializeQueueRequest(requestParams)
	spew.Dump("Request info", request)
	response, err := request.Post()
	if err != nil {
		spew.Dump("Failed to initialize queue")
	}
	spew.Dump("Got response", response)
}

func (c *YoutubeController) CreateInitializeQueueRequestParams(videoID string, listID string) youtube.InitializeQueueRequestParams {
	return youtube.InitializeQueueRequestParams{
		VideoID:             videoID,
		ListID:              listID,
		LoungeID:            c.loungeID,
		RequestCount:        c.requestCounter.GetAndIncrement(),
		SessionRequestCount: c.sessionCounter.GetAndIncrement(),
		SessionID:           c.sessionID,
		GSessionID:          c.gSessionID,
	}
}

func (c *YoutubeController) onStatus(message *api.CastMessage) {
	spew.Dump("Got youtube status message")
	response := &youtube.ScreenStatus{}
	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)
	if err != nil {
		spew.Dump("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
		return
	}
	select {
	case c.incoming <- &response.Data.ScreenID:
		spew.Dump("Delivered status. %v", response)
	case <-time.After(time.Second):
		spew.Dump("Incoming youtube status, but we aren't listening. %v", response)
	}
}

func (c *YoutubeController) getScreenID(timeout time.Duration) (*string, error) {

	waitCh := make(chan bool)
	var screenID *string
	go func() {
		spew.Dump("Listening for incoming youtube status")
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

	return nil, errors.New("Shouldn't ever get here.")
}
