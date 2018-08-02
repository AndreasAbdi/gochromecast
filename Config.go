package castv2

//ChromecastServiceName is the name of the service to lookup via mDNS for finding chromecast devices
const chromecastServiceName = "_googlecast._tcp"

const defaultChromecastReceiverID = "receiver-0"
const defaultChromecastSenderID = "sender-0"

//APPIDs are the hardcoded IDs for different applications in chromecast devices

//media receiver is a generic media player for urls. Can play images, videos, music, etc.
const mediaReceiverAppID string = "CC1AD845"

const youtubeAppID string = "233637DE"
const spotifyAppID string = "CC32E753"

//back drop is a back drop of images usually displayed as the default for when you run your chromecast
const backdropAppID string = "E8C28D3C"
