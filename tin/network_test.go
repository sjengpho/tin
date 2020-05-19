package tin

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"reflect"
	"testing"
)

type essidLookupMock struct {
	returnError bool
}

func (n essidLookupMock) Lookup() (ESSID, error) {
	if n.returnError {
		return "", errors.New("error")
	}

	return "Network name", nil
}

type publicIPLookupMock struct {
	returnError bool
}

func (p publicIPLookupMock) Lookup() (PublicIP, error) {
	if p.returnError {
		return PublicIP{}, errors.New("error")
	}

	return PublicIP{}, nil
}

func TestNewNetworkService(t *testing.T) {
	tt := []struct {
		want *NetworkService
		got  *NetworkService
	}{
		{
			want: &NetworkService{},
			got:  NewNetworkService(essidLookupMock{}, publicIPLookupMock{}, log.New(ioutil.Discard, "", log.Flags())),
		},
		{
			want: &NetworkService{},
			got:  NewNetworkService(essidLookupMock{}, publicIPLookupMock{}, log.New(ioutil.Discard, "", log.Flags())),
		},
		{
			want: &NetworkService{},
			got:  NewNetworkService(essidLookupMock{returnError: true}, publicIPLookupMock{returnError: true}, log.New(ioutil.Discard, "", log.Flags())),
		},
		{
			want: &NetworkService{},
			got:  NewNetworkService(nil, nil, log.New(ioutil.Discard, "", log.Flags())),
		},
	}

	for _, tc := range tt {
		if reflect.TypeOf(tc.got) != reflect.TypeOf(tc.want) {
			t.Errorf("want %v, got %v", reflect.TypeOf(tc.want), reflect.TypeOf(tc.got))
		}
	}
}

func TestNetworkSubscribe(t *testing.T) {
	s := NewNetworkService(essidLookupMock{}, publicIPLookupMock{}, log.New(ioutil.Discard, "", log.Flags()))
	want := StateSubscription{}
	got := s.Subscribe()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}

func TestNetworkName(t *testing.T) {
	withState := NewNetworkService(nil, publicIPLookupMock{}, log.New(ioutil.Discard, "", log.Flags()))
	withState.SetName("Network name")

	tt := []struct {
		service *NetworkService
		want    ESSID
	}{
		{
			service: withState,
			want:    "Network name",
		},
		{
			service: NewNetworkService(nil, publicIPLookupMock{}, log.New(ioutil.Discard, "", log.Flags())),
			want:    "",
		},
	}

	for _, tc := range tt {
		got := tc.service.Name()

		if got != tc.want {
			t.Errorf("want %v, got %v", tc.want, got)
		}
	}
}

func TestNetworkSetName(t *testing.T) {
	s := NewNetworkService(nil, publicIPLookupMock{}, log.New(ioutil.Discard, "", log.Flags()))
	s.SetName("name")

	want := ESSID("name")
	got := s.Name()

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestNetworkIP(t *testing.T) {
	withState := NewNetworkService(nil, nil, log.New(ioutil.Discard, "", log.Flags()))
	withState.SetIP(net.IPv4(0, 0, 0, 0))

	tt := []struct {
		service *NetworkService
		want    string
	}{
		{
			service: withState,
			want:    "0.0.0.0",
		},
		{
			service: NewNetworkService(nil, nil, log.New(ioutil.Discard, "", log.Flags())),
			want:    "Unknown",
		},
	}

	for _, tc := range tt {
		got := tc.service.IP()

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("want %v, got %v", tc.want, got)
		}
	}
}
