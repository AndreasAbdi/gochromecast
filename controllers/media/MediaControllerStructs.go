package media

import "github.com/AndreasAbdi/go-castv2/primitives"

type MediaCommand struct {
	primitives.PayloadHeaders
	MediaSessionID int `json:"mediaSessionId"`
}

type MediaStatusResponse struct {
	primitives.PayloadHeaders
	Status []*MediaStatus `json:"status,omitempty"`
}

type LoadCommand struct {
	primitives.PayloadHeaders
	Media       MediaData              `json:"media"`
	Autoplay    bool                   `json:"autoplay,omitempty"`
	CurrentTime float64                `json:"currentTime,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
}

//Generic enum type for media data

//StreamType is a type for media data defining what type of stream the data is supposed to be.
type streamType string

//NoneStreamType is for when you don't want to define a stream type.
const NoneStreamType streamType = "NONE"

//BufferedStreamType is for a stream that should be buffered/loaded.
const BufferedStreamType streamType = "BUFFERED"

//LiveStreamType is for videos that are livestreaming. (Twitch/Youtube livestreams, etc)
const LiveStreamType streamType = "LIVE"

//MediaData is data format for message to send to chromecast to play a (vid/image/tvshow/music video/etc) via generic media player.
//https://developers.google.com/cast/docs/reference/messages#MediaData is the general info.
type MediaData struct {
	//ContentID is the identifier for the content to be loaded by the current receiver application in the chromecast.
	//Usually this is just the URL.
	//ContentType is the MIME type of the media
	ContentID   string                 `json:"contentId"`
	ContentType string                 `json:"contentType"` //data MIME
	StreamType  string                 `json:"streamType"`  // (NONE, BUFFERED, or LIVE)
	Duration    *float64               `json:"duration,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"` //stores a mediadata
	CustomData  map[string]interface{} `json:"customData,omitempty"`
}

//StandardMediaMetadata is standard part for all metadata objects parts for mediadata objects.
type StandardMediaMetadata struct {
	MetadataType int     `json:"metadataType"`
	Title        *string `json:"title,omitempty"`
}

type genericMediaMetadata struct {
	StandardMediaMetadata
	Images      []string `json:"images,omitempty"`
	Subtitle    *string  `json:"subtitle,omitempty"`
	ReleaseDate *string  `json:"releaseDate,omitempty"`
}

//TODO
type MovieMediaMetadata struct {
	StandardMediaMetadata
	Studio      string
	Subtitle    string
	Images      []string
	ReleaseDate string
}

//TODO
type TvShowMediaMetadata struct {
	StandardMediaMetadata
	Images          []string
	SeriesTitle     string
	Season          int
	Episode         int
	OriginalAirDate string
}

//TODO
type MusicTrackMediaMetadata struct {
	StandardMediaMetadata
	Images      []string
	AlbumName   string
	AlbumArtist string
	Artist      string
	Composer    string
	TrackNumber int
	DiscNumber  int
	ReleaseDate string
}

//TODO
type PhotoTrackMediaMetadata struct {
	StandardMediaMetadata
	Artist           string
	Location         string
	Latitude         float64
	Longitude        float64
	Width            int64
	Height           int64
	CreationDateTime string
}

type Volume struct {
	Level *float64 `json:"level,omitempty"`
	Muted *bool    `json:"muted,omitempty"`
}

type MediaStatus struct {
	primitives.PayloadHeaders
	MediaSessionID         int                    `json:"mediaSessionId"`
	PlaybackRate           float64                `json:"playbackRate"`
	PlayerState            string                 `json:"playerState"`
	CurrentTime            float64                `json:"currentTime"`
	SupportedMediaCommands int                    `json:"supportedMediaCommands"`
	Volume                 *Volume                `json:"volume,omitempty"`
	CustomData             map[string]interface{} `json:"customData"`
	IdleReason             string                 `json:"idleReason"`
	Media                  MediaData              `json:"media,omitempty"`
}
