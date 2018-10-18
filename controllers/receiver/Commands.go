package receiver

import "github.com/AndreasAbdi/gochromecast/primitives"

//LaunchRequest is a command to the receiver controller namespace to launch an application.
type LaunchRequest struct {
	primitives.PayloadHeaders
	AppID *string `json:"appId,omitempty"`
}

//StopRequest is a command to the receiver controller namespace to launch an application.
type StopRequest struct {
	primitives.PayloadHeaders
	SessionID *string `json:"sessionID,omitempty"`
}
