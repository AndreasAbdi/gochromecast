package youtube

import (
	"errors"
	"net/url"
	"regexp"

	"github.com/AndreasAbdi/gochromecast/generic"
	"github.com/imroc/req"
)

/*
Requests to the bind API return objects of type.
"[[0,[\"c\",\"19AB39151763497F\",\"\",8]\n]\n,[1,[\"S\",\"d6CNYWDUZb40UcroBuzH6QZJti79F-mc\"]]\n,[2,[\"loungeStatus\",{\"devices\":\"[{\\\"app\\\":\\\"lb-v4\\\",\\\"capabilities\\\":\\\"dsp,que,mus\\\",\\\"clientName\\\":\\\"tvhtml5\\\",\\\"experiments\\\":\\\"\\\",\\\"name\\\":\\\"Chromecast\\\",\\\"id\\\":\\\"1ed072b4-b75a-4878-88d0-fe9e6625d9ec\\\",\\\"type\\\":\\\"LOUNGE_SCREEN\\\",\\\"hasCc\\\":\\\"true\\\"},{\\\"app\\\":\\\"GOCAST_REMOTE_APP\\\",\\\"pairingType\\\":\\\"cast\\\",\\\"capabilities\\\":\\\"que,mus\\\",\\\"clientName\\\":\\\"unknown\\\",\\\"experiments\\\":\\\"\\\",\\\"name\\\":\\\"21b78ce1-4311-4c5e-8ef5-0101eddf5671\\\",\\\"remoteControllerUrl\\\":\\\"\\\",\\\"id\\\":\\\"21b78ce1-4311-4c5e-8ef5-0101eddf5671\\\",\\\"type\\\":\\\"REMOTE_CONTROL\\\",\\\"localChannelEncryptionKey\\\":\\\"wMphRtC_eiqqMvJk61EWvN-k1rA7IA72NzG2KMqPxPU\\\"}]\"}]]\n,[3,[\"playlistModified\",{\"videoIds\":\"\"}]]\n,[4,[\"onAutoplayModeChanged\",{\"autoplayMode\":\"UNSUPPORTED\"}]]\n,[5,[\"onPlaylistModeChanged\",{\"shuffleEnabled\":\"false\",\"loopEnabled\":\"false\"}]]\n]\n"

Where c is the sessionID, and S is the gsessionID.
Hard to parse this into json unmarshallable form.
*/
const sessionIDRegex = `"c",\s*?"(.*?)",\"`
const gSessionIDRegex = `"S",\s*?"(.*?)"]`

//Spoof data to be passed for binding the screen and the token
const defaultDeviceType = "REMOTE_CONTROL"
const defaultDeviceName = "GOCAST_REMOTE_CONTROL"
const defaultDeviceID = "GOCAST"
const bindPairingType = "cast"
const defaultAppName = "GOCAST_REMOTE_APP"

var bindData = map[string][]string{
	"device":       {defaultDeviceType},
	"id":           {defaultDeviceID},
	"name":         {defaultDeviceID},
	"mdx-version":  {string(3)},
	"pairing_type": {bindPairingType},
	"app":          {defaultAppName},
}

//CreateBindRequest creates a bind request from relevant data.
func CreateBindRequest(requestID int, loungeToken string) generic.RequestComponents {
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

	requestComponents := generic.RequestComponents{
		URL:    bindURL,
		Body:   url.Values(bindData).Encode(),
		Header: header,
		Params: params,
	}
	return requestComponents
}

//ParseResponse attempts to grab relevant session info from a youtube screen/loungetoken bind request response.
func ParseResponse(bindResponse *req.Resp) (sessionID string, gSessionID string, err error) {
	responseString, err := bindResponse.ToString()
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

//parseSessionID attempts to grab the sessionid from a bindrequest response string.
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

//parseGSessionID attempts to grab the gsessionid from a bindrequest response string.
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
