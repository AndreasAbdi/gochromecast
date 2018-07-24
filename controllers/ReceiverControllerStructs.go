package controllers

import castv2 "github.com/AndreasAbdi/go-castv2"

//TODO: figure out which one of these we need to be public. receiverstatus and application session definitely.

type ApplicationSession struct {
	AppID       *string      `json:"appId,omitempty"`
	DisplayName *string      `json:"displayName,omitempty"`
	Namespaces  []*Namespace `json:"namespaces"`
	SessionID   *string      `json:"sessionId,omitempty"`
	StatusText  *string      `json:"statusText,omitempty"`
	TransportId *string      `json:"transportId,omitempty"`
}

type Namespace struct {
	Name string `json:"name"`
}

type Volume struct {
	Level *float64 `json:"level,omitempty"`
	Muted *bool    `json:"muted,omitempty"`
}

type StatusResponse struct {
	castv2.PayloadHeaders
	Status *ReceiverStatus `json:"status,omitempty"`
}

type ReceiverStatus struct {
	castv2.PayloadHeaders
	Applications []*ApplicationSession `json:"applications"`
	Volume       *Volume               `json:"volume,omitempty"`
}

type LaunchRequest struct {
	castv2.PayloadHeaders
	AppID *string `json:"appId,omitempty"`
}

type StopRequest struct {
	castv2.PayloadHeaders
	SessionID *string `json:"sessionID,omitempty"`
}
