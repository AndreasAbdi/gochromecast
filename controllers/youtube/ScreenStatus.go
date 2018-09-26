package youtube

import (
	"github.com/AndreasAbdi/go-castv2/primitives"
)

//ScreenStatus is returned status of the youtube chromecast channel.
type ScreenStatus struct {
	primitives.PayloadHeaders
	Data ScreenStatusData `json:"data"`
}

//ScreenStatusData is the internal data of the returned status.
type ScreenStatusData struct {
	//ScreenId for the screen that has a youtube session.
	ScreenID string `json:"screenId"`
}
