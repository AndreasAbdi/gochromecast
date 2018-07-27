package castv2

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/AndreasAbdi/go-castv2/api"
	"github.com/gogo/protobuf/proto"
)

type Client struct {
	realConn *tls.Conn
	conn     *packetStream
	channels []*Channel
}

type PayloadHeaders struct {
	Type      string `json:"type"`
	RequestID *int   `json:"requestID,omitempty"`
}

func (h *PayloadHeaders) setRequestID(id int) {
	h.RequestID = &id
}

func (h *PayloadHeaders) getRequestID() int {
	return *h.RequestID
}

func NewClient(host net.IP, port int) (*Client, error) {

	log.Printf("connecting to %s:%d ...", host, port)

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Chromecast. Error:%s", err)
	}

	wrapper := NewPacketStream(conn)

	client := &Client{
		realConn: conn,
		conn:     wrapper,
		channels: make([]*Channel, 0),
	}

	go func() {
		for {
			packet := wrapper.Read()

			message := &api.CastMessage{}
			err = proto.Unmarshal(*packet, message)
			if err != nil {
				log.Fatalf("Failed to unmarshal CastMessage: %s", err)
			}

			var headers PayloadHeaders

			err := json.Unmarshal([]byte(*message.PayloadUtf8), &headers)

			if err != nil {
				log.Fatalf("Failed to unmarshal message: %s", err)
			}

			for _, channel := range client.channels {
				channel.receiveMessage(message, &headers)
			}

		}
	}()

	return client, nil
}

func (c *Client) Close() {
	c.realConn.Close()
}

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

func (c *Client) Send(message *api.CastMessage) error {

	proto.SetDefaults(message)

	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(&data)

	return err

}
