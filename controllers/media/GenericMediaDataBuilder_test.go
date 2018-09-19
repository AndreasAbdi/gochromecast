package media

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/fatih/structs"
)

//testing to get error message
func TestDefaultBuild(t *testing.T) {
	builder, er := NewGenericMediaDataBuilder("1", "image/JPEG", NoneStreamType)
	if er != nil {
		t.Errorf("Failed to build a generic media data object.")
	}
	result, err := builder.Build()
	if err != nil {
		t.Error("Failed to build generic media data with basic data")
	}
	expectedMetadata := genericMediaMetadata{
		StandardMediaMetadata{
			genericMediaMetadataType,
			nil,
		},
		nil,
		nil,
		nil,
	}
	expected := MediaData{
		"1",
		"image/JPEG",
		string(NoneStreamType),
		nil,
		structs.Map(expectedMetadata),
		nil,
	}
	if !cmp.Equal(result, expected) {
		t.Errorf("Result != expected.\nResult: (%v). \nExpected: (%v)\n", result, expected)
	}
}

func TestPartialParameterizedBuild(t *testing.T) {

	testSubtitle := "By Blender Foundation"
	testTitle := "The Big Buck Bunny"
	testDuration := float64(0)
	expectedMetadata := genericMediaMetadata{
		StandardMediaMetadata{
			genericMediaMetadataType,
			&testTitle,
		},
		nil,
		&testSubtitle,
		nil,
	}
	expected := MediaData{
		"102312318",
		"video/mpeg",
		string(BufferedStreamType),
		&testDuration,
		structs.Map(expectedMetadata),
		nil,
	}

	builder, er := NewGenericMediaDataBuilder("102312318", "video/mpeg", BufferedStreamType)
	if er != nil {
		t.Errorf("Failed to build a generic media data object.")
	}

	builder.SetSubtitle(&testSubtitle)
	builder.SetDuration(&testDuration)
	builder.SetTitle(&testTitle)
	result, err := builder.Build()
	if err != nil {
		t.Error("Failed to build generic media data with basic data")
	}

	if !cmp.Equal(result, expected) {
		t.Errorf("Result != expected.\nResult: (%v). \nExpected: (%v)\n", result, expected)
	}
}

func TestFullParameterizedBuild(t *testing.T) {

	testSubtitle := "By Google"
	testTitle := "For Bigger Blazes"
	testDuration := float64(12312313)
	testReleaseDate := randomDate()
	testReleaseDateFormatted := string(testReleaseDate.UTC().Format(time.RFC3339))
	testImageURLs := make([]string, 1)
	testImageURLs[0] = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/ForBiggerBlazes.jpg"
	expectedMetadata := genericMediaMetadata{
		StandardMediaMetadata{
			genericMediaMetadataType,
			&testTitle,
		},
		testImageURLs,
		&testSubtitle,
		&testReleaseDateFormatted,
	}
	expected := MediaData{
		"102312318",
		"video/mpeg",
		string(BufferedStreamType),
		&testDuration,
		structs.Map(expectedMetadata),
		nil,
	}

	builder, err := NewGenericMediaDataBuilder("102312318", "video/mpeg", BufferedStreamType)
	if err != nil {
		t.Errorf("Failed to build a generic media data object.")
	}
	builder.SetTitle(&testTitle)
	builder.SetImageURLs(testImageURLs)
	builder.SetReleaseDate(&testReleaseDate)
	builder.SetSubtitle(&testSubtitle)
	builder.SetDuration(&testDuration)
	result, err := builder.Build()
	if err != nil {
		t.Error("Failed to build generic media data with basic data")
	}

	if !cmp.Equal(result, expected) {
		t.Errorf("Result != expected.\nResult:\n (%v). \nExpected:\n (%v)\n", result, expected)
	}
}

func randomDate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
