package media

import (
	"errors"
	"testing"
)

var tests = []struct {
	input       string
	expected    contentType
	errorResult error
}{
	{"", "", errors.New(`mime: no media type`)},
	{"MIMETYPE", "MIMETYPE", nil},
	{"image/JPEG", "image/JPEG", nil},
	{"audio/mpeg", "audio/mpeg", nil},
	{"audio/*", "audio/*", nil},
	{"image/gif", "image/gif", nil},
	{"video/mp4", "video/mp4", nil},
	{"video//mp4", "", errors.New("mime: expected token after slash")},
	{"video/", "", errors.New("mime: expected token after slash")},
	{`video//mp4; id="oc==jwioefj"`, "", errors.New("mime: expected token after slash")},
	{`Message/Partial; number=2; total=3; id="oc=jpbe0M2Yt4s@thumper.bellcore.com";`, `Message/Partial; number=2; total=3; id="oc=jpbe0M2Yt4s@thumper.bellcore.com";`, nil},
}

func TestBuildContentType(t *testing.T) {
	for num, test := range tests {
		actual, err := NewContentType(test.input)
		if err != test.errorResult {
			if err != nil && test.errorResult != nil && err.Error() == test.errorResult.Error() {
				continue
			}
			t.Errorf("Error for test %v returned different from expected. (%v) expected, (%v) received", num, test.errorResult, err)
		}
		if test.expected != actual {
			t.Errorf("Result for test %v returned different from expected. (%v) expected, (%v) received", num, test.expected, actual)
		}
	}
	t.Log("all tests completed")
}
