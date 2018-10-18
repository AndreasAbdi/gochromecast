package controllers

import (
	"fmt"
	"time"

	"github.com/AndreasAbdi/gochromecast/primitives"
	"github.com/AndreasAbdi/gochromecast/controllers/dashcast"
)

/*
	Related to http://stestagg.github.io/dashcast/
	Plays a url and then reloads it a few times. Rather than having to cast your own browser, this is easier.
*/

//DashCastController is the controller for the commands unique to the dashcast.
type DashCastController struct {
	channel *primitives.Channel
}

//NewDashCastController is a constructor for a dash cast controller.
func NewDashCastController(client primitives.Client, sourceID, destinationID string) *DashCastController {
	return &DashCastController{
		channel: client.NewChannel(sourceID, destinationID, dashcastControllerNamespace),
	}
}

/*Load loads a url, reloads it every few seconds. Forces launch of the url if you tell it to.
Setting force to True may load pages which block iframe embedding, but will prevent reload from
working and will cause calls to load_url() to reload the app.
*/
func (d *DashCastController) Load(url string, reloadTime time.Duration, forceLaunch bool) error {
	reload := !(forceLaunch || reloadTime == 0)
	if reload {
		reloadTime = 0
	}
	err := d.channel.Send(&dashcast.LoadCommand{
		URL: url, Force: forceLaunch, Reload: reload, ReloadTime: int64(reloadTime)})
	if err != nil {
		return fmt.Errorf("Failed to send play command: %s", err)
	}
	return nil
}
