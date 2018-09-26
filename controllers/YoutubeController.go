package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req"

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

//Spoof data to be passed for binding the screen and the token
const defaultDeviceType = "REMOTE_CONTROL"
const defaultDeviceName = "GOCAST"
const defaultDeviceID = "aaaaaaaaaaaaaaaaaaaaaaaa"
const bindPairingType = "cast"
const defaultAppName = "GOCAST_REMOTE_APP"

var bindData = map[string][]string{
	"device":       []string{defaultDeviceType},
	"id":           []string{defaultDeviceID},
	"name":         []string{defaultDeviceID},
	"mdx-version":  []string{string(3)},
	"pairing_type": []string{bindPairingType},
	"app":          []string{defaultAppName},
}

var defaultHeaders = map[string]string{
	"Origin":       youtubeBaseURL,
	"Content-Type": "application/x-www-form-urlencoded"}

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

const listIDKey = "_listId"
const actionKey = "__sc"
const currentTimeKey = "_currentTime"
const currentIndexKey = "_currentIndex"
const audioOnlyKey = "_audioOnly"
const countKey = "count"
const videoIDKey = "_videoId"

const defaultTime = "0"
const defaultIndex = -1
const defaultAudioOnlySetting = "false"
const defaultCount = 1

const requestPrefixFormat = "req%d"

/*
Requests to the bind API return objects of type.
"[[0,[\"c\",\"19AB39151763497F\",\"\",8]\n]\n,[1,[\"S\",\"d6CNYWDUZb40UcroBuzH6QZJti79F-mc\"]]\n,[2,[\"loungeStatus\",{\"devices\":\"[{\\\"app\\\":\\\"lb-v4\\\",\\\"capabilities\\\":\\\"dsp,que,mus\\\",\\\"clientName\\\":\\\"tvhtml5\\\",\\\"experiments\\\":\\\"\\\",\\\"name\\\":\\\"Chromecast\\\",\\\"id\\\":\\\"1ed072b4-b75a-4878-88d0-fe9e6625d9ec\\\",\\\"type\\\":\\\"LOUNGE_SCREEN\\\",\\\"hasCc\\\":\\\"true\\\"},{\\\"app\\\":\\\"GOCAST_REMOTE_APP\\\",\\\"pairingType\\\":\\\"cast\\\",\\\"capabilities\\\":\\\"que,mus\\\",\\\"clientName\\\":\\\"unknown\\\",\\\"experiments\\\":\\\"\\\",\\\"name\\\":\\\"21b78ce1-4311-4c5e-8ef5-0101eddf5671\\\",\\\"remoteControllerUrl\\\":\\\"\\\",\\\"id\\\":\\\"21b78ce1-4311-4c5e-8ef5-0101eddf5671\\\",\\\"type\\\":\\\"REMOTE_CONTROL\\\",\\\"localChannelEncryptionKey\\\":\\\"wMphRtC_eiqqMvJk61EWvN-k1rA7IA72NzG2KMqPxPU\\\"}]\"}]]\n,[3,[\"playlistModified\",{\"videoIds\":\"\"}]]\n,[4,[\"onAutoplayModeChanged\",{\"autoplayMode\":\"UNSUPPORTED\"}]]\n,[5,[\"onPlaylistModeChanged\",{\"shuffleEnabled\":\"false\",\"loopEnabled\":\"false\"}]]\n]\n"

Where c is the sessionID, and S is the gsessionID.
Hard to parse this into json unmarshallable form.
*/
const sessionIDRegex = `"c",\s*?"(.*?)",\"`
const gSessionIDRegex = `"S",\s*?"(.*?)"]`
const screenIDsKey = "screen_ids"

//TODO. This will handle the internal operations of executing commands against the chromecast.
type youtubeSession struct {
}

//YoutubeController is the controller for the commands unique to the dashcast.
type YoutubeController struct {
	connection     *mediaConnection
	screenID       string
	youtubeSession youtubeSession
	sessionID      string
	gSessionID     string
	loungeID       string
	incoming       chan *string
	requestCounter generic.Counter
	sessionCounter generic.Counter
}

type requestComponents struct {
	URL    string
	Body   interface{}
	Header req.Header
	Params req.Param
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

	loungeToken, err := c.getLoungeToken(screenID)
	c.loungeID = loungeToken
	if err != nil {
		spew.Dump("Failed to get lounge token")
	}

	c.bind(screenID, loungeToken)
	c.initializeQueue("cwQgjq0mCdE", "")
}

func (c *YoutubeController) bind(screenID *string, loungeToken string) error {
	c.requestCounter.Reset()
	c.sessionCounter.Reset()
	requestID := c.requestCounter.GetAndIncrement()
	request := createBindRequest(requestID, loungeToken)
	response, err := postBindRequest(request)
	sessionID, gSessionID, err := parseResponse(response)
	if err != nil {
		return err
	}
	c.sessionID = sessionID
	c.gSessionID = gSessionID
	//assign self data to sid and gsiddata.
	return nil

}

func postBindRequest(reqComponents requestComponents) (*req.Resp, error) {
	return req.Post(
		reqComponents.URL,
		reqComponents.Header,
		reqComponents.Params,
		reqComponents.Body,
	)
}

