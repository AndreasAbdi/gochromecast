package youtube

import (
	"net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/imroc/req"
)

const screenIDsKey = "screen_ids"

func GetLoungeToken(screenID string) (string, error) {
	payload := url.Values{screenIDsKey: {screenID}}.Encode()
	response, err := req.Post(
		loungeTokenURL,
		req.Header(defaultHeaders),
		payload)
	if err != nil {
		spew.Dump("Failed to get the lounge ID", err)
		return "", err
	}
	tokenResponse := LoungeTokenResponse{}
	err = response.ToJSON(&tokenResponse)
	if err != nil {
		spew.Dump("Failed to unmarshal the token response")
		return "", err
	}
	return tokenResponse.Screens[0].LoungeToken, nil
}
