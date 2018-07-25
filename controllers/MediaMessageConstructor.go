package controllers

type MediaMessageConstructor struct{}

func (constructor MediaMessageConstructor) CreateImageMediaMessage(url string, mimeType string, streamType StreamType) MediaData {
	return nil
}