func parseResponse(bindResponse *req.Resp) (sessionID string, gSessionID string, err error) {
	responseString, err := bindResponse.ToString()
	spew.Dump(responseString)
	if err != nil {
		return "", "", err
	}
	sessionID, err = parseSessionID(responseString)
	if err != nil {
		return "", "", err
	}
	gSessionID, err = parseGSessionID(responseString)
	if err != nil {
		return "", "", err
	}
	return sessionID, gSessionID, err
}

func parseSessionID(bindResponse string) (string, error) {
	regex, err := regexp.Compile(sessionIDRegex)
	if err != nil {
		return "", err
	}
	matches := regex.FindStringSubmatch(bindResponse)
	if len(matches) <= 0 {
		return "", errors.New("Failed to find sessionID inside bind response")
	}
	return matches[1], nil
}

func parseGSessionID(bindResponse string) (string, error) {
	regex, err := regexp.Compile(gSessionIDRegex)
	if err != nil {
		return "", err
	}
	matches := regex.FindStringSubmatch(bindResponse)
	if len(matches) <= 0 {
		return "", errors.New("Failed to find GSessionID inside bind response")
	}
	return matches[1], nil
}

func createBindRequest(requestID int, loungeToken string) requestComponents {
	header := req.Header{
		loungeIDHeader: loungeToken,
	}
	for k, v := range defaultHeaders {
		header[k] = v
	}

	params := req.Param{
		requestIDKey: requestID,
		versionKey:   bindVersion,
		cVersionKey:  bindCVersion,
	}

	requestComponents := requestComponents{
		URL:    bindURL,
		Body:   url.Values(bindData).Encode(),
		Header: header,
		Params: params,
	}
	return requestComponents
}

func (c *YoutubeController) getLoungeToken(screenID *string) (string, error) {
	payload := url.Values{screenIDsKey: {*screenID}}.Encode()
	response, err := req.Post(loungeTokenURL,
		req.Header(defaultHeaders),
		payload)
	if err != nil {
		spew.Dump("Failed to get the lounge ID", err)
		return "", err
	}
	tokenResponse := &youtube.LoungeTokenResponse{}
	err = response.ToJSON(tokenResponse)
	if err != nil {
		spew.Dump("Failed to unmarshal the token response")
		return "", err
	}
	return tokenResponse.Screens[0].LoungeToken, nil
}

//TODO
func (c *YoutubeController) PlayVideo(videoID int) {
	//id, err := c.getLoungeID(c.screenID)
	//if err != nil {
	//	c.loungeID = string(id)
	//}
	//c.bind()
}

func (c *YoutubeController) basicPostRequest(payload url.Values, headers map[string]string, parameters map[string]string) (*http.Response, error) {
	requestHeaders := defaultHeaders
	for k, v := range headers {
		requestHeaders[k] = v
	}

	request, err := http.NewRequest("POST", loungeTokenURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, err
	}

	for k, v := range requestHeaders {
		request.Header.Set(k, v)
	}

	return http.DefaultClient.Do(request)
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
	request := c.createInitializeQueueRequest(videoID, listID)
	spew.Dump("Request info", request)
	response, err := c.sendInitializeQueueRequest(request)
	if err != nil {
		spew.Dump("Failed to initialize queue")
	}
	spew.Dump("Got response", response)
}

//TODO
func (c *YoutubeController) createInitializeQueueRequest(videoID string, listID string) requestComponents {
	requestCount := c.requestCounter.GetAndIncrement()
	header := req.Header{
		loungeIDHeader: c.loungeID,
	}

	for k, v := range defaultHeaders {
		header[k] = v
	}

	params := req.Param{
		sessionIDKey:  c.sessionID,
		gSessionIDKey: c.gSessionID,
		requestIDKey:  requestCount,
		versionKey:    bindVersion,
		cVersionKey:   bindCVersion,
	}
	index := strconv.Itoa(defaultIndex)
	count := strconv.Itoa(defaultCount)
	body := map[string][]string{
		listIDKey:       []string{listID},
		actionKey:       []string{actionSetPlaylist},
		currentTimeKey:  []string{defaultTime},
		currentIndexKey: []string{index},
		audioOnlyKey:    []string{defaultAudioOnlySetting},
		videoIDKey:      []string{videoID},
		countKey:        []string{count},
	}
	spew.Dump("body", body)
	formattedBody := formatSessionParameters(body, c.sessionCounter.GetAndIncrement())
	spew.Dump("Formatted body", formattedBody)
	return requestComponents{
		URL:    bindURL,
		Header: header,
		Params: params,
		Body:   url.Values(formattedBody).Encode(),
	}

}

func formatSessionParameters(params map[string][]string, requestCount int) map[string][]string {
	formattedMap := make(map[string][]string)
	requestPrefix := fmt.Sprintf(requestPrefixFormat, requestCount)
	for key, value := range params {
		newKey := key
		if strings.HasPrefix(newKey, "_") {
			newKey = requestPrefix + newKey
		}
		formattedMap[newKey] = value
	}
	return formattedMap
}

//TODO
func (c *YoutubeController) sendInitializeQueueRequest(reqComponents requestComponents) (*req.Resp, error) {
	return req.Post(
		reqComponents.URL,
		reqComponents.Header,
		reqComponents.Params,
		reqComponents.Body,
	)
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
