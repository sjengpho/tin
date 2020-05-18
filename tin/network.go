package tin

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

// ESSIDLookup is the interface implemented by an object that can
// lookup the network name.
type ESSIDLookup interface {
	Lookup() (ESSID, error)
}

// ESSID represents the network name.
type ESSID string

// PublicIPLookup is the interface implemented by an object that can
// lookup the IP, City and Country.
type PublicIPLookup interface {
	Lookup() (PublicIP, error)
}

// PublicIP represents the IP address.
type PublicIP = net.IP

// Represents tin.StateKeys.
const (
	NetworkName StateKey = "NetworkName"
	IP                   = "IP"
)

// NetworkService provides network information.
type NetworkService struct {
	nameLookup     ESSIDLookup
	publicIPLookup PublicIPLookup
	state          *State
	nameWorker     *Worker
	publicIPWorker *Worker
	logger         *log.Logger
}

// NewNetworkService returns tin.NetworkService.
func NewNetworkService(n ESSIDLookup, p PublicIPLookup, l *log.Logger) *NetworkService {
	s := &NetworkService{
		nameLookup:     n,
		publicIPLookup: p,
		state:          NewState(),
		logger:         l,
	}

	// Worker that lookup the network name on intervals and updates the state.
	if n == nil {
		s.logger.Println(errors.New("failed initializing worker"))
	} else {
		s.nameWorker = NewWorker(time.Minute, func() {
			name, err := s.nameLookup.Lookup()
			if err != nil {
				s.logger.Println(fmt.Errorf("worker failed: %w", err))
			} else {
				s.SetName(name)
			}
		})
	}

	// Worker that lookup the public IP, city and country on intervals and updates the state.
	if p == nil {
		s.logger.Println(errors.New("failed initializing public IP lookup worker"))
	} else {
		s.publicIPWorker = NewWorker(time.Minute, func() {
			publicIP, err := s.publicIPLookup.Lookup()
			if err != nil {
				s.logger.Println(fmt.Errorf("worker failed: %w", err))
			} else {
				s.SetIP(publicIP)
			}
		})
	}

	return s
}

// Subscribe returns a read-only channel.
func (s *NetworkService) Subscribe() <-chan StateValue {
	return s.state.Subscribe()
}

// Name returns a string.
func (s *NetworkService) Name() ESSID {
	v, err := s.state.Get(NetworkName)
	if err != nil {
		return ""
	}

	return v.(ESSID)
}

// SetName updates the state.
func (s *NetworkService) SetName(n ESSID) {
	s.state.Set(NetworkName, n)
}

// IP returns a string.
func (s *NetworkService) IP() string {
	v, err := s.state.Get(IP)
	if err != nil {
		return "Unknown"
	}

	return v.(PublicIP).String()
}

// SetIP updates the state.
func (s *NetworkService) SetIP(ip PublicIP) {
	s.state.Set(IP, ip)
}
