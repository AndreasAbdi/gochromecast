package youtube

import "github.com/AndreasAbdi/go-castv2/primitives"

//Request for getting a screen ID for an existing youtube application.
type GetScreenIDRequest struct {
	primitives.PayloadHeaders
	ScreenID int `json:"screen_ids"`
}
