package controllers

import (
	"time"

	"github.com/AndreasAbdi/go-castv2"
	"github.com/AndreasAbdi/go-castv2/api"
)

//Sends pings and wait for pongs - https://github.com/thibauts/node-castv2-client/blob/master/lib/controllers/heartbeat.js

const interval = time.Second * 5
const timeoutFactor = 3 // timeouts after 3 intervals
//TODO: TimeoutFactor is essentially ignored. We need to change that so that we perform something on timeout.

//HeartbeatController is used to maintain a connection to a chromecast via sending keepalive messages.
type HeartbeatController struct {
	ticker      *time.Ticker
	channel     *castv2.Channel
	pongChannel chan *api.CastMessage
}

var ping = castv2.PayloadHeaders{Type: SystemEventPing}
var pong = castv2.PayloadHeaders{Type: SystemEventPong}

//NewHeartbeatController is a constructor for a heartbeat controller.
func NewHeartbeatController(client *castv2.Client, sourceID, destinationID string) *HeartbeatController {
	controller := &HeartbeatController{
		channel: client.NewChannel(sourceID, destinationID, heartbeatControllerNamespace),
	}
	controller.channel.OnMessage(SystemEventPing, controller.onPing)

	return controller
}

func (c *HeartbeatController) onPing(_ *api.CastMessage) {
	c.channel.Send(pong)
}

//TODO
func (c *HeartbeatController) onTimeout() {

}

/*Start begins the keepalive event stream.
Essentially, we send a ping event, then the chromecast will start sending pongs back.
We would then need to consistently return ping events every specified interval period.
*/
func (c *HeartbeatController) Start() {

	c.ticker = time.NewTicker(interval)
	go func() {
		for {
			<-c.ticker.C
			c.channel.Send(ping)
		}
	}()
	//Process ping events or timeout.
	//TODO
	go func() {
		for {
			select {
			case <-time.After(timeoutFactor * interval):
				return
			case <-c.pongChannel:
				return
			}
		}
	}()
}

//Stop maintaining the keepalive.
func (c *HeartbeatController) Stop() {

	if c.ticker != nil {
		c.ticker.Stop()
		c.ticker = nil
	}

}
