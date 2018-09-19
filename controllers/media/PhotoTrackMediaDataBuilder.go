package media

import (
	"time"

	"github.com/fatih/structs"
)

//TODO
// type PhotoTrackMediaMetadata struct {
// 	StandardMediaMetadata
// 	Artist           string
// 	Location         string
// 	Latitude         float64
// 	Longitude        float64
// 	Width            int64
// 	Height           int64
// 	CreationDateTime string
// }

//const PhotoMediaMetadataType int = 4

//GenericMediaDataBuilder component for building up mediadatabuilders
type PhotoTrackMediaDataBuilder struct {
	standardMediaDataBuilder
	artist           string
	location         string
	latitude         float64
	longitude        float64
	width            int64
	height           int64
	title            string
	creationDateTime string
}

//SetTitle sets the builder's title
func (builder *PhotoTrackMediaDataBuilder) SetTitle(title string) {
	builder.title = title
}

//SetArtist sets the builder's artist
func (builder *PhotoTrackMediaDataBuilder) SetArtist(artist string) {
	builder.artist = artist
}

//TODO
func (builder *PhotoTrackMediaDataBuilder) SetLocation(location string) {
}

//TODO
func (builder *PhotoTrackMediaDataBuilder) SetLatitude(latitude float64) {
}

//TODO
func (builder *PhotoTrackMediaDataBuilder) SetLongitude(longitude float64) {

}

//SetArtist sets the builder's artist
func (builder *PhotoTrackMediaDataBuilder) SetWidth(width int64) {
	builder.width = width
}

//SetArtist sets the builder's artist
func (builder *PhotoTrackMediaDataBuilder) SetHeight(height int64) {
	builder.height = height
}

//TODO
func (builder *PhotoTrackMediaDataBuilder) SetCreationTime(creationTime time.Time) {
}

//Build returns a standard mediadata object from its current data.
func (builder *PhotoTrackMediaDataBuilder) Build() (MediaData, error) {

	metadata := PhotoTrackMediaMetadata{
		StandardMediaMetadata: StandardMediaMetadata{
			MetadataType: PhotoMediaMetadataType,
			Title:        &builder.title,
		},
		Artist:           builder.artist,
		Location:         builder.location,
		Latitude:         builder.latitude,
		Longitude:        builder.longitude,
		Width:            builder.width,
		Height:           builder.height,
		CreationDateTime: builder.creationDateTime,
	}
	return MediaData{
		ContentID:   builder.contentID,
		ContentType: builder.contentType,
		StreamType:  builder.streamType,
		Duration:    builder.duration,
		Metadata:    structs.Map(metadata),
		CustomData:  builder.customData,
	}, nil
}
