# Youtube controller information

## Sample usage

## References

[PYCHROMECAST](https://github.com/balloob/pychromecast/blob/master/pychromecast/controllers/youtube.py)

[CASTTUBE](https://github.com/ur1katz/casttube/blob/master/casttube/YouTubeSession.py)

[GOTUBECAST](https://github.com/CBiX/gotubecast/blob/master/main.go)

[YOUTUBE-REMOTE](https://github.com/mutantmonkey/youtube-remote/blob/master/remote.py)

## Details

You create a screen object that can be used by the youtube lounge/leanback api to use using the chromecast connection. This has a unique screen ID.

Then you create a lounge token based on this screen ID using the loungeToken API.

Then you bind the loungeToken and the screenID using the bind API and then it'll return a gsession and session ID.

You then construct requests to the bind API to play videos/add videos to queue, etc.

Commands supported are:

- initialize playlist.
- add as next video in playlist.
- add video to back of playlist.
- remove video from playlist.