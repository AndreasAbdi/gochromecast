package receiver

import "github.com/AndreasAbdi/gochromecast/primitives"

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
	primitives.PayloadHeaders
	Status *ReceiverStatus `json:"status,omitempty"`
}

type ReceiverStatus struct {
	primitives.PayloadHeaders
	Applications []*ApplicationSession `json:"applications"`
	Volume       *Volume               `json:"volume,omitempty"`
}

type LaunchRequest struct {
	primitives.PayloadHeaders
	AppID *string `json:"appId,omitempty"`
}

type StopRequest struct {
	primitives.PayloadHeaders
	SessionID *string `json:"sessionID,omitempty"`
}

//GetSessionByNamespace attempts to return the first session with a specified namespace.
func (s *ReceiverStatus) GetSessionByNamespace(namespace string) *ApplicationSession {

	for _, app := range s.Applications {
		for _, ns := range app.Namespaces {
			if ns.Name == namespace {
				return app
			}
		}
	}

	return nil
}
