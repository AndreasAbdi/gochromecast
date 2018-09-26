package youtube

import (
	"net/url"
	"strconv"

	"github.com/AndreasAbdi/go-castv2/generic"
	"github.com/imroc/req"
)

type ActionRequestParameters struct {
	VideoID             string
	actionID            string
	LoungeID            string
	RequestCount        int
	SessionRequestCount int
	SessionID           string
	GSessionID          string
}

//CreateActionRequest to be sent to active session.
func CreateActionRequest(params ActionRequestParameters) generic.RequestComponents {
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
		actionKey:  []string{params.actionID},
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
