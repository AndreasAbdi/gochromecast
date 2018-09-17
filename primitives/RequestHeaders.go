package primitives

//PayloadHeaders are general request components for any message to be passed to a chromecast.
type PayloadHeaders struct {
	Type      string `json:"type"`
	RequestID *int   `json:"requestID,omitempty"`
}

func (h *PayloadHeaders) setRequestID(id int) {
	h.RequestID = &id
}

func (h *PayloadHeaders) getRequestID() int {
	return *h.RequestID
}
