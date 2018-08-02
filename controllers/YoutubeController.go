package controllers

import (
	"github.com/AndreasAbdi/go-castv2/primitives"
)

/*
Helps with playing the youtube chromecast app.
See
https://github.com/balloob/pychromecast/blob/master/pychromecast/controllers/youtube.py
and
https://github.com/ur1katz/casttube/blob/master/casttube/YouTubeSession.py.
Essentially, you start a session with the website, and the website/session handles  like any other receiver app.
*/

//TODO. This will handle the internal operations of executing commands against the chromecast.
type youtubeSession struct {
}

//YoutubeController is the controller for the commands unique to the dashcast.
type YoutubeController struct {
	channel        *primitives.Channel
	screenID       int
	youtubeSession youtubeSession
}

//NewYoutubeController is a constructor for a dash cast controller.
func NewYoutubeController(client primitives.Client, sourceID, destinationID string) *DashCastController {
	return &DashCastController{
		channel: client.NewChannel(sourceID, destinationID, dashcastControllerNamespace),
	}
}

//TODO
func (c *YoutubeController) PlayVideo(VideoID int) {

}

//TODO
func (c *YoutubeController) PlayNext(VideoID int) {

}

//TODO
func (c *YoutubeController) AddToQueue(VideoID int) {

}

//TODO
func (c *YoutubeController) RemoveFromQueue(VideoID int) {

}
