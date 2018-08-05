package controllers

/*
After deliberation, this should probably contain configs that are separate from the actual media files. Mostly because I don't really know when google is going to change the API or status messages. Is it a good idea? IDK. suggestions from the internet are to keep constants as close as possible to their usages. That makes sense, but these are sorta related to the message formats that google are using, so I'm reluctant.
*/

const receiverControllerNamespace string = "urn:x-cast:com.google.cast.receiver"
const mediaControllerNamespace string = "urn:x-cast:com.google.cast.media"
const heartbeatControllerNamespace string = "urn:x-cast:com.google.cast.tp.heartbeat"
const connectionControllerNamespace string = "urn:x-cast:com.google.cast.tp.connection"
const dashcastControllerNamespace string = "urn:x-cast:com.madmod.dashcast"
const youtubeControllerNamespace string = "urn:x-cast:com.google.youtube.mdx"

//TODO: Consider if the better name would be like system event blah. (systemeventgetstatus, or systemeventclose)
const receiverControllerSystemEventGetStatus string = "GET_STATUS"
const receiverControllerSystemEventSetVolume string = "SET_VOLUME"
const receiverControllerSystemEventReceiverStatus string = "RECEIVER_STATUS"
const receiverControllerSystemEventLaunch string = "LAUNCH"
const receiverControllerSystemEventStop string = "STOP"
const receiverControllerSystemEventLaunchError string = "LAUNCH_ERROR"

const SystemEventPing string = "PING"
const SystemEventPong string = "PONG"

const connectionControllerSystemEventConnect string = "CONNECT"
const connectionControllerSystemEventClose string = "CLOSE"

//Media play requests need a stream type. These are the valid stream types.

//media receiver is a generic media player for urls. Can play images, videos, music, etc.
const mediaReceiverAppID string = "CC1AD845"

const youtubeAppID string = "233637DE"
const spotifyAppID string = "CC32E753"

//back drop is a back drop of images usually displayed as the default for when you run your chromecast.
const backdropAppID string = "E8C28D3C"
