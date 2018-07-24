package controllers

import castv2 "github.com/AndreasAbdi/go-castv2"

type MediaCommand struct {
	castv2.PayloadHeaders
	MediaSessionID int `json:"mediaSessionId"`
}

type MediaStatusResponse struct {
	castv2.PayloadHeaders
	Status []*MediaStatus `json:"status,omitempty"`
}

//MediaData is data format for message to send to chromecast to play a (vid/image/tvshow/music video/etc) via generic media player.
//https://developers.google.com/cast/docs/reference/messages#MediaData is the general info.
type MediaData struct {
	castv2.PayloadHeaders
	URL               string
	ContentType       string
	Title             string
	ThumbnailURL      string
	CurrentTime       float64
	AutoPlay          bool
	StreamType        string
	Subtitles         string
	SubtitlesLanguage string
	SubtitlesMIME     string
	subtitlesID       int
	Metadata          map[string]interface{}
}

type StandardMediaMetadata struct {
	MetadataType int
	Title        string
}

//TODO
type GenericMediaMetadata struct {
	StandardMediaMetadata
	Images      []string
	Subtitle    string
	ReleaseDate string
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

type MediaStatus struct {
	castv2.PayloadHeaders
	MediaSessionID         int                    `json:"mediaSessionId"`
	PlaybackRate           float64                `json:"playbackRate"`
	PlayerState            string                 `json:"playerState"`
	CurrentTime            float64                `json:"currentTime"`
	SupportedMediaCommands int                    `json:"supportedMediaCommands"`
	Volume                 *Volume                `json:"volume,omitempty"`
	CustomData             map[string]interface{} `json:"customData"`
	IdleReason             string                 `json:"idleReason"`
}
