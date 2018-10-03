package youtube

import (
	"net/url"

	"github.com/imroc/req"
)

const screenIDsKey = "screen_ids"

func getLoungeToken(screenID string) (string, error) {
	payload := url.Values{screenIDsKey: {screenID}}.Encode()
	response, err := req.Post(
		loungeTokenURL,
		req.Header(defaultHeaders),
		payload)
	if err != nil {
		return "", err
	}
	tokenResponse := LoungeTokenResponse{}
	err = response.ToJSON(&tokenResponse)
	if err != nil {
		return "", err
	}
	return tokenResponse.Screens[0].LoungeToken, nil
}
