package media

import (
	"mime"
)

type contentType string

//NewContentType creates a new contentType object from some string.
func NewContentType(cTypeString string) (contentType, error) {
	var cType contentType
	_, _, err := mime.ParseMediaType(cTypeString)
	if err != nil {
		return cType, err
	}
	return contentType(cTypeString), nil
}
