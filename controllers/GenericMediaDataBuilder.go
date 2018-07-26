package controllers

import (
	"fmt"
	"net/url"
	"time"

	"github.com/fatih/structs"
)

//Perhaps considering error returns if builder fails to add operation.

//GenericMediaDataBuilder component for building up mediadatabuilders
type GenericMediaDataBuilder struct {
	StandardMediaDataBuilder
	imageURLs   []string
	subtitleURL string
	title       string
	releaseDate time.Time
}

//SetImageURLs the display image slideshow urls for the media message
func (builder *GenericMediaDataBuilder) SetImageURLs(imageURLs []string) {
	for index := 0; index < len(imageURLs); index++ {
		_, err := url.ParseRequestURI(imageURLs[index])
		if err != nil {
			return
		}
	}
	builder.imageURLs = imageURLs
}

//SetTitle sets the builder's title
func (builder *GenericMediaDataBuilder) SetTitle(title string) {
	builder.title = title
}

//SetSubtitleURL sets the builder's subtitle URL
func (builder *GenericMediaDataBuilder) SetSubtitleURL(subtitleURL string) {
	_, err := url.ParseRequestURI(subtitleURL)
	if err != nil {
		return
	}
	builder.subtitleURL = subtitleURL
}

//SetReleaseDate sets release date for element.
func (builder *GenericMediaDataBuilder) SetReleaseDate(releaseDate time.Time) {
	builder.releaseDate = releaseDate
}

func convertDateToISO8601(date time.Time) string {
	return fmt.Sprint(date.UTC().Format(time.RFC3339))
}

//Build returns a standard mediadata object from its current data.
func (builder *GenericMediaDataBuilder) Build() (MediaData, error) {
	metadata := GenericMediaMetadata{
		StandardMediaMetadata{
			GenericMediaMetadataType,
			builder.title,
		},
		builder.imageURLs,
		builder.subtitleURL,
		convertDateToISO8601(builder.releaseDate),
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
