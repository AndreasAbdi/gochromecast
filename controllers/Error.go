package controllers

//InitializationError struct for initialization failure of youtube/chromecast applications.
type InitializationError struct {
}

func (e InitializationError) Error() string {
	return "Failed to initialize"
}
