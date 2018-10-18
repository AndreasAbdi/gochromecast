package youtube

import (
	"net/url"
	"strconv"

	"github.com/AndreasAbdi/gochromecast/generic"
	"github.com/imroc/req"
)

type actionRequestParameters struct {
	VideoID             string
	ActionID            string
	LoungeID            string
	RequestCount        int
	SessionRequestCount int
	SessionID           string
	GSessionID          string
}

//CreateActionRequest to be sent to active session.
func createActionRequest(params actionRequestParameters) generic.RequestComponents {
	requestCount := params.RequestCount
	header := req.Header{
		loungeIDHeader: params.LoungeID,
	}

	for k, v := range defaultHeaders {
		header[k] = v
	}

	reqParams := req.Param{
		sessionIDKey:  params.SessionID,
		gSessionIDKey: params.GSessionID,
		requestIDKey:  requestCount,
		versionKey:    bindVersion,
		cVersionKey:   bindCVersion,
	}

	count := strconv.Itoa(defaultCount)
	body := map[string][]string{
		actionKey:  []string{params.ActionID},
		videoIDKey: []string{params.VideoID},
		countKey:   []string{count},
	}
	formattedBody := FormatSessionParameters(body, params.SessionRequestCount)
	return generic.RequestComponents{
		URL:    bindURL,
		Header: header,
		Params: reqParams,
		Body:   url.Values(formattedBody).Encode(),
	}
}
