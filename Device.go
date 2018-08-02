package castv2

import (
	"net"

	"github.com/AndreasAbdi/go-castv2/primitives"
)

//Device Object to run basic chromecast commands
type Device struct {
	client primitives.Client
}

//NewDevice is constructor for Device struct
func NewDevice(host net.IP, port int) {

}
