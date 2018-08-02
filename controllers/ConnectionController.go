package controllers

import "github.com/AndreasAbdi/go-castv2"

//ConnectionController is for basic connect and close caommands connection to a chromecast
type ConnectionController struct {
	channel *castv2.Channel
}

var connect = castv2.PayloadHeaders{Type: "CONNECT"}
var close = castv2.PayloadHeaders{Type: "CLOSE"}

//NewConnectionController constructor for connection controllers
func NewConnectionController(client *castv2.Client, sourceID, destinationID string) *ConnectionController {
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
