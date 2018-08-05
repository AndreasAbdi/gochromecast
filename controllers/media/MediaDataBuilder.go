package media

//Perhaps considering error returns if builder fails to add operation.

type StreamType string

const NoneStreamType StreamType = "NONE"
const BufferedStreamType StreamType = "BUFFERED"
const LiveStreamType StreamType = "LIVE"

//MediaDataBuilder interface for generic media player media data messages to pass to chromecast.
type MediaDataBuilder interface {
	SetContentID(string)
	SetContentType(string)
	SetStreamType(string)
	SetDuration(float64)
	SetCustomData(map[string]interface{})
	Build() (MediaData, error)
}

//GenericMediaDataBuilder component for building up mediadatabuilders
type StandardMediaDataBuilder struct {
	contentID   string
	contentType string
	streamType  string
	duration    float64
	customData  map[string]interface{}
}

//SetContentID sets the contentID field
func (builder *StandardMediaDataBuilder) SetContentID(id string) {
	builder.contentID = id
}

//SetContentType sets the contentType field
func (builder *StandardMediaDataBuilder) SetContentType(cType string) {
	builder.contentType = cType
}

//SetStreamType sets the streamType field
func (builder *StandardMediaDataBuilder) SetStreamType(sType string) {
	switch StreamType(sType) {
	case NoneStreamType, BufferedStreamType, LiveStreamType:
		builder.streamType = sType
	}
}

//SetDuration sets the contentType field
func (builder *StandardMediaDataBuilder) SetDuration(duration float64) {
	if duration < 0 {
		return
	}
	builder.duration = duration
}

//SetCustomData sets the custom data field
func (builder *StandardMediaDataBuilder) SetCustomData(customData map[string]interface{}) {
	builder.customData = customData
}

//Build returns a standard mediadata object from its current data.
func (builder *StandardMediaDataBuilder) Build() (MediaData, error) {
	return MediaData{
		builder.contentID,
		builder.contentType,
		builder.streamType,
		builder.duration,
		nil,
		builder.customData,
	}, nil
}
