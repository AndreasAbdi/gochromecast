package castv2

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AndreasAbdi/go-castv2/api"
	"github.com/AndreasAbdi/go-castv2/generic"
)

type Channel struct {
	client        *Client
	sourceId      string
	DestinationId string
	namespace     string
	counter       generic.Counter
	inFlight      map[int]chan *api.CastMessage
	listeners     []channelListener
}

type channelListener struct {
	responseType string
	callback     func(*api.CastMessage)
}

type hasRequestId interface {
	setRequestId(id int)
	getRequestId() int
}

/*
	General workflow is
	1. message is sent via the request method with unique requestID. adds an inflight chan to wait for event to return.
	2. request method calls the send method and wraps the call around some stuff.
	3. Client processes events in its active socket stream. If the requestID matches previous one then it sends the
	unmarshalled message to the message function.
	4. message function processes the event and passes to inflight chan. Also calls any attached listener functions.
	5. Request method returns the unmarshalled chromecast event if it worked, timeout if it didn't receive the event in time.
*/

//Processes message that has been received by the packetstream in the client.
func (c *Channel) message(message *api.CastMessage, headers *PayloadHeaders) {

	if *message.DestinationId != "*" && (*message.SourceId != c.DestinationId || *message.DestinationId != c.sourceId || *message.Namespace != c.namespace) {
		return
	}

	if *message.DestinationId != "*" && headers.RequestId != nil {
		listener, ok := c.inFlight[*headers.RequestId]
		if !ok {
			log.Printf("Warning: Unknown incoming response id: %d to destination:%s", *headers.RequestId, c.DestinationId)
			return
		}
		listener <- message
		delete(c.inFlight, *headers.RequestId)
		return
	}

	if headers.Type == "" {
		log.Printf("Warning: No message type. Don't know what to do. headers:%v message:%v", headers, message)
		return
	}

	for _, listener := range c.listeners {
		if listener.responseType == headers.Type {
			listener.callback(message)
		}
	}

}

//OnMessage adds a callback listener for a certain type of message with specified responseType
func (c *Channel) OnMessage(responseType string, cb func(*api.CastMessage)) {
	c.listeners = append(c.listeners, channelListener{responseType, cb})
}

/*
Send creates a simple message to be sent by the client.
*/
func (c *Channel) Send(payload interface{}) error {

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	payloadString := string(payloadJson)

	message := &api.CastMessage{
		ProtocolVersion: api.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &c.sourceId,
		DestinationId:   &c.DestinationId,
		Namespace:       &c.namespace,
		PayloadType:     api.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadString,
	}

	return c.client.Send(message)
}

/*
Request sends a payload and returns the message the chromecast gives back.
Throws an error if timeout has passed waiting for the message to be returned.
*/
func (c *Channel) Request(payload hasRequestId, timeout time.Duration) (*api.CastMessage, error) {

	payload.setRequestId(c.counter.IncrementAndGet())

	response := make(chan *api.CastMessage)

	c.inFlight[payload.getRequestId()] = response

	err := c.Send(payload)

	if err != nil {
		delete(c.inFlight, payload.getRequestId())
		return nil, err
	}

	select {
	case reply := <-response:
		return reply, nil
	case <-time.After(timeout):
		delete(c.inFlight, payload.getRequestId())
		return nil, fmt.Errorf("Call to cast channel %s - timed out after %d seconds", c.DestinationId, timeout/time.Second)
	}

}
