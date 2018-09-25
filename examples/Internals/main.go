package main

//internal example of how the controllers work together.

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/AndreasAbdi/go-castv2/controllers"
	"github.com/AndreasAbdi/go-castv2/primitives"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/mdns"
)

// A simple example, showing how to find a Chromecast using mdns, and request its status.
func main() {

	castService := "_googlecast._tcp"

	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {

			if !strings.Contains(entry.Name, castService) {
				return
			}

			fmt.Printf("Got new chromecast: %v\n", entry)

			client, err := primitives.NewClient(entry.Addr, entry.Port)

			if err != nil {
				log.Fatalf("Failed to connect to chromecast %s", entry.Addr)
			}

			//_ = controllers.NewHeartbeatController(client, "Tr@n$p0rt-0", "Tr@n$p0rt-0")

			heartbeat := controllers.NewHeartbeatController(client, "sender-0", "receiver-0")
			heartbeat.Start()

			connection := controllers.NewConnectionController(client, "sender-0", "receiver-0")
			connection.Connect()

			receiver := controllers.NewReceiverController(client, "sender-0", "receiver-0")

			testLaunch(client, receiver)

			if err != nil {
				spew.Dump("Failed to launch application:", err)
			}
			response, err := receiver.GetStatus(time.Second * 5)

			spew.Dump("Status response", response, err)
		}
	}()

	go func() {
		mdns.Query(&mdns.QueryParam{
			Service: castService,
			Timeout: time.Second * 30,
			Entries: entriesCh,
		})
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
}

func testLaunch(client *primitives.Client, receiver *controllers.ReceiverController) {
	date := "CC1AD845"
	receiver.LaunchApplication(&date, time.Second*5, false)
	media := controllers.NewMediaController(client, "sender-0", receiver)
	media.Load("http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4", "video/mp4", time.Second*5)

	time.Sleep(time.Second * 5)
	media.Pause(time.Second * 10)

	time.Sleep(time.Second * 5)
	media.Play(time.Second * 10)

	time.Sleep(time.Second * 10)
	media.Stop(time.Second * 10)

}
