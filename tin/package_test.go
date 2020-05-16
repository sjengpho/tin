package tin

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

type packageManagerMock struct {
	returnError bool
}

func (p packageManagerMock) AvailableUpdates() ([]Package, error) {
	if p.returnError {
		return make([]Package, 0), errors.New("error")
	}

	return make([]Package, 1), nil
}

func TestNewPackageManagerService(t *testing.T) {
	tt := []struct {
		want *PackageManagerService
		got  *PackageManagerService
	}{
		{
			want: &PackageManagerService{},
			got:  NewPackageManagerService(packageManagerMock{returnError: false}, log.New(os.Stdout, "", log.Flags())),
		},
		{
			want: &PackageManagerService{},
			got:  NewPackageManagerService(packageManagerMock{returnError: true}, log.New(os.Stdout, "", log.Flags())),
		},
		{
			want: &PackageManagerService{},
			got:  NewPackageManagerService(nil, log.New(os.Stdout, "", log.Flags())),
		},
	}

	for _, tc := range tt {
		if reflect.TypeOf(tc.got) != reflect.TypeOf(tc.want) {
			t.Errorf("want %v, got %v", reflect.TypeOf(tc.want), reflect.TypeOf(tc.got))
		}
	}
}

func TestPackageSubscribe(t *testing.T) {
	s := NewPackageManagerService(packageManagerMock{}, log.New(os.Stdout, "", log.Flags()))
	want := make(<-chan interface{}, 1)
	got := s.Subscribe()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}

func TestPackageAvailableUpdatesCount(t *testing.T) {
	withState := NewPackageManagerService(nil, log.New(ioutil.Discard, "", log.Flags()))
	withState.SetAvailableUpdates(PackageCount(7))

	tt := []struct {
		service *PackageManagerService
		want    PackageCount
	}{
		{
			service: withState,
			want:    PackageCount(7),
		},
		{
			service: NewPackageManagerService(nil, log.New(ioutil.Discard, "", log.Flags())),
			want:    PackageCount(0),
		},
	}

	for _, tc := range tt {
		got := tc.service.AvailableUpdatesCount()

		if got != tc.want {
			t.Errorf("want %v, got %v", tc.want, got)
		}
	}
}
