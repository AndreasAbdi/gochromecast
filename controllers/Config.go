package controllers

const receiverControllerNamespace string = "urn:x-cast:com.google.cast.receiver"
const mediaControllerNamespace string = "urn:x-cast:com.google.cast.media"
const heartbeatControllerNamespace string = "urn:x-cast:com.google.cast.tp.heartbeat"
const connectionControllerNamespace string = "urn:x-cast:com.google.cast.tp.connection"

//TODO: Consider if the better name would be like system event blah. (systemeventgetstatus, or systemeventclose)
const receiverControllerSystemEventGetStatus string = "GET_STATUS"
const receiverControllerSystemEventSetVolume string = "SET_VOLUME"
const receiverControllerSystemEventReceiverStatus string = "RECEIVER_STATUS"
const receiverControllerSystemEventLaunch string = "LAUNCH"
const receiverControllerSystemEventLaunchError string = "LAUNCH_ERROR"

const heartbeatControllerSystemEventPing string = "PING"
const heartbeatControllerSystemEventPong string = "PONG"

const connectionControllerSystemEventConnect string = "CONNECT"
const connectionControllerSystemEventClose string = "CLOSE"
