package youtube

import (
	"fmt"
	"strings"

	"github.com/AndreasAbdi/go-castv2/generic"
)

const requestPrefixFormat = "req%d"

//Session represents a connection to the youtube chromecast api.
type Session struct {
	requestCounter generic.Counter
	sessionCounter generic.Counter
}

func (s *Session) bind() {

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
