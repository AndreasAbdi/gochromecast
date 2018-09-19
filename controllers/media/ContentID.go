package media

import (
	"fmt"
)

type contentID string

//see https://developers.google.com/cast/docs/reference/messages#MediaData
const maxContentIDLength = 1000

//ContentIDLengthError is error caused by input string for contentID constructor being longer than max.
type ContentIDLengthError struct {
	length int
}

func (er *ContentIDLengthError) Error() string {
	return fmt.Sprintf("Length of string too long (%v), must be shorter than %v.", er.length, maxContentIDLength)
}

//NewContentID creates a new contentID object from some string.
func NewContentID(contentIDString string) (contentID, error) {
	if len(contentIDString) > maxContentIDLength {
		return "", &ContentIDLengthError{len(contentIDString)}
	}

	return contentID(contentIDString), nil
}
