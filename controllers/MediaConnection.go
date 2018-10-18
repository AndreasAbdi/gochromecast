package controllers

import (
	"errors"
	"time"

	"github.com/AndreasAbdi/gochromecast/controllers/receiver"

	"github.com/AndreasAbdi/gochromecast/api"
	"github.com/AndreasAbdi/gochromecast/primitives"
)

type listener struct {
	responseType string
	cb           func(*api.CastMessage)
}

//mediaConnection is an abstraction over the more complicated process of setting up a connection for an application session. Treat it effectively as a channel with setup and teardown.
type mediaConnection struct {
	client               *primitives.Client
	receiverController   *ReceiverController
	namespace            string //chromecast namespace for a specific application type
	channel              *primitives.Channel
	connectionController *ConnectionController
	listeners            []listener
	sessionID            string
	sourceID             string //the id for the requesting application.
}

//NewMediaConnection constructor for application sessions.
func NewMediaConnection(client *primitives.Client, receiverController *ReceiverController, namespace string, sourceID string) *mediaConnection {
	connection := mediaConnection{
		client:             client,
		receiverController: receiverController,
		namespace:          namespace,
		sourceID:           sourceID,
	}
	return &connection
}

//Terminate the media connection for an application session's channel.
func (connection *mediaConnection) Terminate(timeout time.Duration) {
	if connection.connectionController != nil {
		connection.connectionController.Close()
	}
}

func (connection *mediaConnection) OnMessage(responseType string, cb func(*api.CastMessage)) {
	connection.listeners = append(connection.listeners, listener{responseType, cb})
	if connection.channel != nil {
		connection.channel.OnMessage(responseType, cb)
	}
}

func (connection *mediaConnection) Request(payload primitives.HasRequestID, timeout time.Duration) (*api.CastMessage, error) {
	err := connection.ensureConnectionActive()
	if err != nil {
		return nil, err
	}
	return connection.channel.Request(payload, timeout)
}

func (connection *mediaConnection) ensureConnectionActive() error {
	appSession, err := connection.getAppSession()
	if err != nil {
		return err
	}
	connection.sessionID = *appSession.SessionID
	if connection.shouldResetConnection(appSession) {
		return connection.refreshConnection()
	}
	return nil
}

func (connection *mediaConnection) getAppSession() (*receiver.ApplicationSession, error) {
	getStatusCh := make(chan bool)
	var session *receiver.ApplicationSession
	go func() {
		status := <-connection.receiverController.Incoming
		session = status.GetSessionByNamespace(connection.namespace)

		getStatusCh <- true
	}()
	_, err := connection.receiverController.GetStatus(defaultTimeout)
	if err != nil {
		return nil, err
	}
	<-getStatusCh
	if session == nil {
		return nil, errors.New("No session by that name")
	}
	return session, nil
}

func (connection *mediaConnection) shouldResetConnection(session *receiver.ApplicationSession) bool {
	if session == nil {
		return true
	}
	return connection.channel == nil || connection.sessionID != *session.SessionID
}

func (connection *mediaConnection) refreshConnection() error {
	connection.performCleanup()
	session, err := connection.getAppSession()
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("Failed to generate a connection")
	}
	connection.setup(*session.TransportId)
	return nil
}

func (connection *mediaConnection) performCleanup() {
	if connection.connectionController != nil {
		connection.connectionController.Close()
	}
	connection.connectionController = nil
}

func (connection *mediaConnection) setup(transportID string) {
	connection.connectionController = setupConnectionController(connection.client, connection.sourceID, transportID)
	connection.channel = setupChannel(connection.client, connection.channel, connection.sourceID, transportID, connection.namespace, connection.listeners)
}

func setupConnectionController(client *primitives.Client, sourceID string, transportID string) *ConnectionController {
	connectionController := NewConnectionController(client, sourceID, transportID)
	connectionController.Connect()
	return connectionController
}

func setupChannel(client *primitives.Client, channel *primitives.Channel, sourceID string, transportID string, namespace string, listeners []listener) *primitives.Channel {
	if channel == nil {
		channel = client.NewChannel(
			sourceID,
			transportID,
			namespace)
		for _, listener := range listeners {
			channel.OnMessage(listener.responseType, listener.cb)
		}
	} else {
		channel.DestinationID = transportID
	}
	return channel
}
