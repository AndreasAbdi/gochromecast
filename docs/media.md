# Media controller construction

The media controller allows you to run commands that are specific to the generic media controller for the chromecast.

This application allows you to play videos, images, and songs so long as the resource is publicly available/accessible via the chromecast device.

Communications specific to the generic media application are done via a separate communication channel like all chromecast applications, so the media controller builds a new chromecast channel whenever it is deployed.

Controls that are supported by the media controller include

- play
- pause
- skip
- stop
- next
- seek
- load video/image/song

The media controller is also in charge of running these media control commands for other applications (youtube, netflix, spotify), so you can use this controller for those use cases as well. 

Launching the actual generic media application is run by the receiver controller. So if you want to watch a video, you'd launch the application via the receiver controller, then you'd run a load command via the media controller. Interestingly enough, subtitle controls are supposed to be via receiver controller.

See examples folder to find out how to use.