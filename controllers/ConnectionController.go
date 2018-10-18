package controllers

import (
	"github.com/AndreasAbdi/gochromecast/primitives"
)

//ConnectionController is for basic connect and close caommands connection to a chromecast
type ConnectionController struct {
	channel *primitives.Channel
}

var connect = primitives.PayloadHeaders{Type: "CONNECT"}
var close = primitives.PayloadHeaders{Type: "CLOSE"}

//NewConnectionController constructor for connection controllers
func NewConnectionController(client *primitives.Client, sourceID, destinationID string) *ConnectionController {
	controller := &ConnectionController{
		channel: client.NewChannel(sourceID, destinationID, connectionControllerNamespace),
	}

	return controller
}

//Connect is initial command to create a connection to a chromecast
func (c *ConnectionController) Connect() {
	c.channel.Send(connect)
}

//Close is the final command to close a connection to a chromecast
func (c *ConnectionController) Close() {
	c.channel.Send(close)
}
