package main

import (
	"errors"
	"regexp"
	"strings"
	"time"

	castv2 "github.com/AndreasAbdi/go-castv2"
	"github.com/AndreasAbdi/go-castv2/configs"
	"github.com/davecgh/go-spew/spew"
)

// A simple example, showing how to play a youtube video.
func main() {
	devices := make(chan *castv2.Device, 100)
	castv2.FindDevices(time.Second*5, devices)
	for device := range devices {
		appID := configs.YoutubeAppID
		device.ReceiverController.LaunchApplication(&appID, time.Second*5, false)
		device.YoutubeController.Test()
		//device.TestYoutube("some fake url")
	}

}

func playFindTheUnbind() {
	//parseString := []byte("{\"screens\":[{\"screenId\":\"4dnsm3coi2p9psaiugc548lv96\",\"loungeToken\":\"AGdO5p_E-j5833kbaHb8mupgjxgS-J0ovj1dTosF1BsSm_J7s4DQZ_MjoqnDUl-wO7laDweBu6kSHwRrir0S4bGfl7CXskMpmdlB-VVkmQc7-lBJvT7kExs\",\"expiration\":1539027648154}]}")
	// tokenResponse := &youtube.LoungeTokenResponse{}
	// json.Unmarshal(parseString, tokenResponse)
	// spew.Dump(tokenResponse)
}

func playFindTheString() {
	parseString := "892\n[[0,[\"c\",\"19AB39151763497F\",\"\",8]\n]\n,[1,[\"S\",\"d6CNYWDUZb40UcroBuzH6QZJti79F-mc\"]]\n,[2,[\"loungeStatus\",{\"devices\":\"[{\\\"app\\\":\\\"lb-v4\\\",\\\"capabilities\\\":\\\"dsp,que,mus\\\",\\\"clientName\\\":\\\"tvhtml5\\\",\\\"experiments\\\":\\\"\\\",\\\"name\\\":\\\"Chromecast\\\",\\\"id\\\":\\\"1ed072b4-b75a-4878-88d0-fe9e6625d9ec\\\",\\\"type\\\":\\\"LOUNGE_SCREEN\\\",\\\"hasCc\\\":\\\"true\\\"},{\\\"app\\\":\\\"GOCAST_REMOTE_APP\\\",\\\"pairingType\\\":\\\"cast\\\",\\\"capabilities\\\":\\\"que,mus\\\",\\\"clientName\\\":\\\"unknown\\\",\\\"experiments\\\":\\\"\\\",\\\"name\\\":\\\"21b78ce1-4311-4c5e-8ef5-0101eddf5671\\\",\\\"remoteControllerUrl\\\":\\\"\\\",\\\"id\\\":\\\"21b78ce1-4311-4c5e-8ef5-0101eddf5671\\\",\\\"type\\\":\\\"REMOTE_CONTROL\\\",\\\"localChannelEncryptionKey\\\":\\\"wMphRtC_eiqqMvJk61EWvN-k1rA7IA72NzG2KMqPxPU\\\"}]\"}]]\n,[3,[\"playlistModified\",{\"videoIds\":\"\"}]]\n,[4,[\"onAutoplayModeChanged\",{\"autoplayMode\":\"UNSUPPORTED\"}]]\n,[5,[\"onPlaylistModeChanged\",{\"shuffleEnabled\":\"false\",\"loopEnabled\":\"false\"}]]\n]\n"
	spew.Dump(len(parseString))
	parts := strings.SplitN(parseString, "\n", 2)

	spew.Dump("message length is %v", parts[0])
	spew.Dump("remaining message is parts: %v", parts[1])
	regex, err := regexp.Compile(`"c","(.*?)",\"`)
	if err != nil {
		spew.Dump("Bad regex for session id")
	}
	matches := regex.FindStringSubmatch(parseString)
	if len(matches) == 0 {
		return
	}
	sessionId := matches[1]
	findGSessionID(parseString)
	spew.Dump("sessionID", sessionId)
}

func findGSessionID(line string) error {
	regex, err := regexp.Compile(`"S",\s*"(.*?)"]`)
	if err != nil {
		spew.Dump("Bad regex for session id")
	}
	matches := regex.FindStringSubmatch(line)
	if len(matches) == 0 {
		return errors.New("Failed to find gsessionID")
	}
	gsessionID := matches[1]
	spew.Dump("gsessionID", gsessionID)
	return nil
}
