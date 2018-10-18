# Information

## workflow

1. You discover your chromecast by running multicast DNS (local DNS service in intranet, may not work in private/cloud settings).
2. You create a commmunication channel to that chromecast for generic communications (launch application/start application/etc).
    - this channel needs to be kept alive via a heartbeat controller.
    - this channel needs to be started up/ terminated via a connection controller.

3. You start up an application via the the receiver controller that communicates over the generic communication channel.
4. You create a new communication channel to the chromecast specifically for controls of your new application (youtube, spotify, generic media application, etc).
    - this channel is the same as the previous in that it needs to be keptalive and started up/terminated.
    - this channel should be closed when your application is closed.

5. You use your new communication channel to do whatever it was that you wanted. (play youtube videos, wtv)

## On channels

Each communication channel has a specified namespace. This defines what type of communications happen over it.

```go
NAMESPACE_CONNECTION = 'urn:x-cast:com.google.cast.tp.connection'
NAMESPACE_RECEIVER = 'urn:x-cast:com.google.cast.receiver'
NAMESPACE_HEARTBEAT = 'urn:x-cast:com.google.cast.tp.heartbeat'
NAMESPACE_MEDIA = 'urn:x-cast:com.google.cast.media'
NAMESPACE_PLEX = 'urn:x-cast:plex'
NAMESPACE_CAST = 'urn:x-cast:com.url.cast'
NAMESPACE_YOUTUBE  = 'urn:x-cast:com.google.youtube.mdx'
```

channel communications are via tcp, and requests are composed of specified formats defined in the chromecast sdk(though we can't use the sdk, the formats are still consistent). The api folder contains the generic protocol information.

Finding other namespaces involves capture of network traffic from chrome browser/phone to chromecast. You can do this via looking at chrome dev tools or alternatively via `chrome://net-internals/#capture`, it should be under the `Tr@n$p0rt` identifier.
