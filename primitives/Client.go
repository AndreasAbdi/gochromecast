package primitives

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"

	"github.com/AndreasAbdi/go-castv2/api"
	"github.com/gogo/protobuf/proto"
)

//Client is a basic connector to the chromecast event stream.
type Client struct {
	realConn *tls.Conn
	conn     *packetStream
	channels []*Channel
}

//NewClient is a constructor for a Client object. Host and Port are for the chromecast's network info.
func NewClient(host net.IP, port int) (*Client, error) {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Chromecast. Error:%s", err)
	}

	wrapper := newPacketStream(conn)

	client := &Client{
		realConn: conn,
		conn:     wrapper,
		channels: make([]*Channel, 0),
	}

	go func() {
		for {
			packet := wrapper.read()

			message := &api.CastMessage{}
			err = proto.Unmarshal(*packet, message)
			if err != nil {
				continue
			}

			var headers PayloadHeaders

			err := json.Unmarshal([]byte(*message.PayloadUtf8), &headers)

			if err != nil {
				continue
			}

			for _, channel := range client.channels {
				channel.receiveMessage(message, &headers)
			}

		}
	}()

	return client, nil
}

//Close closes the real socket connection that the client has to the chromecast.
func (c *Client) Close() {
	c.realConn.Close()
}

/*
NewChannel constructs a communication channel towards a chromecast device with a specified namespace.
namespace are those namespaces that are specified by the chromecast API.
*/
func (c *Client) NewChannel(sourceID, destinationID, namespace string) *Channel {
	channel := &Channel{
		client:        c,
		sourceID:      sourceID,
		DestinationID: destinationID,
		namespace:     namespace,
		listeners:     make([]channelListener, 0),
		inFlight:      make(map[int]chan *api.CastMessage),
	}

	c.channels = append(c.channels, channel)

	return channel
}

//Send sends a message to the chromecast using the actual socket connection.
func (c *Client) Send(message *api.CastMessage) error {

	proto.SetDefaults(message)

	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	_, err = c.conn.write(&data)

	return err

}
