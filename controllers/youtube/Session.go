package youtube

import (
	"fmt"
	"strings"

	"github.com/AndreasAbdi/go-castv2/generic"
)

const requestPrefixFormat = "req%d"

//Session represents a connection to the youtube chromecast api.
type Session struct {
	sessionID      string
	gSessionID     string
	requestCounter generic.Counter
	sessionCounter generic.Counter
}

//Bind a screen with a loungetoken and link operations to this session object.
func (s *Session) Bind(screenID *string, loungeToken string) error {
	s.requestCounter.Reset()
	s.sessionCounter.Reset()
	requestID := s.requestCounter.GetAndIncrement()
	request := CreateBindRequest(requestID, loungeToken)
	response, err := request.Post()
	sessionID, gSessionID, err := ParseResponse(response)
	if err != nil {
		return err
	}
	s.sessionID = sessionID
	s.gSessionID = gSessionID
	return nil
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
