package media

import (
	"net/url"
	"time"

	"github.com/fatih/structs"
)

//Perhaps considering error returns if builder fails to add operation.

const genericMediaMetadataType int = 0
const MovieMetadataType int = 1
const TvShowMetadataType int = 2
const MusicTrackMetadataType int = 3
const PhotoMediaMetadataType int = 4

//genericMediaDataBuilder component for building up mediadatabuilders
type genericMediaDataBuilder struct {
	standardMediaDataBuilder
	imageURLs   []string
	subtitle    *string
	title       *string
	releaseDate *time.Time
}

//NewGenericMediaDataBuilder is a constructor for the media data builder class
func NewGenericMediaDataBuilder(contentID contentID, contentType contentType, sType streamType) (genericMediaDataBuilder, error) {
	builder := genericMediaDataBuilder{}
	builder.SetContentID(contentID)
	builder.SetContentType(contentType)
	builder.SetStreamType(sType)
	return builder, nil
}

//SetImageURLs the display image slideshow urls for the media message
func (builder *genericMediaDataBuilder) SetImageURLs(imageURLs []string) {
	for index := 0; index < len(imageURLs); index++ {
		_, err := url.ParseRequestURI(imageURLs[index])
		if err != nil {
			return
		}
	}
	builder.imageURLs = imageURLs
}

//SetTitle sets the builder's title
func (builder *genericMediaDataBuilder) SetTitle(title *string) {
	builder.title = title
}

//SetSubtitleURL sets the builder's subtitle
func (builder *genericMediaDataBuilder) SetSubtitle(subtitle *string) {
	builder.subtitle = subtitle
}

//SetReleaseDate sets release date for element.
func (builder *genericMediaDataBuilder) SetReleaseDate(releaseDate *time.Time) {
	builder.releaseDate = releaseDate
}

func convertDateToISO8601(date *time.Time) *string {
	if date == nil {
		return nil
	}
	formatted := date.UTC().Format(time.RFC3339)
	return &formatted
}

//Build returns a standard mediadata object from its current data.
func (builder *genericMediaDataBuilder) Build() (MediaData, error) {
	date := convertDateToISO8601(builder.releaseDate)
	metadata := genericMediaMetadata{
		StandardMediaMetadata{
			genericMediaMetadataType,
			builder.title,
		},
		builder.imageURLs,
		builder.subtitle,
		date,
	}
	return MediaData{
		builder.contentID,
		builder.contentType,
		builder.streamType,
		builder.duration,
		structs.Map(metadata),
		builder.customData,
	}, nil
}
