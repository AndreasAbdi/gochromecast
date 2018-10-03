package youtube

import (
	"fmt"
	"strings"

	"github.com/AndreasAbdi/go-castv2/generic"
	"github.com/imroc/req"
)

const requestPrefixFormat = "req%d"

const actionRemoveVideo = "removeVideo"
const actionInsertVideo = "insertVideo"
const actionAddVideo = "addVideo"
const actionClearQueue = "clearPlaylist"

//Session represents a connection to the youtube chromecast api.
type Session struct {
	screenID       string
	sessionID      *string
	gSessionID     *string
	loungeID       *string
	requestCounter generic.Counter
	sessionCounter generic.Counter
}

//NewSession constructs a new session object
func NewSession(screenID string) *Session {
	session := Session{
		screenID: screenID,
	}
	return &session
}

//StartSession initializes the session.
func (s *Session) StartSession() error {
	err := s.setLoungeID(s.screenID)
	if err != nil {
		return err
	}
	return s.bindAndSetVars(s.screenID, *s.loungeID)
}

//PlayNext adds a video to be played next in the current playlist
func (s *Session) PlayNext(videoID string) {
	s.sendAction(actionInsertVideo, videoID)
}

//ClearQueue removes all the videos from the current playlist
func (s *Session) ClearQueue() {
	s.sendAction(actionClearQueue, "")
}

//AddToQueue adds the video to the end of the current playlist
func (s *Session) AddToQueue(videoID string) {
	s.sendAction(actionAddVideo, videoID)
}

//RemoveFromQueue removes a video from the videoplaylist TODO
func (s *Session) RemoveFromQueue(videoID string) {
	s.sendAction(actionRemoveVideo, videoID)
}

//SendAction sends an action command.
func (s *Session) sendAction(actionType string, videoID string) {
	err := s.ensureSessionActive()
	if err != nil {
		return
	}
	actionParams := s.createActionRequestParameters(actionType, videoID)
	request := createActionRequest(actionParams)
	response, err := request.Post()
	if err != nil {
		return
	}
	s.handleBadResponse(response)
}

//InitializeQueue restarts the playlist to something else.
func (s *Session) InitializeQueue(videoID string, listID string) {
	requestParams := s.createInitializeQueueRequestParameters(videoID, listID)
	request := createInitializeQueueRequest(requestParams)
	response, err := request.Post()
	if err != nil {
		return
	}
	s.handleBadResponse(response)
}

func (s *Session) bindAndSetVars(screenID string, loungeID string) error {
	sessionID, gSessionID, err := s.bind(screenID, loungeID)
	if err != nil {
		return err
	}

	s.assignVariables(screenID, loungeID, sessionID, gSessionID)
	return nil
}

//Bind a screen and link operations to this session object.
func (s *Session) bind(screenID string, loungeID string) (sessionID string, gSessionID string, err error) {
	s.resetCounters()

	requestID := s.requestCounter.GetAndIncrement()
	request := CreateBindRequest(requestID, loungeID)
	response, err := request.Post()
	return ParseResponse(response)
}

func (s *Session) resetCounters() {
	s.requestCounter.Reset()
	s.sessionCounter.Reset()
}

func (s *Session) assignVariables(screenID string, loungeToken string, sessionID string, gSessionID string) {
	s.sessionID = &sessionID
	s.gSessionID = &gSessionID
	s.screenID = screenID
	s.loungeID = &loungeToken
}

func (s *Session) createActionRequestParameters(actionID string, videoID string) actionRequestParameters {
	return actionRequestParameters{
		VideoID:             videoID,
		ActionID:            actionID,
		LoungeID:            *s.loungeID,
		RequestCount:        s.requestCounter.GetAndIncrement(),
		SessionRequestCount: s.sessionCounter.GetAndIncrement(),
		SessionID:           *s.sessionID,
		GSessionID:          *s.gSessionID,
	}
}

func (s *Session) createInitializeQueueRequestParameters(videoID string, listID string) initializeQueueRequestParams {
	return initializeQueueRequestParams{
		VideoID:             videoID,
		ListID:              listID,
		LoungeID:            *s.loungeID,
		RequestCount:        s.requestCounter.GetAndIncrement(),
		SessionRequestCount: s.sessionCounter.GetAndIncrement(),
		SessionID:           *s.sessionID,
		GSessionID:          *s.gSessionID,
	}
}

func (s *Session) setLoungeID(screenID string) error {
	loungeToken, err := getLoungeToken(screenID)
	if err != nil {
		return err
	}
	s.loungeID = &loungeToken
	return nil
}

func (s *Session) ensureSessionActive() error {
	if s.inSession() {
		return s.StartSession()
	}
	return s.bindAndSetVars(s.screenID, *s.loungeID)

}

func (s *Session) inSession() bool {
	return s.loungeID != nil && s.gSessionID != nil
}

func (s *Session) handleBadResponse(response *req.Resp) {
	resp := response.Response()
	if resp == nil {
		return
	}
	if resp.StatusCode == 404 || resp.StatusCode == 400 {
		s.bindAndSetVars(s.screenID, *s.loungeID)
	}
}

//FormatSessionParameters formats session parameters to what youtube wants the keys to be.
func FormatSessionParameters(params map[string][]string, requestCount int) map[string][]string {
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
