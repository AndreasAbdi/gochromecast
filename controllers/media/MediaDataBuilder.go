package media

//Perhaps considering error returns if builder fails to add operation.

//MediaDataBuilder interface for generic media player media data messages to pass to chromecast.
type MediaDataBuilder interface {
	SetContentID(string)
	SetContentType(string)
	SetStreamType(string)
	SetDuration(float64)
	SetCustomData(map[string]interface{})
	Build() (MediaData, error)
}

//standardMediaDataBuilder component for building up mediadatabuilders
type standardMediaDataBuilder struct {
	contentID   string
	contentType string
	streamType  string
	duration    *float64
	customData  map[string]interface{}
}

//SetContentID sets the contentID field
func (builder *standardMediaDataBuilder) SetContentID(id contentID) {
	builder.contentID = string(id)
}

func (builder *standardMediaDataBuilder) SetContentType(contentType contentType) {
	builder.contentType = string(contentType)
}

//SetStreamType sets the streamType field
func (builder *standardMediaDataBuilder) SetStreamType(sType streamType) {
	switch sType {
	case NoneStreamType, BufferedStreamType, LiveStreamType:
		builder.streamType = string(sType)
	}
}

//SetDuration sets the contentType field
func (builder *standardMediaDataBuilder) SetDuration(duration *float64) {
	builder.duration = duration
}

//SetCustomData sets the custom data field
func (builder *standardMediaDataBuilder) SetCustomData(customData map[string]interface{}) {
	builder.customData = customData
}

//Build returns a standard mediadata object from its current data.
func (builder *standardMediaDataBuilder) Build() (MediaData, error) {
	return MediaData{
		builder.contentID,
		builder.contentType,
		builder.streamType,
		builder.duration,
		nil,
		builder.customData,
	}, nil
}
