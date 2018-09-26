package youtube

import (
	"fmt"
	"strings"

	"github.com/AndreasAbdi/go-castv2/generic"
	"github.com/davecgh/go-spew/spew"
)

const requestPrefixFormat = "req%d"

//Session represents a connection to the youtube chromecast api.
type Session struct {
	screenID       string
	sessionID      string
	gSessionID     string
	loungeID       string
	requestCounter generic.Counter
	sessionCounter generic.Counter
}

//Bind a screen with a loungetoken and link operations to this session object.
func (s *Session) Bind(screenID *string, loungeToken string) error {
	s.resetCounters()
	requestID := s.requestCounter.GetAndIncrement()
	request := CreateBindRequest(requestID, loungeToken)
	response, err := request.Post()
	sessionID, gSessionID, err := ParseResponse(response)
	if err != nil {
		return err
	}
	s.assignVariables(*screenID, loungeToken, sessionID, gSessionID)
	return nil
}

func (s *Session) InitializeQueue(videoID string, listID string) {
	requestParams := s.createInitializeQueueRequestParams(videoID, listID)
	request := CreateInitializeQueueRequest(requestParams)
	spew.Dump("Request info", request)
	response, err := request.Post()
	if err != nil {
		spew.Dump("Failed to initialize queue")
	}
	spew.Dump("Got response", response)
}

func (s *Session) resetCounters() {
	s.requestCounter.Reset()
	s.sessionCounter.Reset()
}

func (s *Session) assignVariables(screenID string, loungeToken string, sessionID string, gSessionID string) {
	s.sessionID = sessionID
	s.gSessionID = gSessionID
	s.screenID = screenID
	s.loungeID = loungeToken
}
func (s *Session) createInitializeQueueRequestParams(videoID string, listID string) InitializeQueueRequestParams {
	return InitializeQueueRequestParams{
		VideoID:             videoID,
		ListID:              listID,
		LoungeID:            s.loungeID,
		RequestCount:        s.requestCounter.GetAndIncrement(),
		SessionRequestCount: s.sessionCounter.GetAndIncrement(),
		SessionID:           s.sessionID,
		GSessionID:          s.gSessionID,
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
