package castv2

import "net"

//Device Object to run basic chromecast commands
type Device struct {
	host net.IP
	port int
}

//NewDevice is constructor for Device struct
func NewDevice(host net.IP, port int) {

}
