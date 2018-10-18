package receiver

import "github.com/AndreasAbdi/gochromecast/primitives"

//StatusResponse is the main wrapper around a receiver status
type StatusResponse struct {
	primitives.PayloadHeaders
	Status *Status `json:"status,omitempty"`
}

//Status is the general struct containing the volume as well as the current application sessions.
type Status struct {
	primitives.PayloadHeaders
	Applications []*ApplicationSession `json:"applications"`
	Volume       *Volume               `json:"volume,omitempty"`
}

//Volume is the struct containing current volume information
type Volume struct {
	Level *float64 `json:"level,omitempty"`
	Muted *bool    `json:"muted,omitempty"`
}

//ApplicationSession is a descript of the application
type ApplicationSession struct {
	AppID       *string      `json:"appId,omitempty"`
	DisplayName *string      `json:"displayName,omitempty"`
	Namespaces  []*Namespace `json:"namespaces"`
	SessionID   *string      `json:"sessionId,omitempty"`
	StatusText  *string      `json:"statusText,omitempty"`
	TransportID *string      `json:"transportId,omitempty"`
}

//Namespace is the channel namespace of the application
type Namespace struct {
	Name string `json:"name"`
}
