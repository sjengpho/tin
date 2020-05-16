package cli

import (
	"fmt"
	"log"

	"github.com/sjengpho/tin/grpc"
)

// NewNetworkCommander returns a cli.NetworkCommander.
func NewNetworkCommander() NetworkCommander {
	return &networkCommander{}
}

// networkCommander implements cli.NetworkCommander.
type networkCommander struct{}

// ESSID outputs the network name.
func (s *networkCommander) ESSID(c *grpc.Client) {
	v, err := c.ESSID()
	if err != nil {
		log.Printf("failed getting the essid: %v", err)
		return
	}
	fmt.Println(v)
}

// IP outputs the IP address.
func (s *networkCommander) IP(c *grpc.Client) {
	v, err := c.IPAddress()
	if err != nil {
		log.Printf("failed getting the IP address: %v", err)
		return
	}
	fmt.Println(v)
}
