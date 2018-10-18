Try to figure out controls for the chromecast via golang. 
Then we can set up a remote server with a custom alexa skill that calls it. 

notes on running chromecast actions via golang. 

# useful references

## links

[CSharp chromecast API](https://github.com/Tapanila/SharpCaster)

[Node basic communication chromecast API](https://github.com/thibauts/node-castv2)

[Node chromecast API](https://github.com/thibauts/node-castv2-client)

[python chromecast API](https://github.com/balloob/pychromecast)

[golang basic communication chromecast API](https://github.com/ninjasphere/go-castv2)

[Google chromecast reference docs](https://developers.google.com/cast/docs/reference/)

[MediaData/metadata formats](https://developers.google.com/cast/docs/reference/messages#MediaData)

## Information

- no default supported API for chromecast except for the chromecast play button, ios device, and android device. Have to use a different language (nodecast, pychromecast, and go-cast are available).
- use https://github.com/ninjasphere/go-castv2/ for now. 
- chromecasts have a concept called namespaces (need to figure out what that is)
- chromecasts need to first be discovered. Then once discovered you can send requests to them. device discover runs via multicast DNS (local DNS service in intranet, may not work in private/cloud settings)
- interactions are via tcp, and requests are composed of specified formats defined in the chromecast sdk(though we can't use the sdk, the formats are still consistent).
- each request must pass through a channel that needs to be keptalive. The request must then have a specified namespace for the request. 

```
NAMESPACE_CONNECTION = 'urn:x-cast:com.google.cast.tp.connection'
NAMESPACE_RECEIVER = 'urn:x-cast:com.google.cast.receiver'
NAMESPACE_HEARTBEAT = 'urn:x-cast:com.google.cast.tp.heartbeat'
NAMESPACE_MEDIA = 'urn:x-cast:com.google.cast.media'
Plex Channel (urn:x-cast:plex)
Web Channel (urn:x-cast:com.url.cast)
YouTube Channel (urn:x-cast:com.google.youtube.mdx)
```

are the generic namespaces for communications. 
- finding other namespaces would involve screwing with the network capture of media to chromecast. You can do this via looking at chrome dev tools or alternatively via `chrome://net-internals/#capture`, it should be under the `Tr@n$p0rt` identifier. 

Receiver is the one you'd use to communicate with the platform for application running. (run youtube, netflix, etc)


- channels can be targeted for the entire chromecast platform - read `the actual dongle` - rather than the application currently running. 

- to play, you would use a mediacontroller abstraction object and use its commands. These controller objects have channels with the devices for communication. 

- hmmm, so the golang cast thing doesn't permit that i invoke requests directly through the channel because it is, you know, private. May need to modify it so that we can send more requests to it. 

- so to invoke a media request, you send a media request with the receiver controller / receiver channel. 
- the request type is of format https://developers.google.com/cast/docs/reference/messages#MediaData
- it'd be nice if we didn't have to declare MIME types tho. QQ. 

- when sending a request, you need to set a media player. 
- the total list can be seen in "https://clients3.google.com/cast/chromecast/device/baseconfig"
Some default useful ones are 

```
APP_BACKDROP = "E8C28D3C"
APP_YOUTUBE = "233637DE"
APP_MEDIA_RECEIVER = "CC1AD845"
APP_PLEX = "06ee44ee-e7e3-4249-83b6-f5d0b6f07f34_1"
APP_DASHCAST = "84912283"
APP_SPOTIFY = "CC32E753"
```

the media receiver is the easiest one to play with.

- to test, lets try sending an image, as well as a video link.
- parts to test, use the receiver controller to launch an app. use the receiver controller to then play the application. 

- chromecast sdk requires that you sign up. 