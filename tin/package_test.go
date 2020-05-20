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

func (p packageManagerMock) Installed() ([]Package, error) {
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
	want := StateSubscription{}
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

func TestPackagesEqualTrue(t *testing.T) {
	a := Packages{
		Package{Name: "Name", Version: "1.0.0"},
		Package{Name: "Package", Version: "1.0.1"},
	}
	b := Packages{
		Package{Name: "Name", Version: "1.0.0"},
		Package{Name: "Package", Version: "1.0.1"},
	}

	want := true
	got := a.Equal(b)
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestPackagesEqualFalse(t *testing.T) {
	tt := []struct {
		a Packages
		b interface{}
	}{
		{
			a: Packages{Package{Name: "Name", Version: "1.0.0"}},
			b: Packages{Package{Name: "Name", Version: "1.0.1"}},
		},
		{
			a: Packages{Package{Name: "Name", Version: "1.0.0"}},
			b: Packages{
				Package{Name: "Name", Version: "1.0.1"},
				Package{Name: "Name", Version: "1.0.2"},
			},
		},
		{
			a: Packages{Package{Name: "Name", Version: "1.0.0"}},
			b: Packages{Package{Name: "Package", Version: "1.0.0"}},
		},
		{
			a: Packages{Package{Name: "Name", Version: "1.0.0"}},
			b: "package",
		},
	}

	for _, tc := range tt {
		want := false
		got := tc.a.Equal(tc.b)

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}

func TestPackageEqualTrue(t *testing.T) {
	a := Package{Name: "Name", Version: "1.0.0"}
	b := Package{Name: "Name", Version: "1.0.0"}

	want := true
	got := a.Equal(b)
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestPackageEqualFalse(t *testing.T) {
	tt := []struct {
		a Package
		b interface{}
	}{
		{
			a: Package{Name: "Name", Version: "1.0.0"},
			b: Package{Name: "Name", Version: "1.0.1"},
		},
		{
			a: Package{Name: "Name", Version: "1.0.0"},
			b: Package{Name: "Package", Version: "1.0.0"},
		},
		{
			a: Package{Name: "Name", Version: "1.0.0"},
			b: "package",
		},
	}

	for _, tc := range tt {
		want := false
		got := tc.a.Equal(tc.b)

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}

func TestPackageCountEqualTrue(t *testing.T) {
	a := PackageCount(7)
	b := PackageCount(7)

	want := true
	got := a.Equal(b)
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestPackageCountEqualFalse(t *testing.T) {
	tt := []struct {
		a PackageCount
		b interface{}
	}{
		{
			a: PackageCount(7),
			b: PackageCount(17),
		},
		{
			a: PackageCount(7),
			b: 7,
		},
	}

	for _, tc := range tt {
		want := false
		got := tc.a.Equal(tc.b)

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}

func TestPackageInstalled(t *testing.T) {
	withState := NewPackageManagerService(nil, log.New(ioutil.Discard, "", log.Flags()))
	withState.SetInstalled([]Package{{Name: "package", Version: "1.0.0"}})

	tt := []struct {
		service *PackageManagerService
		want    int
	}{
		{
			service: withState,
			want:    1,
		},
		{
			service: NewPackageManagerService(nil, log.New(ioutil.Discard, "", log.Flags())),
			want:    0,
		},
	}

	for _, tc := range tt {
		got := len(tc.service.Installed())

		if got != tc.want {
			t.Errorf("want %v, got %v", tc.want, got)
		}
	}
}
